package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"

	"github.com/reecerussell/monzo-plus-plus/libraries/sse"

	"github.com/reecerussell/monzo-plus-plus/service.job_queue/interface/http"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/interface/persistence"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/interface/rpc"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/processing"
)

// DefaultWorkerLimit is used if the WORKER_LIMIT environment
// variable was not set.
const DefaultWorkerLimit = 3

func main() {
	mb := sse.NewBroker()

	repo := persistence.NewJobRepository()
	processor := processing.NewJobProcessor(repo, getLimit(), mb)
	pctx, cancelProcessing := context.WithCancel(context.Background())
	go processor.Start(pctx)

	web := http.Build(mb)
	go web.Serve()

	rpc := rpc.Build(processor)
	go rpc.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	web.Shutdown(bootstrap.ShutdownGraceful)
	rpc.Shutdown(bootstrap.ShutdownGraceful)
	cancelProcessing()
}

func getLimit() int {
	if v := os.Getenv("WORKER_LIMIT"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			return DefaultWorkerLimit
		}

		return l
	}

	return DefaultWorkerLimit
}
