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

// Process endpoint
func processReceipt(c *gin.Context) {
	id := uuid.New()

	var receipt Receipt
	if err := c.BindJSON(&receipt); err != nil {
		return
	}

	// initialize points as 0
	pts := 0

	// every alphanumeric retailer name
	count := 0
	for _, c := range receipt.Retailer {
		if unicode.IsDigit(c) || unicode.IsLetter(c) {
			count += 1
		}
	}
	pts += count
	fmt.Println("  Alphanumeric:                    ", count)

	// total round
	t, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		fmt.Println("Error parsing float: ", err)
		return
	}
	if t == math.Trunc(t) {
		pts += 50
		fmt.Println("  Round Total:                     50")
	} else {
		fmt.Println("  Round Total:                      0")
	}

	// multiple of 0.25
	if 4*t == math.Trunc(4*t) {
		pts += 25
		fmt.Println("  Multiple of 0.25:                25")
	} else {
		fmt.Println("  Multiple of 0.25:                 0")
	}

	// 5 points for each pair
	pairs := 0
	pairs += int(len(receipt.Items) / 2)
	pts += pairs * 5
	fmt.Println("  5 points each pair:              ", pairs*5)

	// desc multiple of 3
	bonus := 0
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
			bonus += int(math.Ceil(price * 0.2))
		}
	}
	pts += bonus
	fmt.Println("  Desc length multiple of 3:       ", bonus)

	// odd day
	dayStr := receipt.PurchaseDate[len(receipt.PurchaseDate)-2:]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return
	}
	if day%2 == 1 {
		pts += 6
		fmt.Println("  Odd day:                          6")
	} else {
		fmt.Println("  Odd day:                          0")
	}

	// between 2pm & 4pm
	time := strings.Split(receipt.PurchaseTime, ":")
	hour, _ := strconv.Atoi(time[0])
	min, _ := strconv.Atoi(time[1])
	if (hour == 14 && min > 0) || (hour > 14 && hour < 16) {
		pts += 10
		fmt.Println("  After 2:00pm and before 4:00pm:  10")
	} else {
		fmt.Println("  After 2:00pm and before 4:00pm:   0")
	}

	// c.IndentedJSON(http.StatusCreated, receipt)
	points[id.String()] = pts
	fmt.Println("                                  ===")
	fmt.Println("  Points:                         ", pts)
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
	router.POST("/processReceipt", processReceipt)
	router.GET("/receipts/:id/points", getPoints)
	router.Run("localhost:8801")
}
