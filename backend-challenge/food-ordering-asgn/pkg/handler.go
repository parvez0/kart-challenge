package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type RequestHandler struct {
	db *gorm.DB
}

func NewRequestHandler(db *gorm.DB) *RequestHandler {
	return &RequestHandler{db: db}
}

func authorisationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if key := r.Header.Get("Authorization"); key != "apitest" {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	w.Write([]byte("Unauthorized"))
		// 	return
		// }
		next.ServeHTTP(w, r)
	})
}

func urlLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		logger.Info(fmt.Sprintf("Request to %s completed in %d ms", r.URL.Path, time.Since(startTime).Milliseconds()))
	})
}

func (h *RequestHandler) ServeHTTP() http.Handler {
	mux := http.NewServeMux()

	// Apply middleware chain
	handler := urlLoggingMiddleware(authorisationMiddleware(mux))

	// Register routes
	mux.HandleFunc("GET /health", h.HealthCheckHandler)
	mux.HandleFunc("GET /products", h.GetProductsHandler)
	mux.HandleFunc("GET /orders", h.GetOrdersHandler)
	mux.HandleFunc("GET /product/{productId}", h.GetProductByIDHandler)
	mux.HandleFunc("POST /order", h.CreateOrderHandler)

	return handler
}

func (h *RequestHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *RequestHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	var products []Product
	if err := h.db.Find(&products).Error; err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *RequestHandler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var placedOrders []Order
	if err := h.db.Preload("Items").Preload("Products").Find(&placedOrders).Error; err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(placedOrders)
}

func (h *RequestHandler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("productId")
	if productId == "" {
		http.Error(w, "Invalid ID supplied", http.StatusBadRequest)
		return
	}

	var product Product
	result := h.db.First(&product, productId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("No product found with id: %s", productId), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *RequestHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderReq OrderReq
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		logger.Error("Failed to parse order request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if orderReq.CouponCode != "" && !h.isCouponValid(orderReq.CouponCode) {
		http.Error(w, "Validation exception", http.StatusUnprocessableEntity)
		return
	}

	if len(orderReq.Items) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}

	// Get product IDs from request
	var productIDs []string
	for _, item := range orderReq.Items {
		if item.Quantity <= 0 {
			http.Error(w, "Quantity must be greater than zero", http.StatusBadRequest)
			return
		}
		productIDs = append(productIDs, item.ProductID)
	}

	var products []Product
	if err := h.db.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		logger.Error("Failed to fetch products:", err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	if len(products) != len(productIDs) {
		http.Error(w, "One or more products not found", http.StatusBadRequest)
		return
	}

	var order Order
	// Creating order in trasaction to avoid inconsistent state 
	// and rollback on failed order items.
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, item := range orderReq.Items {
			orderItem := OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&order).Association("Items").Append(&orderReq.Items); err != nil {
			return err
		}

		if err := tx.Model(&order).Association("Products").Append(&products); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error("Failed to create order:", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *RequestHandler) isCouponValid(code string) bool {
	var coupon Coupon
	result := h.db.Preload("SourceFile").Where("code = ?", code).First(&coupon)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Info("Coupon not found:", code)
			return false
		}
		logger.Error("Failed to verify coupon:", result.Error)
		return false
	}
	if len(coupon.SourceFile) < 2 {
		logger.Info("Invalid coupon code provided:", code)
		return false
	}
	return true
}