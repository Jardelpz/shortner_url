package domain

import (
	"context"
)

type UrlCache interface {
	Get(ctx context.Context, shortUrl string) (*Url, error)
	Set(ctx context.Context, url *Url) error
	Delete(ctx context.Context, shortUrl string) error
}
