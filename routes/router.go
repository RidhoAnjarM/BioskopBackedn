package routes

import (
	"Bioskop/controllers"
	"Bioskop/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	router := r.Group("/api")
	{
		router.POST("/login", controllers.Login) //login
		router.POST("/register", controllers.Register) //register

		router.POST("/movies", controllers.CreateMovie) //create movie
		router.GET("/movies", controllers.GetAllMovies) //get semua movie
		router.GET("/movies/:id", controllers.GetMovieByID) //get byId movie
		router.PUT("/movies/:id", controllers.UpdateMovie) // update movie
		router.DELETE("/movies/:id", controllers.DeleteMovie) //delete movie

		router.POST("/movies/:id/pesan", middleware.AuthMiddleware(), controllers.BookMovie) //pesan tiket
		router.GET("/user/pesanan", middleware.AuthMiddleware(), controllers.GetUserBookings) //liat riwayat pemesanan
	}
}
