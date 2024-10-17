package web

import "github.com/gin-gonic/gin"

type Middleware func() gin.HandlerFunc

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func (a *App) wrapMiddleware(mw []Middleware, handler Handler) Handler {

	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := a.mw[i]
		if mwFunc != nil {
			a.Use(mwFunc())
		}
	}

	return handler
}
