package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Book models database Book table. Represents a book.
type Book struct {
	gorm.Model
	Title           string    `json:"title" binding:"required"`
	Authors         string    `json:"authors" binding:"required"`
	AverageRating   string    `json:"average_rating"`
	ISBN            string    `json:"isbn" binding:"required"`
	ISBN13          string    `json:"isbn_13"`
	LanguageCode    string    `json:"language_code" binding:"required"`
	NumPages        int       `json:"num_pages"`
	Ratings         int       `json:"ratings"`
	Reviews         int       `json:"reviews"`
	PublicationDate time.Time `json:"publication_date" binding:"required"`
	Publisher       string    `json:"publisher" binding:"required"`
}

// BookUpdater facilitates updating field(s) in a Book record.
type BookUpdater struct {
	Title           string    `json:"title"`
	Authors         string    `json:"authors"`
	AverageRating   string    `json:"average_rating"`
	ISBN            string    `json:"isbn"`
	ISBN13          string    `json:"isbn_13"`
	LanguageCode    string    `json:"language_code"`
	NumPages        int       `json:"num_pages"`
	Ratings         int       `json:"ratings"`
	Reviews         int       `json:"reviews"`
	PublicationDate time.Time `json:"publication_date"`
	Publisher       string    `json:"publisher"`
}
