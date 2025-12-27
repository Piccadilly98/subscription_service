package storage

import (
	"github.com/Piccadilly98/subscription_service/internal/storage/cache"
	"github.com/Piccadilly98/subscription_service/internal/storage/data_base"
)

type Storage struct {
	Db    *data_base.DataBase
	Cache *cache.Cache
}

func NewStorage(db *data_base.DataBase, cache *cache.Cache) *Storage {
	return &Storage{
		Db:    db,
		Cache: cache,
	}
}
