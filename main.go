package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"Bioskop/database"
	"Bioskop/routes"
	"time"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour, 
	}))

	database.ConnectDatabase()

	routes.AuthRoutes(r)

	r.Run(":8080") 
}
