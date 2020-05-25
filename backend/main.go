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
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(database.DBContext)

	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("", Home)
		books := api.Group("/books")
		{
			books.GET("", BooksList)
			books.POST("", BooksCreate)
			books.GET("/:id", BooksGet)
			books.PATCH("/:id", BooksUpdate)
			books.DELETE("/:id", BooksDelete)
			books.GET("/isbn/:isbn", BooksGetByISBN)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	sslCertPath := os.Getenv("SSL_CERT_PATH")
	sslKeyPath := os.Getenv("SSL_KEY_PATH")
	router.Logger.Fatal(router.StartTLS(fmt.Sprintf(":%s", port), sslCertPath, sslKeyPath))
}
