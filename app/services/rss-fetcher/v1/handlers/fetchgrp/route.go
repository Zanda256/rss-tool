package fetchgrp

import (
	"net/http"

	"github.com/Zanda256/rss-tool/foundation/web"
)

func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
