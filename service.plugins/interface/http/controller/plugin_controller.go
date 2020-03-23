package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"

	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/registry"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/usecase"
)

// PluginController contains an array of HTTP handlers used to manage the Plugin domain.
type PluginController struct {
	u usecase.PluginUsecase
}

// NewPluginController is used to create a new instance of PluginController.
func NewPluginController(ctn *di.Container, r *mux.Router) *PluginController {
	u := ctn.Resolve(registry.ServicePluginUsecase).(usecase.PluginUsecase)
	c := &PluginController{u}

	r.HandleFunc("/", c.HandleGetList).Methods("GET")
	r.HandleFunc("/{id}", c.HandleGet).Methods("GET")
	r.HandleFunc("/", c.HandleCreate).Methods("POST")
	r.HandleFunc("/", c.HandleUpdate).Methods("PUT")
	r.HandleFunc("/{id}", c.HandleDelete).Methods("DELETE")

	return c
}

// HandleGet is a HTTP handler used to get an specific plugin.
func (c *PluginController) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	p, err := c.u.Get(id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&p)
}

// HandleGetList is a HTTP handler used to get a list of plugins.
func (c *PluginController) HandleGetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	term := r.URL.Query().Get("term")

	plugins, err := c.u.All(term)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&plugins)
}

// HandleCreate is a HTTP handler used to create a plugin record.
func (c *PluginController) HandleCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.CreatePlugin
	_ = json.NewDecoder(r.Body).Decode(&d)

	plugin, err := c.u.Create(&d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&plugin)
}

// HandleUpdate is a HTTP handler used to update a plugin.
func (c *PluginController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.UpdatePlugin
	_ = json.NewDecoder(r.Body).Decode(&d)

	err := c.u.Update(&d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleDelete is a HTTP handler used to delete individual plugins.
func (c *PluginController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	err := c.u.Delete(id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
