package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	suite "github.com/stretchr/testify/suite"
	gorm "gorm.io/gorm"
)

type HandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *gorm.DB
}

func (suite *HandlerTestSuite) SetupSuite() {
	var err error
	suite.db, err = NewDB(WithSqliteInMemoryDB())
	assert.NoError(suite.T(), err)

	err = SeedDatabase(suite.db)
	assert.NoError(suite.T(), err)

	handler := NewRequestHandler(suite.db)

	suite.server = httptest.NewServer(handler.ServeHTTP())
}

func (suite *HandlerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHealthCheck() {
	resp, err := http.Get(suite.server.URL + "/health")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Equal(suite.T(), "OK", string(body))
}

func (suite *HandlerTestSuite) TestGetProducts() {
	resp, err := http.Get(suite.server.URL + "/products")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), len(products), 0)
}

func (suite *HandlerTestSuite) TestGetProductByID() {
	// First get all products to get a valid ID
	resp, err := http.Get(suite.server.URL + "/products")
	assert.NoError(suite.T(), err)
	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), len(products), 0)

	// Test with valid ID
	productID := products[0].ID
	resp, err = http.Get(suite.server.URL + "/product/" + fmt.Sprintf("%d", productID))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), productID, product.ID)

	// Test with invalid ID
	resp, err = http.Get(suite.server.URL + "/product/999999")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *HandlerTestSuite) TestCreateOrder() {
	// First get a product ID
	resp, err := http.Get(suite.server.URL + "/products")
	assert.NoError(suite.T(), err)
	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), len(products), 0)

	// Create order request
	orderReq := OrderReq{
		Items: []OrderItem{
			{
				ProductID: fmt.Sprintf("%d", products[0].ID),
				Quantity:  2,
			},
		},
	}

	jsonData, err := json.Marshal(orderReq)
	assert.NoError(suite.T(), err)

	// Test creating order
	resp, err = http.Post(suite.server.URL+"/order", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var order Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), order.ID)
	assert.Equal(suite.T(), 1, len(order.Items))
}

func (suite *HandlerTestSuite) TestCreateOrderWithInvalidData() {
	// Test with empty items
	orderReq := OrderReq{
		Items: []OrderItem{},
	}

	jsonData, err := json.Marshal(orderReq)
	assert.NoError(suite.T(), err)

	resp, err := http.Post(suite.server.URL+"/order", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	// Test with invalid product ID
	orderReq = OrderReq{
		Items: []OrderItem{
			{
				ProductID: "invalid_id",
				Quantity:  1,
			},
		},
	}

	jsonData, err = json.Marshal(orderReq)
	assert.NoError(suite.T(), err)

	resp, err = http.Post(suite.server.URL+"/order", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *HandlerTestSuite) TestGetOrders() {
	resp, err := http.Get(suite.server.URL + "/orders")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var orders []Order
	err = json.NewDecoder(resp.Body).Decode(&orders)
	assert.NoError(suite.T(), err)
}