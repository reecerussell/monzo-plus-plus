package processing

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/processing/proto"

	"google.golang.org/grpc"
)

// Queue is a first-in, first-out queue used for executing and proccessing Jobs.
type Queue struct {
	jobs          repository.JobRepository
	hold          chan bool
	wg            sync.WaitGroup
	pool          chan *Worker
	workers       []*Worker
	internalQueue chan *model.Job
}

// NewQueue returns a new instance of Queue.
func NewQueue(repo repository.JobRepository, workerLimit int) *Queue {
	q := &Queue{
		jobs:          repo,
		hold:          make(chan bool, 1),
		wg:            sync.WaitGroup{},
		pool:          make(chan *Worker, workerLimit),
		workers:       make([]*Worker, workerLimit),
		internalQueue: make(chan *model.Job, workerLimit),
	}

	hosts := NewHostProvider()

	for i := range q.workers {
		q.workers[i] = &Worker{
			jobs:  repo,
			hosts: hosts,
		}
	}

	return q
}

// Start begins the queue and pooling.
func (q *Queue) Start() {
	for _, w := range q.workers {
		q.pool <- w
	}

	go q.dispatch()

	for {
		w := <-q.pool
		go func() {
			j := <-q.internalQueue
			defer q.wg.Done()
			w.Process(j)
			q.pool <- w
		}()
	}
}

func (q *Queue) dispatch() {
	for {
		q.wg.Wait()
		jobs, err := q.jobs.GetN(len(q.workers))
		if err != nil {
			log.Printf("\t[ERROR]: failed to get jobs: %s\n", err.Text())
			break
		}
		log.Printf("\tDone.\n")

		if c := len(jobs); c < 1 {
			<-q.hold
			continue
		} else {
			q.wg.Add(c)
		}

		for _, j := range jobs {
			q.internalQueue <- j
		}
	}
}

// Push adds a Job database record and triggers the queue to start processing.
func (q *Queue) Push(j *model.Job) errors.Error {
	log.Printf("Job received, adding to database...")
	err := q.jobs.Add(j)
	if err != nil {
		return err
	}
	log.Printf("\tDone.\n")

	if len(q.hold) < 1 {
		q.hold <- true
		log.Printf("Notified queue.")
	}

	return nil
}

// Worker is used to individually handle the proccessing of a job.
type Worker struct {
	hosts *HostProvider
	jobs  repository.JobRepository
}

// Process is the worker's entrypoint when it's used. It executed and
// updates a Job, depending on the result of the execution.
func (w *Worker) Process(j *model.Job) {
	err := j.Execute(w.internalProcessor)
	if err != nil {
		log.Printf("[JOB:%d][ERROR]: an error occured while executing the job: %s\n", j.GetID(), err.Text())
	}

	err = w.jobs.Update(j)
	if err != nil {
		log.Printf("[JOB:%d][ERROR]: an error occured while updating the job: %s\n", j.GetID(), err.Text())
	}
}

// internalProcessor is a ProcessFunc used by a Job to process itself.
func (w *Worker) internalProcessor(userID, accountID, pluginName, data string) errors.Error {
	var host, accessToken string
	if h, err := w.hosts.Get(pluginName); err == nil {
		host = h
	} else {
		return err
	}

	if ac, err := getAccessToken(userID); err == nil {
		accessToken = ac
	} else {
		return err
	}

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return errors.InternalError(fmt.Errorf("dial: %v", err))
	}
	defer conn.Close()

	client := proto.NewPluginServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload := &proto.SendRequest{
		UserID:      userID,
		AccountID:   accountID,
		AccessToken: accessToken,
		JSONData:    data,
	}

	_, err = client.Send(ctx, payload)
	if err != nil {
		return errors.InternalError(fmt.Errorf("send: %v", err))
	}

	return nil
}

// HostProvider is used to provider and cache plugin hostnames.
type HostProvider struct {
	mu    sync.RWMutex
	hosts map[string]string
}

// NewHostProvider returns a new instance of HostProvider.
func NewHostProvider() *HostProvider {
	return &HostProvider{
		mu:    sync.RWMutex{},
		hosts: make(map[string]string),
	}
}

// Get returns the hostname for the given plugin. If the hostname
// is not in the provider's internal cache, a request will be made to
// the registry for it, via RPC.
func (p *HostProvider) Get(pluginName string) (string, errors.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	host, ok := p.hosts[pluginName]
	if !ok || host == "" {
		host, err := bootstrap.GetHost(pluginName)
		if err != nil {
			return "", errors.InternalError(err)
		}

		if host == "" {
			return "", errors.InternalError(fmt.Errorf("host not found for %s", pluginName))
		}

		p.hosts[pluginName] = fmt.Sprintf("%s:8080", host)
	}

	return host, nil
}

func getAccessToken(userID string) (string, errors.Error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_RPC_HOST"), grpc.WithInsecure())
	if err != nil {
		return "", errors.InternalError(fmt.Errorf("dial: %v", err))
	}
	defer conn.Close()

	client := proto.NewPermissionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := &proto.AccessTokenRequest{
		UserID: userID,
	}

	res, err := client.GetMonzoAccessToken(ctx, payload)
	if err != nil {
		return "", errors.InternalError(fmt.Errorf("send: %v", err))
	}

	return res.GetAccessToken(), nil
}
