package controller

import (
	"encoding/json"
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"

	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/registry"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/usecase"
)

type PreferencesController struct {
	pu usecase.PreferencesUsecase
}

func NewPreferencesController(ctn *di.Container, r *routing.Router) *PreferencesController {
	pu := ctn.Resolve(registry.PreferencesUsecaseService).(usecase.PreferencesUsecase)

	c := &PreferencesController{
		pu: pu,
	}

	r.GetFunc("/preferences", c.HandleGet)
	r.PutFunc("/preferences", c.HandleUpdate)

	return c
}

func (pc *PreferencesController) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.NotFound(w, r)
		return
	}

	p, err := pc.pu.Get(r.Context(), userID)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&p)
}

func (pc *PreferencesController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data dto.Preferences
	_ = json.NewDecoder(r.Body).Decode(&data)

	err := pc.pu.Update(r.Context(), &data)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	}
}
