package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error occured while initialising the database")
	}
	m.Run()
}

func createTable() {
	createTableQuery := ` CREATE TABLE IF NOT EXISTS products(
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		quantity int, 
		price float(10,7),
		PRIMARY KEY(id)
	);`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	// execute a query to delete everything from the table
	a.DB.Exec("DELETE FROM products;")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v','%v','%v')", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 5000)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, recieved: %v", expectedStatusCode, actualStatusCode)
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name": "chair", "quantity": 1, "price": 100}`)
	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, got %v", "chair", m["name"])
	}
	if m["quantity"] != 1 {
		t.Errorf("Expected quantity: %v, got %v", 1.0, m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest(http.MethodGet, "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest(http.MethodDelete, "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest(http.MethodGet, "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest(http.MethodGet, "/product/1", nil)
	response := sendRequest(req)

	var oldVal map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldVal)
	var product = []byte(`{"name": "connector", "quantity": 1, "price":10}`)
	req, _ = http.NewRequest(http.MethodPut, "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response = sendRequest(req)

	var newVal map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newVal)

	if oldVal["id"] != newVal["id"] {
		t.Errorf("Expected id: %v, got %v", oldVal["id"], newVal["id"])
	}

	if oldVal["price"] != newVal["price"] {
		t.Errorf("Expected price: %v, got %v", oldVal["price"], newVal["price"])
	}

	if oldVal["quantity"] == newVal["quantity"] {
		t.Errorf("Expected quantity: %v, got %v", oldVal["quantity"], newVal["quantity"])
	}
}
