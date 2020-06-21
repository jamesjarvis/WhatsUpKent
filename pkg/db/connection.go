package db

import (
	"context"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

// NewClient sets up a gRPC and returns a new dgraph connection
func NewClient(url string) (*dgo.Dgraph, error) {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	dialOpts := append([]grpc.DialOption{},
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	d, err := grpc.Dial(url, dialOpts...)

	if err != nil {
		return nil, err
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	), nil
}

// Setup initiates the schema into the database
func Setup(c *dgo.Dgraph) error {
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err := c.Alter(context.Background(), &api.Operation{
		Schema: Schema,
	})
	return err
}
