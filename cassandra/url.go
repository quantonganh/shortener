package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/quantonganh/base62"
	"github.com/quantonganh/snowflake"

	"github.com/quantonganh/shortener"
)

type urlService struct {
	db *DB
}

func NewURLService(db *DB) shortener.URLService {
	return &urlService{
		db: db,
	}
}

func (s *urlService) CreateShortURL(longURL string) (string, error) {
	worker := snowflake.NewWorker(1, 1)
	id, err := worker.NextID()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate unique ID")
	}

	shortURL := base62.Encode(int(id))
	query := s.db.session.Query("INSERT INTO shortener.url (id, short_url, long_url) VALUES (?, ?, ?)", id, shortURL, longURL)
	if err := query.Exec(); err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *urlService) GetLongURL(shortURL string) (string, error) {
	id := base62.Decode(shortURL)
	query := s.db.session.Query("SELECT long_url FROM shortener.url WHERE id = ? LIMIT 1", id)
	var longURL string
	if err := query.Consistency(gocql.One).Scan(&longURL); err != nil {
		return "", err
	}
	return longURL, nil
}
