package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"backend/models"
)

// Home serves JSON response for home route.
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "hello world"})
}

// BooksList serves JSON response for books list route.
func BooksList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var books []models.Book
	db.Find(&books)
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// BooksCreate creates a new book record.
func BooksCreate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	book := models.Book{}
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&book)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// BooksGet serves JSON response containing a single book by ID.
func BooksGet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	book := models.Book{}
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record does not exist."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// BooksUpdate updates a book record.
func BooksUpdate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	book := models.Book{}
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record does not exist."})
		return
	}

	input := models.BookUpdater{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Model(&book).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// BooksDelete deletes a single book by ID.
func BooksDelete(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	book := models.Book{}
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record does not exist."})
	}

	db.Delete(&book)
	c.JSON(http.StatusAccepted, gin.H{"data": true})
}
