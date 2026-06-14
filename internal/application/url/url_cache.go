package url

import (
	"context"
	"short_url/internal/domain"
)

type UrlCache interface {
	Get(ctx context.Context, shortUrl string) (*domain.Url, error)
	Set(ctx context.Context, url *domain.Url) error
	Delete(ctx context.Context, shortUrl string) error
}
