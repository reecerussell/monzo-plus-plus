package controller

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
)

type PreferencesController struct {
}

func NewPreferencesController(ctn *di.Container, m *http.ServeMux) *PreferencesController {
	return &PreferencesController{}
}
