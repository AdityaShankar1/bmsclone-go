package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"backend/db"
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type SeatStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"` // "available", "locked", "booked"
}

func GetSeats(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		showIDStr := c.Query("show_id")
		
		showID, _ := strconv.Atoi(showIDStr)
		if showID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "show_id is required"})
			return
		}

		totalSeats := 20
		var seatStatuses []SeatStatus

		// 1. Get all booked seats from DB for this show
		var bookedTickets []models.Ticket
		db.DB.Where("show_id = ? AND status = ?", showID, "Confirmed").Find(&bookedTickets)
		bookedMap := make(map[string]bool)
		for _, t := range bookedTickets {
			bookedMap[t.SeatID] = true
		}

		// 2. Build seat list
		for i := 1; i <= totalSeats; i++ {
			seatID := fmt.Sprintf("A%d", i)
			status := "available"

			if bookedMap[seatID] {
				status = "booked"
			} else {
				// Check Redis lock
				lockKey := fmt.Sprintf("lock:show:%d:seat:%s", showID, seatID)
				exists, _ := rdb.Exists(ctx, lockKey).Result()
				if exists > 0 {
					status = "locked"
				}
			}

			seatStatuses = append(seatStatuses, SeatStatus{
				ID:     seatID,
				Status: status,
			})
		}

		c.JSON(http.StatusOK, gin.H{"seats": seatStatuses})
	}
}
