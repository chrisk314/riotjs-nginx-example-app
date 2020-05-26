package main

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"backend/models"
	"backend/utils"
)

const pageLimit int = 50

// JSONResp is a helper type alias to map[string]interface{} for returning JSON data.
type JSONResp map[string]interface{}

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

func isAllowedSortColumn(col string) bool {
	switch col {
	case "id", "title", "authors", "average_rating", "ratings", "reviews",
		"publication_date":
		return true
	}
	return false
}

type querySpec struct {
	field string
	op    string
	value string
}

type querySpecs []querySpec

// Fields returns an ordered list of all querySpec.field values in qss.
func (qss querySpecs) Fields() []string {
	fields := make([]string, 0)
	for _, qs := range qss {
		fields = append(fields, qs.field)
	}
	return fields
}

// FieldSet returns the set of all querySpec.field values in qss.
func (qss querySpecs) FieldSet() utils.StringSet {
	return utils.NewStringSet(qss.Fields())
}

func getQuerySpecFromTok(tok, sep string) (querySpec, error) {
	_x := strings.SplitN(tok, sep, 2)
	if len(_x) != 2 {
		return querySpec{}, fmt.Errorf("Invalid specifier: %s", _x)
	}
	return querySpec{field: _x[0], value: _x[1]}, nil
}

func getAfterSpecFromTok(tok string) (querySpec, error) {
	qs, err := getQuerySpecFromTok(tok, "=")
	if err != nil {
		return qs, err
	}
	if !(isAllowedSortColumn(qs.field)) {
		return querySpec{}, fmt.Errorf("Unsupported sort column: %s", qs.field)
	}
	return qs, nil
}

func getSortSpecFromTok(tok string) (querySpec, error) {
	qs, err := getQuerySpecFromTok(tok, ".")
	if err != nil {
		return qs, err
	}
	if !(isAllowedSortColumn(qs.field)) {
		return querySpec{}, fmt.Errorf("Unsupported sort column: %s", qs.field)
	}
	if !(qs.value == "asc" || qs.value == "desc") {
		return querySpec{}, fmt.Errorf("Invalid sort order: %s", qs.value)
	}
	return qs, nil
}

func getQuerySpecs(s string, qExtract func(string) (querySpec, error)) (querySpecs, error) {
	specs := querySpecs{}
	if s != "" {
		toks := strings.Split(s, ",")
		for _, tok := range toks {
			spec, err := qExtract(tok)
			if err != nil {
				return nil, err
			}
			specs = append(specs, spec)
		}
	}
	return specs, nil
}

func getAfterSpecs(s string) (qss querySpecs, err error) {
	qss, err = getQuerySpecs(s, getAfterSpecFromTok)
	if err != nil {
		return qss, err
	}
	return qss, nil
}

func getSortSpecs(s string) (qss querySpecs, err error) {
	qss, err = getQuerySpecs(s, getSortSpecFromTok)
	if err != nil {
		return qss, err
	}
	if len(qss) == 0 || !qss.FieldSet().Contains("id") {
		qss = append(qss, querySpec{field: "id", value: "asc"})
	}
	return qss, nil
}

var recognisedFilters = utils.NewStringSet(
	[]string{
		"average_rating",
		"language_code",
		"num_pages",
		"ratings",
		"reviews",
		"publication_date",
		"publisher",
	},
)

var recognisedFilterOps = map[string]string{
	"gt":  ">",
	"gte": ">=",
	"lt":  "<",
	"lte": "<=",
}

func getFilterSpecs(c echo.Context) (qss querySpecs, err error) {
	qss = querySpecs{}
	for field := range recognisedFilters {
		if _filters, ok := c.QueryParams()[field]; ok {
			for _, _filter := range _filters {
				_x := strings.SplitN(_filter, ":", 2)
				if len(_x) != 2 {
					return querySpecs{}, fmt.Errorf("Invalid specifier: %s", _x)
				}
				var op string
				if op, ok = recognisedFilterOps[_x[0]]; !ok {
					return querySpecs{}, fmt.Errorf("Invalid filter operation: %s", _x[0])
				}
				qss = append(qss, querySpec{field: field, op: op, value: _x[1]})
			}
		}
	}
	return qss, nil
}

func buildAfterQueryString(qss querySpecs, book models.Book) string {
	refBook := reflect.ValueOf(book)
	afterArgs := []string{}
	for _, ss := range qss {
		fieldVal := refBook.FieldByName(book.GetFieldByJSONTag(ss.field))
		fieldValURL := url.QueryEscape(fmt.Sprintf("%v", fieldVal))
		afterArgs = append(afterArgs, fmt.Sprintf("%s=%v", ss.field, fieldValURL))
	}
	return "after=" + strings.Join(afterArgs, ",")
}

func buildPaginationLinks(c echo.Context, books []models.Book, limit int,
	sortSpecs, filterSpecs querySpecs, prev bool) JSONResp {

	params := c.QueryParams()
	baseQs := fmt.Sprintf("limit=%d", limit)
	if _, ok := params["pretty"]; ok {
		baseQs = "pretty&" + baseQs
	}
	filters := make([]string, 0)
	for field := range recognisedFilters {
		if _filters, ok := params[field]; ok {
			for _, _filter := range _filters {
				filters = append(filters, fmt.Sprintf("%s=%s", field, _filter))
			}
		}
	}
	if len(filters) > 0 {
		baseQs = baseQs + "&" + strings.Join(filters, "&")
	}
	if v, ok := params["sort_by"]; ok {
		baseQs = baseQs + "&sort_by=" + v[0]
	}

	baseLink := fmt.Sprintf("%s/?%s", c.Path(), baseQs)
	selfLink, prevLink, nextLink := baseLink, "", ""
	_, hasAfter := params["after"]
	if !hasAfter {
		// first page of results
		afterQs := buildAfterQueryString(sortSpecs, books[len(books)-1])
		nextLink = selfLink + "&" + afterQs
	} else if hasAfter && !prev {
		// intermediate result going forwards
		afterQs := buildAfterQueryString(sortSpecs, books[0])
		selfLink = selfLink + "&" + afterQs
		prevLink = selfLink + "&previous"
		if len(books) == limit { // Fetched max records => assume more pages
			afterQs = buildAfterQueryString(sortSpecs, books[len(books)-1])
			nextLink = selfLink + "&" + afterQs
		}
	} else if hasAfter && prev {
		// intermediate result going backwards
		afterQs := buildAfterQueryString(sortSpecs, books[len(books)-1])
		nextLink = selfLink + "&" + afterQs
		if len(books) == limit { // Fetched max records => assume more pages
			afterQs = buildAfterQueryString(sortSpecs, books[0])
			selfLink = selfLink + "&" + afterQs
			prevLink = selfLink + "&previous"
		}
	}

	return JSONResp{"self": selfLink, "next": nextLink, "prev": prevLink}
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
	afterSpecs, err := getAfterSpecs(c.QueryParam("after"))
	if err != nil {
		return badRequest(c, err.Error())
	}
	sortSpecs, err := getSortSpecs(c.QueryParam("sort_by"))
	if err != nil {
		return badRequest(c, err.Error())
	}
	if len(sortSpecs) > 0 {
		if !sortSpecs.FieldSet().ContainsSet(afterSpecs.FieldSet()) {
			err = fmt.Errorf("incompatible `sort_by` and `after` specifications")
			return badRequest(c, err.Error())
		}
	}
	filterSpecs, err := getFilterSpecs(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	prev := false
	if _, ok := c.QueryParams()["previous"]; ok {
		prev = true
	}

	// Compose query
	q := db.Limit(limit)
	for _, ss := range sortSpecs {
		order := "desc"
		if (ss.value == "desc" && prev) || (ss.value == "asc" && !prev) {
			order = "asc"
		}
		q = q.Order(fmt.Sprintf("%s %s", ss.field, order))
	}
	for _, fs := range filterSpecs {
		q = q.Where(fmt.Sprintf("%s %s '%s'", fs.field, fs.op, fs.value))
	}
	// Seek query must be composed manually due to current limitations in Gorm
	// with complex AND/OR combinations.
	if len(afterSpecs) > 0 {
		var stmt string
		for i := 0; i < len(afterSpecs); i++ {
			op := "<"
			if (sortSpecs[i].value == "desc" && prev) || (sortSpecs[i].value == "asc" && !prev) {
				op = ">"
			}
			as := afterSpecs[i]
			orStmt := fmt.Sprintf("%s %s '%s'", as.field, op, as.value)
			for j := i - 1; j >= 0; j-- {
				as = afterSpecs[j]
				orStmt = fmt.Sprintf("%s = '%s' AND %s", as.field, as.value, orStmt)
			}
			if i == 0 {
				stmt = fmt.Sprintf("(%s)", orStmt)
			} else {
				stmt = fmt.Sprintf("%s OR (%s)", stmt, orStmt)
			}
		}
		q = q.Where(stmt)
	}

	// Execute query
	var books []models.Book
	q.Find(&books)
	if prev { // Put data in correct order if this was a request for previous page.
		for i, j := 0, len(books)-1; i < j; i, j = i+1, j-1 {
			books[i], books[j] = books[j], books[i]
		}
	}

	links := buildPaginationLinks(c, books, limit, sortSpecs, filterSpecs, prev)

	resp := JSONResp{"data": books, "count": len(books), "_links": links}
	return utils.JSON(c, http.StatusOK, resp, false)
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
