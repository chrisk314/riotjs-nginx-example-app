package database

import (
	"github.com/gin-gonic/gin"
)

// DBContext is Gin middleware which makes db available in context.
func DBContext(c *gin.Context) {
	c.Set("db", db)
	c.Next()
}
