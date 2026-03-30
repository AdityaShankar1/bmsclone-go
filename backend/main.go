package main

import (
	"context"
	"log"

	"backend/db"
	"backend/handlers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	db.InitDB()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Movie Routes
	r.GET("/movies", handlers.GetMovies)
	r.GET("/seats", handlers.GetSeats(rdb))
	r.POST("/book/:seatId", handlers.BookSeat(rdb))
	r.POST("/confirm-booking", handlers.ConfirmBooking(rdb))

	// Cinema Routes
	r.GET("/chains", handlers.GetChains)
	r.GET("/theatres", handlers.GetTheatres)
	r.GET("/shows", handlers.GetShows)

	// Finance & User Routes
	r.GET("/account", handlers.GetAccount)
	r.GET("/transactions", handlers.GetTransactions)
	r.POST("/invest", handlers.Invest)
	r.GET("/bookings", handlers.GetBookings)
	r.POST("/cancel-booking/:id", handlers.CancelBooking)
	r.POST("/rate-ticket/:id", handlers.RateTicket)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
