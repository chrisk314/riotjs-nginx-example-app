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
		books := api.Group("/books")
		{
			books.GET("/", BooksList)
			books.POST("/", BooksCreate)
			books.GET("/:id", BooksGet)
			books.PATCH("/:id", BooksUpdate)
			books.DELETE("/:id", BooksDelete)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found."})
	})

	router.Run()
}
