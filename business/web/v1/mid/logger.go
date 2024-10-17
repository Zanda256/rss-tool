package mid

import (
	"fmt"
	"time"

	"github.com/Zanda256/rss-tool/foundation/logger"
	"github.com/Zanda256/rss-tool/foundation/web"
	"github.com/gin-gonic/gin"
)

func Logger(log *logger.Logger) web.Middleware {
	m := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			t := time.Now()
			r := c.Request
			path := c.Request.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}
			log.Info(c.Request.Context(), "request ended", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr)

			// before request

			c.Next()

			// after request
			latency := time.Since(t)

			log.Info(c.Request.Context(), "request completed", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr, "statuscode" /*v.StatusCode*/, "since", latency)
		}
	}
	return m
}
