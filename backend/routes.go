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

// BooksGet serves JSON response containing a single book by ID.
func BooksGet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	book := models.Book{}
	if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record does not exist."})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": book})
	}
}
