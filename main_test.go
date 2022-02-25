package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

const CREATE_STATEMENT = `CREATE TABLE IF NOT EXISTS product(
	id SERIAL,
	name TEXT NOT NULL,
	price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
	CONSTRAINT product_PID PRIMARY KEY(id))
	`

var app App

func TestMain(m *testing.M) {

	app.Initialise(
		os.Getenv("APP_DB_USER"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExist()
	code := m.Run()
	dbCleanUp()
	os.Exit(code)

}

/* 							Tests								*/
func TestEmptyTable(t *testing.T) {
	dbCleanUp()

	request, _ := http.NewRequest("GET", "/product/all", nil)
	response := executeRequest(request)

	checkStatusCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistingProduct(t *testing.T) {
	dbCleanUp()

	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(request)

	checkStatusCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	dbCleanUp()

	requestBody := []byte(`{"name: "test product", "price": 11.22}`)
	request, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")

	response := executeRequest(request)

	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != "11.22" {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}
}

func TestGetProduct(t *testing.T) {
	dbCleanUp()
	addProducts(1)

	request, _ := http.NewRequest("GET", "/product/1", nil)

	response := executeRequest(request)

	checkStatusCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {

	dbCleanUp()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkStatusCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	dbCleanUp()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

/* 					Helper Functions 							*/

func ensureTableExist() {
	if _, err := app.DB.Exec(CREATE_STATEMENT); err != nil {
		log.Fatal(err)
	}
}

func dbCleanUp() {
	app.DB.Exec("DELETE FROM product")
	app.DB.Exec("ALTER SEQUENCE product_id_seq RESTART WITH 1")
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {

	responseRecorder := httptest.NewRecorder()
	app.Router.ServeHTTP(responseRecorder, request)
	return responseRecorder
}

func checkStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected %d but got %d", expected, actual)
	}
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		app.DB.Exec("INSERT INTO product(name, price)  VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}
