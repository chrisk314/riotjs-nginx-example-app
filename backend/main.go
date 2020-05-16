package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/database"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	db := database.GetDB()
	defer db.Close()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	api := router.Group("/api/v1")
	{
		api.GET("/", Home)
		api.GET("/books", BooksList)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found."})
	})

	router.Run()
}
