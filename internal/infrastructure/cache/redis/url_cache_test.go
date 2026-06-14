//go:build integration

package redis_test

import (
	"context"
	"errors"
	"short_url/internal/domain"
	redisinfra "short_url/internal/infrastructure/cache/redis"
	"testing"
)

func TestUrlCache_SetAndGet(t *testing.T) {
	client := redisinfra.ConnectionCache()
	defer client.Close()

	cache := redisinfra.NewUrlCache(client)
	ctx := context.Background()

	url := &domain.Url{LongUrl: "https://example.com", ShortUrl: "test1234"}
	if err := cache.Set(ctx, url); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	found, err := cache.Get(ctx, "test1234")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if found.LongUrl != url.LongUrl {
		t.Errorf("expected longUrl %s, got %s", url.LongUrl, found.LongUrl)
	}
	if found.ShortUrl != url.ShortUrl {
		t.Errorf("expected shortUrl %s, got %s", url.ShortUrl, found.ShortUrl)
	}
}

func TestUrlCache_Get_NotFound(t *testing.T) {
	client := redisinfra.ConnectionCache()
	defer client.Close()

	cache := redisinfra.NewUrlCache(client)
	ctx := context.Background()

	_, err := cache.Get(ctx, "nonexistent_key_xyz")
	if !errors.Is(err, domain.ErrUrlNotFound) {
		t.Errorf("expected ErrUrlNotFound, got %v", err)
	}
}

func TestUrlCache_Delete(t *testing.T) {
	client := redisinfra.ConnectionCache()
	defer client.Close()

	cache := redisinfra.NewUrlCache(client)
	ctx := context.Background()

	url := &domain.Url{LongUrl: "https://example.com", ShortUrl: "del12345"}
	_ = cache.Set(ctx, url)

	if err := cache.Delete(ctx, "del12345"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := cache.Get(ctx, "del12345")
	if !errors.Is(err, domain.ErrUrlNotFound) {
		t.Errorf("expected ErrUrlNotFound after delete, got %v", err)
	}
}
