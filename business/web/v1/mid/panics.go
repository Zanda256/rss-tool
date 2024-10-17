package mid

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Zanda256/rss-tool/foundation/web"
	"github.com/gin-gonic/gin"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics() web.Middleware {
	m := func() gin.HandlerFunc {
		h := func(c *gin.Context) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err := fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
					c.AbortWithError(http.StatusBadGateway, err)
				}
			}()
		}

		return h
	}

	return m
}
