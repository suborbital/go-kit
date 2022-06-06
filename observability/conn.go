package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcConnRetries           = 5
	grpcConnBackoffBaseSecond = 1
)

// GrpcConnection returns a configured connection to the collector on the specified endpoint. That connection can then
// be used for both metrics and tracing setup.
func GrpcConnection(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	var err error
	var conn *grpc.ClientConn

	for i := 0; i < grpcConnRetries; i++ {
		conn, err = grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err == nil {
			return conn, nil
		}

		// grpcConnBackoffBaseSecond * (1 << i) is short form for base * 2 ^ i for that exponential backoff.
		time.Sleep(time.Second * grpcConnBackoffBaseSecond * (1 << i))
	}

	return nil, errors.Wrap(err, fmt.Sprintf("grpc.DialContext after %d retries", grpcConnRetries))
}
