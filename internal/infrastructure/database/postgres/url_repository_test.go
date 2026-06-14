//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"short_url/internal/domain"
	"short_url/internal/infrastructure/database/postgres"
	"testing"
	"time"
)

func TestUrlRepository_InsertAndFind(t *testing.T) {
	db := postgres.ConnectionDatabase()
	defer db.Close()

	repo := postgres.NewUrlRepository(db)
	ctx := context.Background()

	shortUrl := fmt.Sprintf("%d", time.Now().UnixNano())[:8]
	url := &domain.Url{LongUrl: "https://integration-test.com", ShortUrl: shortUrl}

	if err := repo.Insert(ctx, url); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	found, err := repo.Find(ctx, shortUrl)
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}
	if found.LongUrl != url.LongUrl {
		t.Errorf("expected longUrl %s, got %s", url.LongUrl, found.LongUrl)
	}
	if found.ShortUrl != url.ShortUrl {
		t.Errorf("expected shortUrl %s, got %s", url.ShortUrl, found.ShortUrl)
	}
}

func TestUrlRepository_Find_NotFound(t *testing.T) {
	db := postgres.ConnectionDatabase()
	defer db.Close()

	repo := postgres.NewUrlRepository(db)
	ctx := context.Background()

	_, err := repo.Find(ctx, "nonexistent")
	if !errors.Is(err, domain.ErrUrlNotFound) {
		t.Errorf("expected ErrUrlNotFound, got %v", err)
	}
}
