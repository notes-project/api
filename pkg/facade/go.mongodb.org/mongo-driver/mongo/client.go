package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error)
}

type client struct{}

var clientInstace Client = client{}

func SetInstance(c Client) {
	clientInstace = c
}

func GetInstace() Client {
	return clientInstace
}

func (c client) Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return mongo.Connect(ctx, opts...)
}
