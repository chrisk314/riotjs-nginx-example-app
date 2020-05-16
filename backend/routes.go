package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"backend/models"
)

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "hello world"})
}

func BooksList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var books []models.Book
	db.Find(&books)
	c.JSON(http.StatusOK, gin.H{"data": books})
}
