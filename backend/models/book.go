package models

import (
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var bookTagToField = make(map[string]string)

// Book models database Book table. Represents a book.
type Book struct {
	gorm.Model
	Title           string    `json:"title" binding:"required"`
	Authors         string    `json:"authors" binding:"required"`
	AverageRating   string    `json:"average_rating"`
	ISBN            string    `json:"isbn" binding:"required" gorm:"index"`
	ISBN13          string    `json:"isbn_13"`
	LanguageCode    string    `json:"language_code" binding:"required"`
	NumPages        int       `json:"num_pages"`
	Ratings         int       `json:"ratings"`
	Reviews         int       `json:"reviews"`
	PublicationDate time.Time `json:"publication_date" binding:"required"`
	Publisher       string    `json:"publisher" binding:"required"`
}

func buildBookTagToFieldMap() {
	ref := reflect.TypeOf(Book{})
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		bookTagToField[tag] = field.Name
	}
	bookTagToField["id"] = "ID"
}

// GetFieldByJSONTag returns struct field name with corresponding json tag
func (book *Book) GetFieldByJSONTag(tag string) string {
	if len(bookTagToField) == 0 {
		buildBookTagToFieldMap()
	}
	return bookTagToField[tag]
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
