package rpc

import (
	"context"
	"errors"

	"github.com/reecerussell/monzo-plus-plus/service.job_queue/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/interface/rpc/proto"
	"github.com/reecerussell/monzo-plus-plus/service.job_queue/processing"
)

// JobsService is a service which handles incompling RPC calls.
type JobsService struct {
	jp processing.JobProcessor
}

// NewJobsService returns a new instance of JObsService.
func NewJobsService(jp processing.JobProcessor) *JobsService {
	return &JobsService{
		jp: jp,
	}
}

// Push handles RPCs to add a new job to the queue.
func (js *JobsService) Push(ctx context.Context, in *proto.PushRequest) (*proto.EmptyPushResponse, error) {
	j := model.NewJob(in.GetUserID(), in.GetPluginID(), in.GetData())

	err := js.jp.Push(j)
	if err != nil {
		return nil, errors.New(err.Text())
	}

	return &proto.EmptyPushResponse{}, nil
}
