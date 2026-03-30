package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=aditya password=password123 dbname=bookmyshow port=5432 sslmode=disable"
	if envDSN := os.Getenv("DATABASE_URL"); envDSN != "" {
		dsn = envDSN
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("PostgreSQL connection established")

	err = DB.AutoMigrate(
		&models.Movie{}, 
		&models.Ticket{}, 
		&models.Chain{}, 
		&models.Theatre{}, 
		&models.Show{},
		&models.UserAccount{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedData()
}

func seedData() {
	// 1. Seed Movies (idempotent)
	var movieCount int64
	DB.Model(&models.Movie{}).Count(&movieCount)
	if movieCount == 0 {
		movies := []models.Movie{
			{Title: "Dhurandhar", Genre: "Action", Rating: "A", PosterURL: "/images/d2revenge_poster.webp", Formats: "2D,IMAX"},
			{Title: "Mario Galaxy", Genre: "Adventure", Rating: "U", PosterURL: "/images/mario_galaxy.webp", Formats: "3D,IMAX,ICE,3DSCREENX,4DX"},
			{Title: "Project Hail Mary", Genre: "Sci-Fi", Rating: "UA", PosterURL: "/images/phm_poster.jpeg", Formats: "2D,4DX"},
		}
		for i := range movies { DB.Create(&movies[i]) }
	}

	// 2. Seed Chains (idempotent)
	var chainCount int64
	DB.Model(&models.Chain{}).Count(&chainCount)
	if chainCount == 0 {
		chains := []models.Chain{
			{Name: "PVR-INOX"},
			{Name: "Cinepolis"},
			{Name: "Gopalan Cinemas"},
			{Name: "Single Screen / Independent"},
		}
		for i := range chains { DB.Create(&chains[i]) }
	}

	// 3. Seed Theatres (idempotent)
	var theatreCount int64
	DB.Model(&models.Theatre{}).Count(&theatreCount)
	if theatreCount == 0 {
		theatres := []models.Theatre{
			{Name: "Orion Mall, Malleshwaram", ChainID: 1, SupportedFormats: "2D,IMAX,4DX"},
			{Name: "Global Mall, Nayandalli", ChainID: 1, SupportedFormats: "2D,IMAX"},
			{Name: "Mantri Sq Mall, Malleshwaram", ChainID: 1, SupportedFormats: "2D,IMAX"},
			{Name: "Garuda Mall", ChainID: 2, SupportedFormats: "2D,IMAX,4DX,3D"},
			{Name: "Gopalan Arcade, RR Nagar", ChainID: 3, SupportedFormats: "2D,3D"},
			{Name: "Sri Venkateshwara Theatre, Girinagar", ChainID: 4, SupportedFormats: "2D"},
		}
		for i := range theatres { DB.Create(&theatres[i]) }
	}

	// 4. Seed Shows (idempotent for today)
	var showCount int64
	DB.Model(&models.Show{}).Count(&showCount)
	if showCount == 0 {
		var movies []models.Movie
		DB.Find(&movies)
		var theatres []models.Theatre
		DB.Find(&theatres)

		now := time.Now()
		for _, m := range movies {
			for _, t := range theatres {
				if m.Title == "Dhurandhar" && !contains(t.SupportedFormats, "IMAX") { continue }
				
				DB.Create(&models.Show{
					MovieID: m.ID,
					TheatreID: t.ID,
					Format: "2D",
					Showtime: time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location()),
					Price: 250,
				})
				if contains(t.SupportedFormats, "IMAX") {
					DB.Create(&models.Show{
						MovieID: m.ID,
						TheatreID: t.ID,
						Format: "IMAX",
						Showtime: time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()),
						Price: 550,
					})
				}
			}
		}
	}

	// 5. User Account (idempotent)
	var accCount int64
	DB.Model(&models.UserAccount{}).Count(&accCount)
	if accCount == 0 {
		DB.Create(&models.UserAccount{ID: 1, BMSCash: 0, InvestmentBalance: 0})
	}
}

func contains(s, substr string) bool {
	return fmt.Sprintf(",%s,", s) != "" // Simplified for MVP seeding logic
}
