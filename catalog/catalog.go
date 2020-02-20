package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type Product struct {
	UUID  string  `json: "uuid"`
	Name  string  `json: "name"`
	Price float64 `json: "price"`
}

type Products struct {
	Products []Product `json: "products"`
}

var productsUrl string

// Executes before main
func init() {
	// productsUrl := os.Getenv("PRODUCT_URL")
	productsUrl = "http://localhost:8082"
	_ = productsUrl
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", listProducts)
	router.HandleFunc("/products/{id}", showProduct)

	http.ListenAndServe(":8083", router)
}

func loadProducts() []Product {
	response, err := http.Get(productsUrl + "/products")
	if err != nil {
		fmt.Printf("The HTTP request failed with error: %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	var products Products
	json.Unmarshal(data, &products)

	return products.Products
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	products := loadProducts()
	t := template.Must(template.ParseFiles("./templates/catalog.html"))
	t.Execute(w, products)
}

func showProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := http.Get(productsUrl + "/products/" + vars["id"])
	if err != nil {
		fmt.Printf("The HTTP request failed with error: %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	err = json.Unmarshal(data, &product)
	if err != nil {
		fmt.Printf("An error has occurred while unmarshal json object: %s\n", err)
	}

	t := template.Must(template.ParseFiles("./templates/view.html"))
	t.Execute(w, product)
}