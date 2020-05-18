package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"backend/models"
)

type JSONResp map[string]interface{}

func badRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, JSONResp{"error": msg})
}

func getRecordByParamKey(i interface{}, key string, db *gorm.DB, c echo.Context) error {
	return db.Where(fmt.Sprintf("%s = ?", key), c.Param(key)).First(i).Error
}

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
	if err := c.Bind(&book); err != nil { // Need to add input validation.
		return badRequest(c, err.Error())
	}
	db.Create(&book)
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksGet serves JSON response containing a single book by ID.
func BooksGet(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := getRecordByParamKey(&book, "id", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksUpdate updates a book record.
func BooksUpdate(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := getRecordByParamKey(&book, "id", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	input := models.BookUpdater{}
	if err := c.Bind(&input); err != nil {
		return badRequest(c, err.Error())
	}
	db.Model(&book).Updates(input)
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksDelete deletes a single book by ID.
func BooksDelete(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := getRecordByParamKey(&book, "id", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	db.Delete(&book)
	return c.JSON(http.StatusAccepted, JSONResp{"data": true})
}

// BooksGetByISBN serves JSON response containing a single book by ISBN.
func BooksGetByISBN(c echo.Context) error {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err := getRecordByParamKey(&book, "isbn", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}
