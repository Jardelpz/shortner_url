package postgres

import (
	"database/sql"
	"short_url/internal/domain"
)

type UrlRepository struct {
	db *sql.DB
}

func NewUrlRepository(conn *sql.DB) *UrlRepository {
	return &UrlRepository{db: conn}
}

func (ur *UrlRepository) Insert(url *domain.Url) error {
	return nil
}

func (ur *UrlRepository) Find(shortUrl string) (*domain.Url, error) {
	return &domain.Url{}, nil
}
