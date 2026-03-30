package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"backend/db"
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func BookSeat(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		seatID := c.Param("seatId")
		showID := c.Query("show_id")

		if showID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "show_id is required"})
			return
		}

		// 1. Check if already booked in DB
		var ticket models.Ticket
		if err := db.DB.Where("seat_id = ? AND show_id = ? AND status = ?", seatID, showID, "Confirmed").First(&ticket).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Seat already booked!"})
			return
		}

		// 2. Attempt to lock seat for 5 mins (SETNX)
		lockKey := fmt.Sprintf("lock:show:%s:seat:%s", showID, seatID)
		success, err := rdb.SetNX(ctx, lockKey, "user_123", 5*time.Minute).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error: " + err.Error()})
			return
		}

		if !success {
			c.JSON(http.StatusConflict, gin.H{"error": "Seat is currently being booked by someone else!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Seat locked. Proceed to pay.", "lock_id": lockKey})
	}
}

func ConfirmBooking(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		var req struct {
			SeatID string `json:"seat_id" binding:"required"`
			ShowID uint   `json:"show_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// 1. Verify Redis lock exists
		lockKey := fmt.Sprintf("lock:show:%d:seat:%s", req.ShowID, req.SeatID)
		_, err := rdb.Get(ctx, lockKey).Result()
		if err == redis.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Session expired. Please select seat again."})
			return
		}

		// 2. Calculate Final Price with Discounts
		var show models.Show
		db.DB.First(&show, req.ShowID)

		var account models.UserAccount
		db.DB.First(&account, 1)

		originalPrice := show.Price
		finalPrice := originalPrice

		// Inflation Protection / Frozen Rates (USP)
		if account.InvestmentBalance >= 500 {
			finalPrice = originalPrice * 0.8
		}

		// Apply BMSCash
		if account.BMSCash > 0 {
			if account.BMSCash >= finalPrice {
				account.BMSCash -= finalPrice
				finalPrice = 0
			} else {
				finalPrice -= account.BMSCash
				account.BMSCash = 0
			}
		}

		// 3. Insert ticket into DB
		ticket := models.Ticket{
			ShowID:              req.ShowID,
			UserID:              account.ID,
			SeatID:              req.SeatID,
			PricePaid:           finalPrice,
			Status:              "Confirmed",
			IsPaymentSuccessful: true,
			CreatedAt:           time.Now(),
		}

		if err := db.DB.Create(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
			return
		}

		// Update Account
		db.DB.Save(&account)

		// Log Transaction if price > 0
		if finalPrice > 0 {
			db.DB.Create(&models.Transaction{
				UserID:      account.ID,
				Amount:      finalPrice,
				Type:        "Payment",
				Description: fmt.Sprintf("Booking for Show #%d", req.ShowID),
				CreatedAt:   time.Now(),
			})
		}

		// 4. Delete Redis lock
		rdb.Del(ctx, lockKey)

		c.JSON(http.StatusOK, gin.H{
			"message":   "Booking confirmed!",
			"ticket_id": ticket.ID,
			"seat_id":   ticket.SeatID,
			"price":     ticket.PricePaid,
		})
	}
}
