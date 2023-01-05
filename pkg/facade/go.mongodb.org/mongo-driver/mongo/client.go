package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client interface {
	Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error)
	Ping(client *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error
	Database(client *mongo.Client, name string, opts ...*options.DatabaseOptions) *mongo.Database
}

type client struct{}

var clientInstace Client = client{}

func SetClientInstance(c Client) {
	clientInstace = c
}

func GetClientInstace() Client {
	return clientInstace
}

func (c client) Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return mongo.Connect(ctx, opts...)
}

func (c client) Ping(client *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error {
	return client.Ping(ctx, rp)
}

func (c client) Database(client *mongo.Client, name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return client.Database(name, opts...)
}
