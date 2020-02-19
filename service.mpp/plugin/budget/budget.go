package budget

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/libraries/util"

	"github.com/reecerussell/monzo-plus-plus/service.mpp/plugin"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/plugin/budget/proto"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/registry"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/usecase"

	"google.golang.org/grpc"
)

func init() {
	plugin.Register(&BudgetPlugin{})
}

// Environment variables.
var (
	BudgetRPCHost  = os.Getenv("BUDGET_RPC_HOST")
	BudgetHTTPHost = os.Getenv("BUDGET_HTTP_HOST")
)

const baseAPIURL = "/api/plugin/budget"

type BudgetPlugin struct {
	usecase usecase.UserUsecase
}

func (bp *BudgetPlugin) Name() string {
	return "budget"
}

func (bp *BudgetPlugin) Description() string {
	return "a daily budget calculator"
}

func (bp *BudgetPlugin) Build(ctn *di.Container) {
	usecase := ctn.Resolve(registry.UserUsecaseService).(usecase.UserUsecase)

	bp.usecase = usecase
}

func (bp *BudgetPlugin) Handler() http.Handler {
	remote, _ := url.Parse(BudgetHTTPHost)
	proxy := httputil.NewSingleHostReverseProxy(remote)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/api/plugin/budget", "", 1)

		proxy.ServeHTTP(w, r)
	})
}

func (bp *BudgetPlugin) TransactionCreated(ctx context.Context, t *monzo.Transaction) error {
	userID := ctx.Value(util.ContextKey("user_id")).(string)
	accessToken := ctx.Value(util.ContextKey("access_token")).(string)

	err := calculate(userID, t.AccountID, accessToken)
	if err != nil {
		return fmt.Errorf("t.c calc: %v", err)
	}

	return nil
}

// calculate makes an rpc call to the plugin server.
func calculate(userID, accountID, accessToken string) error {
	log.Printf("Insecurely dialing: %s\n", BudgetRPCHost)
	conn, err := grpc.Dial(BudgetRPCHost, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewBudgetServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := &proto.CalculateData{
		UserID:      userID,
		AccountID:   accountID,
		AccessToken: accessToken,
	}

	log.Printf("Making calculate call...\n")
	_, err = client.Calculate(ctx, data)
	if err != nil {
		log.Printf("\tFailed to calculate: %v\n", err)
		return fmt.Errorf("calc: %v", err)
	}

	log.Printf("Successfully made call.\n")

	return nil
}
