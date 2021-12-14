package redis

import (
	"context"

	"github.com/quantonganh/shortener"
)

type urlCache struct {
	db *DB
}

func NewURLCache(db *DB) shortener.URLCache {
	return &urlCache{
		db: db,
	}
}

func (s *urlCache) Set(shortURL, longURL string) error {
	return s.db.redisClient.Set(context.Background(), shortURL, longURL, 0).Err()
}

func (s *urlCache) Get(shortURL string) (string, error) {
	longURL, err := s.db.redisClient.Get(context.Background(), shortURL).Result()
	if err != nil {
		return "", err
	}
	return longURL, nil
}
