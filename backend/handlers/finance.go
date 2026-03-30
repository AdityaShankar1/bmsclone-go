package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend/db"
	"backend/models"
	"github.com/gin-gonic/gin"
)

func GetAccount(c *gin.Context) {
	var account models.UserAccount
	db.DB.First(&account, 1)
	c.JSON(http.StatusOK, account)
}

func Invest(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	var account models.UserAccount
	db.DB.First(&account, 1)
	account.InvestmentBalance += req.Amount
	db.DB.Save(&account)

	// Log Transaction
	db.DB.Create(&models.Transaction{
		UserID:      account.ID,
		Amount:      req.Amount,
		Type:        "Investment",
		Description: "Wallet Top-up",
		CreatedAt:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Investment successful", "balance": account.InvestmentBalance})
}

func CancelBooking(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, _ := strconv.Atoi(ticketIDStr)

	var ticket models.Ticket
	if err := db.DB.First(&ticket, ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	if ticket.Status == "Cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket already cancelled"})
		return
	}

	// 70% refund to BMSCash
	refundAmount := ticket.PricePaid * 0.7
	
	var account models.UserAccount
	db.DB.First(&account, 1)
	account.BMSCash += refundAmount
	db.DB.Save(&account)

	// Log Transaction
	db.DB.Create(&models.Transaction{
		UserID:      account.ID,
		Amount:      refundAmount,
		Type:        "Refund",
		Description: fmt.Sprintf("Refund for Ticket #%d", ticket.ID),
		CreatedAt:   time.Now(),
	})

	ticket.Status = "Cancelled"
	db.DB.Save(&ticket)

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled", "refund": refundAmount, "bms_cash": account.BMSCash})
}

func GetTransactions(c *gin.Context) {
	var txs []models.Transaction
	db.DB.Order("created_at desc").Find(&txs)
	c.JSON(http.StatusOK, gin.H{"transactions": txs})
}
