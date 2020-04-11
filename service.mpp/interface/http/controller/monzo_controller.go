package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/routing"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/libraries/util"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/plugin"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/registry"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/usecase"
)

type MonzoController struct {
	usecase usecase.UserUsecase
}

func NewMonzoController(ctn *di.Container, r *routing.Router) *MonzoController {
	u := ctn.Resolve(registry.UserUsecaseService).(usecase.UserUsecase)

	c := &MonzoController{
		usecase: u,
	}

	r.PostFunc("/monzo/hook", c.HandleEvent)

	return c
}

func (c *MonzoController) HandleEvent(w http.ResponseWriter, r *http.Request) {
	var data monzo.TransactionEventWrapper
	_ = json.NewDecoder(r.Body).Decode(&data)

	ctx := context.Background()
	userID := r.URL.Query().Get("userId")
	accessToken, err := c.usecase.GetAccessToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ctx = context.WithValue(ctx, util.ContextKey("user_id"), userID)
	ctx = context.WithValue(ctx, util.ContextKey("access_token"), accessToken)

	log.Printf("[HOOK] New request for:\n\tUser: %s\n\tAccount: %s\n", userID, data.Data.AccountID)

	handlers := plugin.TransactionCreatedHandler()
	wg := &sync.WaitGroup{}

	for i, h := range handlers {
		go func(f plugin.TransactionFunc, idx int) {
			wg.Add(1)
			defer wg.Done()

			log.Printf("[%d]: Calling plugin method...\n", idx)

			err := f(ctx, &data.Data)
			if err != nil {
				log.Printf("[WEBHOOK] an event handler failed with error: %v\n", err)
			}

			log.Printf("\tsuccess.")
		}(h, i)
	}

	wg.Wait()
	w.WriteHeader(http.StatusOK)
}
