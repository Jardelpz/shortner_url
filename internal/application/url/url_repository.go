package url

import (
	"context"
	"short_url/internal/domain"
)

type UrlRepository interface {
	Insert(ctx context.Context, url *domain.Url) error
	Find(ctx context.Context, shortUrl string) (*domain.Url, error)
}
