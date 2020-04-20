package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/jobs"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/provider"
)

// MonzoController is used to handle incoming HTTP requests from Monzo's webhooks.
type MonzoController struct {
	plugins *provider.UserPluginProvider
	jobs    *jobs.Service
}

// NewMonzoController returns a new instance of MonzoController, as well
// as registering a new route with the router.
func NewMonzoController(r *routing.Router) *MonzoController {
	c := &MonzoController{
		plugins: provider.NewUserPluginProvider(),
		jobs:    jobs.NewService(),
	}

	r.PostFunc("/monzo/hook", c.HandleEvent)

	return c
}

// HandleEvent is used to handle POST requests from Monzo's wehbook.
func (c *MonzoController) HandleEvent(w http.ResponseWriter, r *http.Request) {
	var data monzo.TransactionEvent
	_ = json.NewDecoder(r.Body).Decode(&data)
	userID := r.URL.Query().Get("userId")

	pluginIDs, err := c.plugins.GetUserPlugins(userID)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	var wg sync.WaitGroup

	log.Printf("No. of Plugins: %d\n", len(pluginIDs))
	for _, id := range pluginIDs {
		log.Printf("\t- %s\n", id)
	}

	for _, pluginID := range pluginIDs {
		go func(id string) {
			wg.Add(1)
			defer wg.Done()

			err := c.jobs.Create(r.Context(), userID, id, data.Data)
			if err != nil {
				log.Printf("[ERROR]: WEBHOOK ERROR\n\tUser Id: %s\n\tPlugin Id: %s\n\tError: %s", userID, id, err.Text())
			}
		}(pluginID)
	}

	wg.Wait()
}
