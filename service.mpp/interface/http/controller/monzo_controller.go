package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/util"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/monzo"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/plugin"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/registry"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/usecase"
)

type MonzoController struct {
	usecase usecase.UserUsecase
}

func NewMonzoController() *MonzoController {
	return &MonzoController{}
}

func (c *MonzoController) Apply(ctn *di.Container, m *http.ServeMux) {
	c.usecase = ctn.Resolve(registry.UserUsecaseService).(usecase.UserUsecase)

	m.HandleFunc("/monzo", c.Login)
	m.HandleFunc("/monzo/callback", c.LoginCallback)
	m.HandleFunc("/monzo/success", c.Success)
	m.HandleFunc("/monzo/hook", c.HandleEvent)
}

func (c *MonzoController) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Before new user")
	u, err := c.usecase.New()
	if err != nil {
		panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("Before monzo login")

	monzo.Login(w, r, u)
}

func (c *MonzoController) LoginCallback(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	code := vals.Get("code")
	state := vals.Get("state")

	log.Printf("[%s] /monzo/callback\n", r.Method)
	log.Printf("\t code: %s\n", code)
	log.Printf("\t state: %s\n", state)

	u, err := c.usecase.FindByStateToken(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ac, err := monzo.RequestAccessToken(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = c.usecase.SetAccessToken(u, ac)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/monzo/success", http.StatusPermanentRedirect)
}

func (c *MonzoController) Success(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("success"))
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
