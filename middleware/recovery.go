package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogFailedRequests(c *gin.Context, recovered interface{}) {

	if err, ok := recovered.(string); ok {
		FailedRequestLogger(c)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
