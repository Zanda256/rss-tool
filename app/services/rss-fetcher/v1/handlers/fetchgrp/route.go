package fetchgrp

import (
	"net/http"

	"github.com/Zanda256/rss-tool/business/core/rss"
	"github.com/Zanda256/rss-tool/business/core/rss/stores/rssdb"
	"github.com/Zanda256/rss-tool/business/data/searchIndex/es"
	"github.com/Zanda256/rss-tool/foundation/logger"
	"github.com/Zanda256/rss-tool/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *logger.Logger
	DB  *es.EsClient
}

func Routes(app *web.App, cfg Config) {

	rssCore := rss.NewCore(rssdb.NewStore(cfg.Log, cfg.DB), cfg.Log)
	hdl := New(rssCore)

	app.Handle(http.MethodGet, "/hack", hdl.Hack)
}
