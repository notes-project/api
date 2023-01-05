package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexView interface {
	CreateOne(indexView mongo.IndexView, ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
}

type indexView struct{}

var indexViewInstace IndexView = indexView{}

func SetIndexViewInstance(c IndexView) {
	indexViewInstace = c
}

func GetIndexViewInstace() IndexView {
	return indexViewInstace
}

func (iv indexView) CreateOne(indexView mongo.IndexView, ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return indexView.CreateOne(ctx, model)
}
