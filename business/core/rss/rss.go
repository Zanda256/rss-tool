package rss

import (
	"context"

	"github.com/Zanda256/rss-tool/foundation/logger"
)

type Storer interface {
	SaveEntry(context.Context) error
}

type Core struct {
	store Storer
	log   *logger.Logger
}

func NewCore(store Storer, log *logger.Logger) *Core {
	return &Core{
		store: store,
		log:   log,
	}
}
