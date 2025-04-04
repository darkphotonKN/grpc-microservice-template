package discovery

import (
	"context"
	"errors"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {

	// discover the other services
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	length := rand.Intn(len(addrs))

	if length == 0 {
		return nil, errors.New("There are no services to discover now.")
	}

	return grpc.Dial(addrs[length], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
