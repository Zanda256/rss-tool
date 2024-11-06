package rssdb

import (
	"context"

	"github.com/Zanda256/rss-tool/business/data/searchIndex/es"
	"github.com/Zanda256/rss-tool/foundation/logger"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  *es.EsClient
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *es.EsClient) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s *Store) SaveEntry(ctx context.Context) error {
	return nil
}
