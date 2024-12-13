package controllers

import (
	"Bioskop/database"
	"Bioskop/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func calculateStatus(tanggalRilis string) (string, error) {
	tanggal, err := time.Parse("2006-01-02", tanggalRilis)
	if err != nil {
		return "", err
	}

	today := time.Now()

	today = today.Local().Truncate(24 * time.Hour)

	if tanggal.After(today) {
		return "Akan Datang", nil
	} else if tanggal.Equal(today) {
		return "Sedang Berlangsung", nil
	} else {
		return "Selesai", nil
	}
}

func CreateMovie(c *gin.Context) {
	var movie models.Movies

	file, err := c.FormFile("photo")
	if err != nil {
		movie.Photo = "" 
	} else {
		uploadPath := fmt.Sprintf("./uploads/%s", file.Filename)

		if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("./uploads", os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
				return
			}
		}

		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed upload photo"})
			return
		}

		movie.Photo = fmt.Sprintf("/uploads/%s", file.Filename)
	}

	movie.Judul = c.PostForm("judul")
	if movie.Judul == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul harus diisi"})
		return
	}

	movie.Deskripsi = c.PostForm("deskripsi")
	movie.TanggalRilis = c.PostForm("tanggal_rilis")
	status, err := calculateStatus(movie.TanggalRilis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal rilis tidak valid"})
		return
	}
	movie.Status = status

	movie.Trailer = c.PostForm("trailer_url")
	movie.Genre = c.PostForm("genre")

	harga, err := strconv.Atoi(c.PostForm("harga"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harga must be a valid number"})
		return
	}
	movie.Harga = harga

	jumlahTiket, err := strconv.Atoi(c.PostForm("jumlah_tiket"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah tiket must be a valid number"})
		return
	}
	movie.JumlahTiket = jumlahTiket

	if err := database.DB.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal save movie ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie berhasil dibuat",
		"data":    movie,
	})
}

func GetAllMovies(c *gin.Context) {
	var movies []models.Movies

	if err := database.DB.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	for i := range movies {
		status, err := calculateStatus(movies[i].TanggalRilis)
		if err != nil {
			movies[i].Status = "Tanggal tidak valid"
		} else {
			movies[i].Status = status
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
	})
}

func GetMovieByID(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movies

	if err := database.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie tidak ditemukan"})
		return
	}

	status, err := calculateStatus(movie.TanggalRilis)
	if err != nil {
		movie.Status = "Tanggal tidak valid"
	} else {
		movie.Status = status
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movie,
	})
}


func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movies

	if err := database.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie tidak ditemukan"})
		return
	}

	file, err := c.FormFile("photo")
	if err == nil {
		uploadPath := fmt.Sprintf("./uploads/%s", file.Filename)
		if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("./uploads", os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
				return
			}
		}
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload photo"})
			return
		}
		movie.Photo = fmt.Sprintf("/uploads/%s", file.Filename)
	}

	movie.Judul = c.PostForm("judul")
	movie.Deskripsi = c.PostForm("deskripsi")
	movie.TanggalRilis = c.PostForm("tanggal_rilis")
	movie.Trailer = c.PostForm("trailer_url")
	movie.Genre = c.PostForm("genre")

	harga, err := strconv.Atoi(c.PostForm("harga"))
	if err == nil {
		movie.Harga = harga
	}

	jumlahTiket, err := strconv.Atoi(c.PostForm("jumlah_tiket"))
	if err == nil {
		movie.JumlahTiket = jumlahTiket
	}

	if err := database.DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate movie"})
		return
	}

	status, err := calculateStatus(movie.TanggalRilis)
	if err != nil {
		movie.Status = "Tanggal tidak valid"
	} else {
		movie.Status = status
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie berhasil diupdate",
		"data": gin.H{
			"id":           movie.ID,
			"judul":        movie.Judul,
			"deskripsi":    movie.Deskripsi,
			"tanggal_rilis": movie.TanggalRilis,
			"photo":        movie.Photo,
			"harga":        movie.Harga,
			"trailer_url":  movie.Trailer,
			"jumlah_tiket": movie.JumlahTiket,
			"genre":        movie.Genre,
			"status":       movie.Status,
		},
	})
}

func DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movies

	if err := database.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie tidak ditemukan"})
		return
	}

	if err := database.DB.Delete(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie Berhasil dihapus",
	})
}
