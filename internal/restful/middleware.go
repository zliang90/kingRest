package restful

import (
	"github.com/zliang90/kingRest/internal/restful/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zliang90/kingRest/pkg/log"
	"github.com/zliang90/kingRest/pkg/util/uuid"
)

func handlerRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// api errors
				api.Failure(c, err)
			}
		}()

		c.Next()
	}
}

func handlerLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if raw != "" {
			path = path + "?" + raw
		}
		// access log
		log.Debugf("reqId: %s, %s \"%s %s %s\" %d %v %s %s\n",
			api.GetRequestId(c),
			c.ClientIP(),
			c.Request.Proto,
			c.Request.Method,
			path,
			c.Writer.Status(),
			time.Now().Sub(start),
			c.Request.UserAgent(),
			comment)
	}
}

func handlerRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// from http header
		requestId := c.GetHeader("Request-Id")
		if requestId == "" {
			// from query params
			requestId = c.DefaultQuery("Request-Id", "")
		}
		if !uuid.VerifyUUID(requestId) {
			requestId = uuid.New()
		}
		c.Set("Request-Id", requestId)

		c.Next()
	}
}
