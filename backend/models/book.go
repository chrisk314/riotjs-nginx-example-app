package models

type Book struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Title         string `json:"title"`
	Authors       string `json:"authors"`
	AverageRating string `json:"average_rating"`
	ISBN          string `json:"isbn"`
	ISBN13        string `json:"isbn_13"`
	LanguageCode  string `json:"language_code"`
	NumPages      int    `json:"num_pages"`
	Ratings       int    `json:"ratings"`
	Reviews       int    `json:"reviews"`
}
