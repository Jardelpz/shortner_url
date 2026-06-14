package http

import (
	"context"
	"short_url/internal/domain"
)

type UrlService interface {
	InsertUrl(ctx context.Context, longUrl string) (*domain.Url, error)
	GetLongUrl(ctx context.Context, shortUrl string) (*domain.Url, error)
}
