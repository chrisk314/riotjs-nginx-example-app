package models

import (
	"encoding/json"
	"testing"
)

func TestMarshallBook(t *testing.T) {
	jsonData := `{
		"title": "Childhood's End",
		"authors": "Arthur C. Clarke",
		"isbn": "9780330514019",
		"publisher": "Pan Macmillan",
		"publication_date": "2010-05-07T00:00:00Z",
		"language_code": "eng"
	}`
	book := Book{}
	if err := json.Unmarshal([]byte(jsonData), &book); err != nil {
		t.Errorf("Unmarshal failed with error: %s", err)
	}
	t.Log("book:", book)
}
