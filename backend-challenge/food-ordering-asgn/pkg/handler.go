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
	mux.HandleFunc("GET /product/{productId}", h.GetProductByIDHandler)
	mux.HandleFunc("POST /orders", h.CreateOrderHandler)

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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create order logic here
	// TODO: Implement order creation

	w.WriteHeader(http.StatusCreated)
}