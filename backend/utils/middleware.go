package utils

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

const defaultIndent string = "    "

func encodeJSON(c echo.Context, code int, i interface{}, indent string, escapeJSON bool) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(code)
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	if !escapeJSON {
		enc.SetEscapeHTML(false)
	}
	return enc.Encode(i)
}

// JSON returns a JSON response for the passed interface.
// HTML safe escaping of characters can be switched off by passing
// espaceJSON=true.
// Note : this helper class exists because the JSON encoder in the golang
// standard library performs HTML safe escaping by default. This is a pretty
// big and inconvenient assumption about any downstream consumers of JSON.
func JSON(c echo.Context, code int, i interface{}, escapeJSON bool) (err error) {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; c.Echo().Debug || pretty {
		indent = defaultIndent
	}
	return encodeJSON(c, code, i, indent, escapeJSON)
}
