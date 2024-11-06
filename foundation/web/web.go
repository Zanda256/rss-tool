package web

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Handler handles http requests within this web package
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*gin.Engine
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(channel chan os.Signal, mw ...Middleware) *App {
	return &App{
		Engine:   gin.New(),
		shutdown: channel,
		mw:       mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method, path string, handler Handler, mw ...Middleware) { // This function overrides the mux's Handle method
	a.wrapMiddleware(mw, handler)

	h := func(c *gin.Context) {
		if err := handler(c.Request.Context(), c.Writer, c.Request); err != nil {
			c.Error(err)
		}
	}

	handlerChain := make([]gin.HandlerFunc, 0)
	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := a.mw[i]
		if mwFunc != nil {
			handlerChain = append(handlerChain, mwFunc())
		}
	}

	a.Engine.Handle(method, path, append(handlerChain, h)...)
}
