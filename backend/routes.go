package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"backend/models"
)

const pageLimit int = 50

// MaxInt returns the greater of two ints
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MinInt returns the lesser of two ints
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// JSONResp is a helper type alias to map[string]interface{} for returning JSON data.
type JSONResp map[string]interface{}

func badRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, JSONResp{"error": msg})
}

func getRecordByParamKey(i interface{}, key string, db *gorm.DB, c echo.Context) error {
	return db.Where(fmt.Sprintf("%s = ?", key), c.Param(key)).First(i).Error
}

// Home serves JSON response for home route.
func Home(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, JSONResp{"data": "Hello world!"})
}

type sortSpec struct {
	orderBy string
	order   string
}

func isAllowedSortColumn(col string) bool {
	switch col {
	case "id", "title", "authors", "average_rating", "ratings", "reviews",
		"publication_date":
		return true
	}
	return false
}

func getSortSpecFromTok(tok string) (sortSpec, error) {
	_x := strings.Split(tok, ".")
	if len(_x) != 2 {
		return sortSpec{}, fmt.Errorf("Invalid sort specifier: %s", _x)
	}
	orderBy, order := _x[0], _x[1]
	if !(order == "asc" || order == "desc") {
		return sortSpec{}, fmt.Errorf("Invalid sort order: %s", order)
	}
	if !(isAllowedSortColumn(orderBy)) {
		return sortSpec{}, fmt.Errorf("Unsupported sort column: %s", orderBy)
	}
	return sortSpec{orderBy, order}, nil
}

func getSortSpecs(s string) ([]sortSpec, error) {
	sorters := make([]sortSpec, 0)
	toks := strings.Split(s, ",")
	for _, tok := range toks {
		ss, err := getSortSpecFromTok(tok)
		if err != nil {
			return nil, err
		}
		sorters = append(sorters, ss)
	}
	return sorters, nil
}

type afterSpec struct {
	field string
	value string
}

func getAfterSpecFromTok(tok string) (afterSpec, error) {
	_x := strings.Split(tok, "=")
	if len(_x) != 2 {
		return afterSpec{}, fmt.Errorf("Invalid specifier: %s", _x)
	}
	field, value := _x[0], _x[1]
	return afterSpec{field, value}, nil
}

func getAfterSpecs(s string) ([]afterSpec, error) {
	specs := make([]afterSpec, 0)
	toks := strings.Split(s, ",")
	for _, tok := range toks {
		spec, err := getAfterSpecFromTok(tok)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

// BooksList serves JSON response for books list route.
func BooksList(c echo.Context) (err error) {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	// Process query params
	limit := pageLimit
	if _limit := c.QueryParam("limit"); _limit != "" {
		limit, err = strconv.Atoi(_limit)
		if err != nil {
			return badRequest(c, fmt.Sprintf("Error parsing `limit`: %s", err.Error()))
		}
		limit = MinInt(MaxInt(1, limit), pageLimit)
	}
	afterSpecs := []afterSpec{}
	if _after := c.QueryParam("after"); _after != "" {
		afterSpecs, err = getAfterSpecs(_after)
		if err != nil {
			return badRequest(c, err.Error())
		}
	}
	sortSpecs := []sortSpec{}
	if _sortBy := c.QueryParam("sort_by"); _sortBy != "" {
		sortSpecs, err = getSortSpecs(_sortBy)
		if err != nil {
			return badRequest(c, err.Error())
		}
	}
	if len(afterSpecs) > 0 && len(afterSpecs) != len(sortSpecs) {
		err = fmt.Errorf("incompatible `sort_by` and `after` specifications")
		return badRequest(c, err.Error())
	}
	// Compose query
	q := db.Limit(limit)
	for _, ss := range sortSpecs {
		q = q.Order(fmt.Sprintf("%s %s", ss.orderBy, ss.order))
	}
	for i := 0; i < len(afterSpecs); i++ {
		var op string
		if sortSpecs[i].order == "asc" {
			op = ">"
		} else {
			op = "<"
		}
		as := afterSpecs[i]
		stmt := fmt.Sprintf("%s %s '%s'", as.field, op, as.value)
		for j := i - 1; j >= 0; j-- {
			as = afterSpecs[j]
			stmt = fmt.Sprintf("%s = '%s' AND %s", as.field, as.value, stmt)
		}
		if i == 0 {
			q = q.Where(stmt)
		} else {
			q = q.Or(stmt)
		}
	}
	// Execute query
	var books []models.Book
	q.Find(&books)
	return c.JSON(http.StatusOK, JSONResp{"data": books, "count": len(books)})
}

// BooksCreate creates a new book record.
func BooksCreate(c echo.Context) (err error) {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err = c.Bind(&book); err != nil { // Need to add input validation.
		return badRequest(c, err.Error())
	}
	db.Create(&book)
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksGet serves JSON response containing a single book by ID.
func BooksGet(c echo.Context) (err error) {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err = getRecordByParamKey(&book, "id", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksUpdate updates a book record.
func BooksUpdate(c echo.Context) (err error) {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err = getRecordByParamKey(&book, "id", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	input := models.BookUpdater{}
	if err = c.Bind(&input); err != nil {
		return badRequest(c, err.Error())
	}
	db.Model(&book).Updates(input)
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}

// BooksDelete deletes a single book by ID.
func BooksDelete(c echo.Context) (err error) {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err = getRecordByParamKey(&book, "id", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	db.Delete(&book)
	return c.JSON(http.StatusAccepted, JSONResp{"data": true})
}

// BooksGetByISBN serves JSON response containing a single book by ISBN.
func BooksGetByISBN(c echo.Context) (err error) {
	db := c.Get("db").(*gorm.DB) // TODO : How to ensure db not nil?
	book := models.Book{}
	if err = getRecordByParamKey(&book, "isbn", db, c); err != nil {
		return badRequest(c, err.Error())
	}
	return c.JSON(http.StatusOK, JSONResp{"data": book})
}
