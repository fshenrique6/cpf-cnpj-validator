package status

import (
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	startTime   = time.Now()
	requestCount atomic.Int64
)

// CountRequests conta apenas as consultas feitas as rotas de /documents
func CountRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.FullPath(), "/documents") {
			requestCount.Add(1)
		}
		c.Next()
	}
}

func Handler (c *gin.Context) {
	uptime := time.Since(startTime).Round(time.Second).String()
	
	c.JSON(http.StatusOK, gin.H{
		"uptime":        uptime,
		"request_count": requestCount.Load(),
	})
}