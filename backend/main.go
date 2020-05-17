package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"backend/database"
)

func main() {
	// Init database
	db := database.GetDB()
	defer db.Close()

	// Echo instance
	router := echo.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(database.DBContext)

	// Routes
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
			books.GET("/isbn/:isbn", BooksGetByISBN)
		}
	}

	// Start server
	router.Logger.Fatal(router.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
