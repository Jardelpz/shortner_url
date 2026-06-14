package redis

import (
	"context"
	"encoding/json"
	"errors"
	"short_url/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

const CacheTtl = 60 * time.Minute

type UrlCache struct {
	client *redis.Client
}

func NewUrlCache(client *redis.Client) *UrlCache {
	return &UrlCache{client: client}
}

func (c *UrlCache) Get(ctx context.Context, shortUrl string) (*domain.Url, error) {
	val, err := c.client.Get(ctx, cacheKey(shortUrl)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, domain.ErrUrlNotFound
	}
	if err != nil {
		return nil, err
	}

	var url domain.Url
	if err := json.Unmarshal([]byte(val), &url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (c *UrlCache) Set(ctx context.Context, url *domain.Url) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, cacheKey(url.ShortUrl), data, CacheTtl).Err()
}

func (c *UrlCache) Delete(ctx context.Context, shortUrl string) error {
	return c.client.Del(ctx, cacheKey(shortUrl)).Err()
}

func cacheKey(shortUrl string) string {
	return "url:" + shortUrl
}
