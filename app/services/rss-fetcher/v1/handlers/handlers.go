package handlers

import (
	"github.com/Zanda256/rss-tool/app/services/rss-fetcher/v1/handlers/fetchgrp"
	v1 "github.com/Zanda256/rss-tool/business/web/v1"
	"github.com/Zanda256/rss-tool/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	fetchgrp.Routes(app)
}
