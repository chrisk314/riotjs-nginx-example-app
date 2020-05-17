package main

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"backend/models"
)

type JSONResp map[string]interface{}

// Home serves JSON response for home route.
func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, JSONResp{"data": "Hello world!"})
}

// BooksList serves JSON response for books list route.
func BooksList(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	var books []models.Book
	db.Find(&books)
	return c.JSON(http.StatusOK, JSONResp{"data": books})
}

// BooksCreate creates a new book record.
func BooksCreate(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, JSONResp{"error": err.Error()})
	}
	db.Create(&book)
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksGet serves JSON response containing a single book by ID.
func BooksGet(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		return c.JSON(http.StatusBadRequest, JSONResp{"error": "Record does not exist."})
	}
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksUpdate updates a book record.
func BooksUpdate(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		return c.JSON(http.StatusBadRequest, JSONResp{"error": "Record does not exist."})
	}
	input := models.BookUpdater{}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, JSONResp{"error": err.Error()})
	}
	db.Model(&book).Updates(input)
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksDelete deletes a single book by ID.
func BooksDelete(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		return c.JSON(http.StatusBadRequest, JSONResp{"error": "Record does not exist."})
	}
	db.Delete(&book)
	return c.JSON(http.StatusAccepted, JSONResp{"data": true})
}

// BooksGetByISBN serves JSON response containing a single book by ISBN.
func BooksGetByISBN(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := db.Where("isbn = ?", c.Param("isbn")).First(&book).Error; err != nil {
		return c.JSON(http.StatusBadRequest, JSONResp{"error": "Record does not exist."})
	}
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}
