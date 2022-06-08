package database

import (
	"context"
	"github.com/uptrace/bun"
)

type DbCatalogStore struct {
	db  *bun.DB
	ctx context.Context
}

func NewDbCatalogStore(db *bun.DB, ctx context.Context) *DbCatalogStore {
	return &DbCatalogStore{
		db:  db,
		ctx: ctx,
	}
}
