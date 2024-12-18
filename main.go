package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	// "github.com/gin-gonic/gin"
	// "errors"
)

type Item struct{
	ShortDescription string `json:"shortDescription"`
	Price string `json:"price"`
}

type Receipt struct{
	Retailer string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total string `json:"total"`
	Items []Item `json:"items"`
}


func allReceipts(w http.ResponseWriter, r *http.Request) {

	type Items []Item

	items := Items{
		Item{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
		Item{ShortDescription: "qwe - 5-oz", Price: "2.25"},
	}

	receipts := Receipt{
		Retailer: "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total: "1.25",
		Items: items,
	}

	fmt.Println("Endpoint Hit: All Receipts Endpoint")
	json.NewEncoder(w).Encode(receipts)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/receipts", allReceipts)
	log.Fatal(http.ListenAndServe(":8801", nil))
}

func main(){
	handleRequests()
}