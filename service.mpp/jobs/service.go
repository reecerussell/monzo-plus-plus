package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/reecerussell/monzo-plus-plus/service.mpp/jobs/proto"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"google.golang.org/grpc"
)

// Environment variables.
var (
	JobsHost = os.Getenv("JOBS_HOST")
)

// Service is used to create jobs and add them to the queue.
type Service struct{}

// NewService returns a new instance of Service.
func NewService() *Service {
	return new(Service)
}

// Create is sued to push jobs to the job queue.
func (s *Service) Create(ctx context.Context, userID, pluginID string, data *monzo.Transaction) errors.Error {
	conn, err := grpc.Dial(JobsHost, grpc.WithInsecure())
	if err != nil {
		return errors.InternalError(fmt.Errorf("dial: %v", err))
	}
	defer conn.Close()

	client := proto.NewJobsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return errors.InternalError(fmt.Errorf("failed to marshal data: %v", err))
	}

	payload := &proto.PushRequest{
		UserID:   userID,
		PluginID: pluginID,
		Data:     string(dataBytes),
	}

	_, err = client.Push(ctx, payload)
	if err != nil {
		return errors.InternalError(fmt.Errorf("failed to push job: %v", err))
	}

	return nil
}
