package controllers

import (
	"Bioskop/database"
	"Bioskop/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func BookMovie(c *gin.Context) {
	userID, userOk := c.Get("userID")
	role, roleOk := c.Get("role")

	if !userOk || !roleOk {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role or userID not found in context"})
		return
	}

	if role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "hanya role user yang dapat memesan"})
		return
	}

	id := c.Param("id")
	var movie models.Movies

	if err := database.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	pembayaran := c.PostForm("pembayaran")
	tiket, err := strconv.Atoi(c.PostForm("tiket"))
	if err != nil || tiket <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tiket must be a valid positive number"})
		return
	}

	if tiket > movie.JumlahTiket {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough tiket available"})
		return
	}

	movie.JumlahTiket -= tiket
	if err := database.DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie data"})
		return
	}

	booking := models.Booking{
		UserID:     userID.(uint),
		MovieID:    movie.ID,
		Pembayaran: pembayaran,
		Tiket:      tiket,
		MovieTitle: movie.Judul,
	}

	if err := database.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Tiket berhasil dipesan",
		"pembayaran": booking.Pembayaran,
		"tiket":      booking.Tiket,
		"movie_title": booking.MovieTitle,
		"created_at": booking.CreatedAt,
		"sisa_tiket": movie.JumlahTiket,
	})
}

func GetUserBookings(c *gin.Context) {
	userID, _ := c.Get("userID")

	var bookings []models.Booking
	if err := database.DB.Preload("Movie").Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	for i := range bookings {
		if bookings[i].Movie.ID != 0 { 
			status, err := calculateStatus(bookings[i].Movie.TanggalRilis)
			if err != nil {
				bookings[i].Movie.Status = "Tanggal tidak valid"
			} else {
				bookings[i].Movie.Status = status
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
	})
}


