package controller

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
)

// UI is used to serve UI content.
func UI(r *routing.Router) {
	fs := http.FileServer(http.Dir("ui"))
	r.Get("/", fs)
}
