package domain

import "context"

type UrlRepository interface {
	Insert(ctx context.Context, url *Url) error
	Find(ctx context.Context, shortUrl string) (*Url, error)
}
