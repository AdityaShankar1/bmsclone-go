package handlers

import (
	"net/http"

	"backend/db"
	"backend/models"
	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	db.DB.Find(&movies)
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

func GetBookings(c *gin.Context) {
	type BookingDetail struct {
		models.Ticket
		MovieTitle string `json:"movie_title"`
		Showtime   string `json:"showtime"`
		Format     string `json:"format"`
	}

	var results []BookingDetail
	db.DB.Table("tickets").
		Select("tickets.*, movies.title as movie_title, shows.showtime, shows.format").
		Joins("join shows on shows.id = tickets.show_id").
		Joins("join movies on movies.id = shows.movie_id").
		Order("tickets.created_at desc").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"bookings": results})
}
