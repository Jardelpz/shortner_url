package domain

import (
	"context"
	"time"
)

type UrlCache interface {
	Get(ctx context.Context, shortUrl string) (*Url, error)
	Set(ctx context.Context, url *Url, ttl time.Duration) error
	Delete(ctx context.Context, shortUrl string) error
}
