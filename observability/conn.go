package observability

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcConnection returns a configured connection to the collector on the specified endpoint. That connection can then
// be used for both metrics and tracing setup.
func GrpcConnection(ctx context.Context, endpoint string, tlsConfigs ...*tls.Config) (*grpc.ClientConn, error) {
	// Make sure there's only zero or one of the tls configs passed in.
	if len(tlsConfigs) > 1 {
		return nil, errors.New("more than 1 tls config passed in. Either one, or zero is accepted")
	}

	// by default set the grpc connection to be insecure, and only upgrade it if we have a tls connection.
	transportCredentials := insecure.NewCredentials()
	if len(tlsConfigs) == 1 {
		transportCredentials = credentials.NewTLS(tlsConfigs[0])
	}

	conn, err := grpc.DialContext(
		ctx,
		endpoint,
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.4,
				Jitter:     32,
				MaxDelay:   15 * time.Second,
			},
			MinConnectTimeout: 20 * time.Second,
		}),
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithBlock(),
	)
	if err == nil {
		return conn, nil
	}

	return nil, errors.Wrap(err, "grpc.DialContext after retrying with backoff")
}
