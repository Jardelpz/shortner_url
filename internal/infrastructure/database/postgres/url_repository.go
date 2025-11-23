package postgres

import (
	"context"
	"database/sql"
	"errors"
	"short_url/internal/domain"
)

type UrlRepository struct {
	db *sql.DB
}

func NewUrlRepository(conn *sql.DB) *UrlRepository {
	return &UrlRepository{db: conn}
}

func (ur *UrlRepository) Insert(ctx context.Context, url *domain.Url) error {
	return nil
}

func (ur *UrlRepository) Find(ctx context.Context, shortUrl string) (*domain.Url, error) {
	var url domain.Url
	row := ur.db.QueryRowContext(ctx, `
		SELECT long_url, short_url
		FROM table_url
		WHERE short_url = $1
	`, shortUrl)

	err := row.Scan(&url.LongUrl, &url.ShortUrl)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUrlNotFound
		}

		return nil, err

	}

	return &url, nil
}
