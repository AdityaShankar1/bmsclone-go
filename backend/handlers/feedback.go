package handlers

import (
	"net/http"
	"strconv"

	"backend/db"
	"backend/models"
	"github.com/gin-gonic/gin"
)

func RateTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, _ := strconv.Atoi(ticketIDStr)

	var req struct {
		Rating int `json:"rating" binding:"required,min=1,max=5"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating (1-5)"})
		return
	}

	var ticket models.Ticket
	if err := db.DB.First(&ticket, ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	ticket.UserRating = req.Rating
	db.DB.Save(&ticket)

	c.JSON(http.StatusOK, gin.H{"message": "Rating saved!"})
}
