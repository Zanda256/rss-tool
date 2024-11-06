package v1

import (
	"net/http"
	"os"

	"github.com/Zanda256/rss-tool/business/data/searchIndex/es"
	"github.com/Zanda256/rss-tool/business/web/v1/mid"
	"github.com/Zanda256/rss-tool/foundation/logger"
	"github.com/Zanda256/rss-tool/foundation/web"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	WebClient *http.Client
	Build     string
	Shutdown  chan os.Signal
	Log       *logger.Logger
	Db        *es.EsClient
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg APIMuxConfig)
}

// APIMux constructs a http.Handler with all application routes defined.
// Note: The order in which the middle wares is passed is outer to inner,
// i.e the inner most middleware is on the extreme right
func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) http.Handler {
	mux := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log) /*mid.Errors(cfg.Log)*/, mid.Panics())

	routeAdder.Add(mux, cfg)

	return mux
}
