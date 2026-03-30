package handlers

import (
	"net/http"
	"strconv"

	"backend/db"
	"backend/models"
	"github.com/gin-gonic/gin"
)

func GetChains(c *gin.Context) {
	var chains []models.Chain
	db.DB.Find(&chains)
	c.JSON(http.StatusOK, gin.H{"chains": chains})
}

func GetTheatres(c *gin.Context) {
	chainIDStr := c.Query("chain_id")
	movieIDStr := c.Query("movie_id")
	format := c.Query("format")

	var theatres []models.Theatre
	query := db.DB.Model(&models.Theatre{})

	if chainIDStr != "" {
		chainID, _ := strconv.Atoi(chainIDStr)
		query = query.Where("chain_id = ?", chainID)
	}

	if movieIDStr != "" {
		movieID, _ := strconv.Atoi(movieIDStr)
		// Filter theatres showing this movie
		query = query.Joins("JOIN shows ON shows.theatre_id = theatres.id").Where("shows.movie_id = ?", movieID)
		if format != "" {
			query = query.Where("shows.format = ?", format)
		}
	}

	query.Distinct().Find(&theatres)
	c.JSON(http.StatusOK, gin.H{"theatres": theatres})
}

func GetShows(c *gin.Context) {
	movieIDStr := c.Query("movie_id")
	theatreIDStr := c.Query("theatre_id")

	var shows []models.Show
	query := db.DB.Model(&models.Show{})

	if movieIDStr != "" {
		movieID, _ := strconv.Atoi(movieIDStr)
		query = query.Where("movie_id = ?", movieID)
	}
	if theatreIDStr != "" {
		theatreID, _ := strconv.Atoi(theatreIDStr)
		query = query.Where("theatre_id = ?", theatreID)
	}

	query.Find(&shows)
	c.JSON(http.StatusOK, gin.H{"shows": shows})
}
