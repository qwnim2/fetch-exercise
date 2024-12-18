package main

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var points = make(map[string]int)

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required,regexp=^[0-9]+\\.[0-9]{2}$"`
}

type Receipt struct {
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Total        string `json:"total" binding:"required"`
	Items        []Item `json:"items" binding:"required,min=1"`
}

func validateReceipt(r Receipt) error {

	retailerPattern := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !retailerPattern.MatchString(r.Retailer) {
		return fmt.Errorf("invalid retailer")
	}

	decimalRegex := regexp.MustCompile(`^[0-9]+\.[0-9]{2}$`)
	if !decimalRegex.MatchString(r.Total) {
		return fmt.Errorf("invalid total")
	}
	// check date
	if _, err := time.Parse("2006-01-02", r.PurchaseDate); err != nil {
		return fmt.Errorf("invalid date")
	}

	// check time
	if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
		return fmt.Errorf("invalid time")
	}

	// check shortDescription and price
	itemDescPattern := regexp.MustCompile(`^[\w\s\-]+$`)
	for _, i := range r.Items {
		if !itemDescPattern.MatchString(i.ShortDescription) {
			return fmt.Errorf("invalid description")
		}
		if !decimalRegex.MatchString(i.Price) {
			return fmt.Errorf("invalid price")
		}
	}

	return nil
}

// Process endpoint
func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"BadRequest": "The receipt is invalid."})
		return
	}

	if err := validateReceipt(receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"BadRequest": "The receipt is invalid."})
		return
	}

	id := uuid.New()
	pts := 0

	// Count alphanumeric chars in the retailer name
	alphanumericCount := 0
	for _, r := range receipt.Retailer {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			alphanumericCount += 1
		}
	}
	pts += alphanumericCount

	// if total is round
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		fmt.Println("Error parsing float: ", err)
		return
	}
	roundPoints := 0
	if totalFloat == math.Trunc(totalFloat) {
		roundPoints = 50
	}
	pts += roundPoints

	// if total is multiple of 0.25
	multipleOfQuarterPoints := 0
	if 4*totalFloat == math.Trunc(4*totalFloat) {
		multipleOfQuarterPoints = 25
	}
	pts += multipleOfQuarterPoints

	// 5 points for each pair
	pairCount := len(receipt.Items) / 2
	pairPoints := pairCount * 5
	pts += pairPoints

	// Desc len multiple of 3
	descBonus := 0
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		length := len(trimmedDesc)

		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return
		}

		if length%3 == 0 {
			descBonus += int(math.Ceil(price * 0.2))
		}
	}
	pts += descBonus

	// Odd day
	dayStr := receipt.PurchaseDate[len(receipt.PurchaseDate)-2:]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return
	}
	oddDayPoints := 0
	if day%2 == 1 {
		oddDayPoints = 6
	}
	pts += oddDayPoints

	// After 2 before 4
	timeParts := strings.Split(receipt.PurchaseTime, ":")
	hour, _ := strconv.Atoi(timeParts[0])
	min, _ := strconv.Atoi(timeParts[1])
	timePoints := 0
	if (hour == 14 && min > 0) || (hour > 14 && hour < 16) {
		timePoints = 10
	}
	pts += timePoints

	points[id.String()] = pts

	fmt.Println("Points Breakdown:")
	fmt.Printf("  Alphanumeric:                    %d\n", alphanumericCount)
	fmt.Printf("  Round Total:                     %d\n", roundPoints)
	fmt.Printf("  Multiple of 0.25:                %d\n", multipleOfQuarterPoints)
	fmt.Printf("  5 points each pair:              %d\n", pairPoints)
	fmt.Printf("  Desc length multiple of 3:       %d\n", descBonus)
	fmt.Printf("  Odd day:                         %d\n", oddDayPoints)
	fmt.Printf("  After 2:00pm and before 4:00pm:  %d\n", timePoints)
	fmt.Println("--------------------------------------")
	fmt.Printf("Total Points:                      %d\n", pts)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Getter endpoint
func getPoints(c *gin.Context) {
	id := c.Param("id")

	pts, exist := points[id]
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"NotFound": "No receipt found for that ID"})
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
