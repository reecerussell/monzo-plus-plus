package controller

import (
	"encoding/json"
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/dto"

	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/registry"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/usecase"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
)

type PreferencesController struct {
	pu usecase.PreferencesUsecase
}

func NewPreferencesController(ctn *di.Container, m *http.ServeMux) *PreferencesController {
	pu := ctn.Resolve(registry.PreferencesUsecaseService).(usecase.PreferencesUsecase)

	pc := &PreferencesController{
		pu: pu,
	}

	m.HandleFunc("/preferences", pc.router)

	return pc
}

func (pc *PreferencesController) router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pc.HandleGet(w, r)
		break
	case http.MethodPut:
		pc.HandleUpdate(w, r)
		break
	default:
		w.WriteHeader(http.StatusNotFound)
		break
	}
}

func (pc *PreferencesController) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.NotFound(w, r)
		return
	}

	p, err := pc.pu.Get(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(&p)
}

func (pc *PreferencesController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data dto.Preferences
	_ = json.NewDecoder(r.Body).Decode(&data)

	err := pc.pu.Update(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
