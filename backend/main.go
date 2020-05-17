package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/database"
)

func main() {
	db := database.GetDB()
	defer db.Close()

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(database.DBContext)

	api := router.Group("/api/v1")
	{
		api.GET("/", Home)
		api.GET("/books", BooksList)
		api.GET("/books/:id", BooksGetByID)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found."})
	})

	router.Run()
}
