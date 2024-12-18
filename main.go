package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "errors"
)

var points = make(map[string]int)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

func allReceipts(c *gin.Context) {
	// dictionary["hello"] = 2
	id := uuid.New()
	points[id.String()] = 1
	type Items []Item

	items := Items{
		Item{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
		Item{ShortDescription: "qwe - 5-oz", Price: "2.25"},
	}

	receipts := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        items,
	}

	fmt.Println("Endpoint Hit: All Receipts Endpoint")
	fmt.Println(points)
	c.IndentedJSON(http.StatusOK, receipts)
}

func main() {
	router := gin.Default()
	router.GET("/receipts", allReceipts)
	router.Run("localhost:8801")

}
