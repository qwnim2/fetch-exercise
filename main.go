package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// func allReceipts(c *gin.Context) {
// 	// dictionary["hello"] = 2
// 	id := uuid.New()
// 	points[id.String()] = 1
// 	type Items []Item

// 	items := Items{
// 		Item{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
// 		Item{ShortDescription: "qwe - 5-oz", Price: "2.25"},
// 	}

// 	receipts := Receipt{
// 		Retailer:     "Target",
// 		PurchaseDate: "2022-01-02",
// 		PurchaseTime: "13:13",
// 		Total:        "1.25",
// 		Items:        items,
// 	}

// 	fmt.Println("Endpoint Hit: All Receipts Endpoint")
// 	c.IndentedJSON(http.StatusOK, receipts)
// }

// Process endpoint
func processReceipt(c *gin.Context) {
	id := uuid.New()

	var receipt Receipt
	if err := c.BindJSON(&receipt); err != nil {
		return
	}

	pts := 0
	// fmt.Println(receipt.Retailer)
	// fmt.Println(receipt.PurchaseDate)
	// fmt.Println(receipt.PurchaseTime)
	// fmt.Println(receipt.Total)
	// fmt.Println(receipt.Items)

	// every alphanumeric retailer name
	count := 0
	for _, c := range receipt.Retailer {
		if unicode.IsDigit(c) || unicode.IsLetter(c) {
			count += 1
		}
	}
	pts += count
	fmt.Println(pts)

	// total round
	t, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		fmt.Println("Error parsing float: ", err)
		return
	}
	// fmt.Println(int(t))
	if t == math.Trunc(t) {
		pts += 50
	}

	fmt.Println(pts)

	// multiple of 0.25
	if 4*t == math.Trunc(4*t) {
		pts += 25
	}
	fmt.Println(pts)

	// 5 points for each pair
	pairs := 0
	pairs += int(len(receipt.Items) / 2)
	pts += pairs * 5
	fmt.Println(pts)

	// desc multiple of 3
	for _, item := range receipt.Items {
		desc := item.ShortDescription
		priceStr := item.Price

		trimmedDesc := strings.TrimSpace(desc)
		length := len(trimmedDesc)

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return
		}

		if length%3 == 0 {
			pts += int(math.Ceil(price * 0.2))
		}
	}
	fmt.Println(pts)

	// odd day
	dayStr := receipt.PurchaseDate[len(receipt.PurchaseDate)-2:]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return
	}
	if day%2 == 1 {
		pts += 6
	}
	fmt.Println(pts)

	// between 2pm & 4pm
	time := strings.Split(receipt.PurchaseTime, ":")
	hour, _ := strconv.Atoi(time[0])
	min, _ := strconv.Atoi(time[1])
	if (hour == 14 && min > 0) || (hour > 14 && hour < 16) {
		pts += 10
	}
	fmt.Println(pts)
	// c.IndentedJSON(http.StatusCreated, receipt)
	points[id.String()] = pts
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Getter endpoint
func getPoints(c *gin.Context) {
	id := c.Param("id")

	pts, exist := points[id]
	if !exist {
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": pts})
}

func main() {
	router := gin.Default()
	// router.GET("/receipts", allReceipts)
	router.POST("/processReceipt", processReceipt)
	router.GET("/receipts/:id/points", getPoints)
	router.Run("localhost:8801")

}
