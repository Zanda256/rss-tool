package fetchgrp

import (
	"context"
	"net/http"

	"github.com/Zanda256/rss-tool/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// if n := rand.Intn(100) % 2; n == 0 {
	// 	panic("PANIC! This value will be stolen from the defer function in the panics middleware using closures and named return arguments")
	// 	//return response.NewError(errors.New("TRUST ERROR"), http.StatusBadRequest)
	// }
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK) //json.NewEncoder(w).Encode(status)
}
