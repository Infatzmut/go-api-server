package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	Id       string
	Name     string
	Quantity int
	Price    float64
}

var Products []Product

func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint hit: products")
	// Convert the products data into json and write it as reponse
	json.NewEncoder(w).Encode(Products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint hit: product")
	vars := mux.Vars(r)
	key := vars["id"]
	for _, product := range Products {
		if product.Id == key {
			json.NewEncoder(w).Encode(product)
		}
	}
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/products", returnAllProducts)
	myRouter.HandleFunc("/product/{id}", getProduct)
	myRouter.HandleFunc("/", homePage)
	http.ListenAndServe("localhost:3000", nil)
}

func main() {

	Products = []Product{
		{Id: "1", Name: "Chair", Quantity: 100, Price: 100.00},
		{Id: "2", Name: "Desk", Quantity: 200, Price: 200.00},
	}
	handleRequest()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to homepage")
	fmt.Println("Endpoint hit: homepage")
}
