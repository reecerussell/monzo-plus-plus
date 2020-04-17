package processing

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/sse"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/processing/proto"

	"google.golang.org/grpc"
)

type JobProcessor interface {
	Start(ctx context.Context) errors.Error
	Push(j *model.Job) errors.Error
}

type jobProcessor struct {
	mu    *sync.RWMutex
	hosts map[string]string
	jobs  repository.JobRepository
	con   chan int
	mb    *sse.Broker

	WorkerLimit int
}

// NewJobProcessor returns a new instance of JobRepository.
func NewJobProcessor(repo repository.JobRepository, workerLimit int, mb *sse.Broker) JobProcessor {
	return &jobProcessor{
		mu:          &sync.RWMutex{},
		hosts:       make(map[string]string),
		jobs:        repo,
		con:         make(chan int, 1),
		mb:          mb,
		WorkerLimit: workerLimit,
	}
}

func (jp *jobProcessor) Start(ctx context.Context) errors.Error {
	var wg sync.WaitGroup

	for true {
		select {
		case <-ctx.Done():
			return nil
		default:
			break
		}

		jobs, err := jp.jobs.GetN(jp.WorkerLimit)
		if err != nil {
			return err
		}

		c := len(jobs)
		if c < 1 {
			<-jp.con
			continue
		}

		jp.mb.Notifier <- []byte(strconv.Itoa(c))

		for i := 0; i < c; i++ {
			go func(j *model.Job, idx int) {
				wg.Add(1)
				defer wg.Done()

				err := j.Execute(jp.process)
				if err != nil {
					log.Printf("[ERROR]: %v", err)
				}

				jp.mb.Notifier <- []byte(strconv.Itoa(c - idx))
			}(jobs[i], i)
		}

		wg.Wait()
	}

	return nil
}

func (jp *jobProcessor) Push(j *model.Job) errors.Error {
	err := jp.jobs.Add(j)
	if err != nil {
		return err
	}

	if len(jp.con) < 1 {
		jp.con <- 1
	}

	return nil
}

func (jp *jobProcessor) process(userID, pluginName, data string) errors.Error {
	host, hErr := jp.getPluginHost(pluginName)
	if hErr != nil {
		return hErr
	}

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return errors.InternalError(fmt.Errorf("dial: %v", err))
	}
	defer conn.Close()

	client := proto.NewPluginServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := &proto.SendRequest{
		UserID:   userID,
		JSONData: data,
	}

	_, err = client.Send(ctx, payload)
	if err != nil {
		return errors.InternalError(fmt.Errorf("send: %v", err))
	}

	return nil
}

func (jp *jobProcessor) getPluginHost(pluginName string) (string, errors.Error) {
	jp.mu.RLock()
	defer jp.mu.RUnlock()

	host, ok := jp.hosts[pluginName]
	if !ok {
		host, err := bootstrap.GetHost(pluginName)
		if err != nil {
			return "", errors.InternalError(err)
		}

		jp.hosts[pluginName] = host
	}

	return host, nil
}
