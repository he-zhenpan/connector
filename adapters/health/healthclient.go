package health

import (
	"context"
	util "github.com/aldelo/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"fmt"
	"time"
)

type HealthClient struct {
	hcClient grpc_health_v1.HealthClient
}

func NewHealthClient(conn *grpc.ClientConn) (*HealthClient, error) {
	if conn == nil {
		return nil, fmt.Errorf("New Health Client Failed: %s", "gRPC Client Connection Nil")
	}

	return &HealthClient{
		hcClient: grpc_health_v1.NewHealthClient(conn),
	}, nil
}

func (h *HealthClient) Check(svcName string, timeoutDuration ...time.Duration) (grpc_health_v1.HealthCheckResponse_ServingStatus, error) {
	if h.hcClient == nil {
		return grpc_health_v1.HealthCheckResponse_UNKNOWN, fmt.Errorf("Health Check Failed: %s", "Health Check Client Nil")
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if len(timeoutDuration) > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeoutDuration[0])
	} else {
		ctx = context.Background()
	}

	in := &grpc_health_v1.HealthCheckRequest{}

	if util.LenTrim(svcName) > 0 {
		in.Service = svcName
	}

	if cancel != nil {
		defer cancel()
	}

	if resp, err := h.hcClient.Check(ctx, in); err != nil {
		return grpc_health_v1.HealthCheckResponse_UNKNOWN, fmt.Errorf("Health Check Failed: (Call Health Server Error) %s", err.Error())
	} else {
		return resp.Status, nil
	}
}