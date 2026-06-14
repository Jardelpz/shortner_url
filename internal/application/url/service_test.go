package url

import (
	"context"
	"errors"
	"short_url/internal/domain"
	"testing"
)

type mockUrlRepository struct {
	insertFn func(ctx context.Context, url *domain.Url) error
	findFn   func(ctx context.Context, shortUrl string) (*domain.Url, error)
}

func (m *mockUrlRepository) Insert(ctx context.Context, url *domain.Url) error {
	return m.insertFn(ctx, url)
}

func (m *mockUrlRepository) Find(ctx context.Context, shortUrl string) (*domain.Url, error) {
	return m.findFn(ctx, shortUrl)
}

type mockUrlCache struct {
	getFn    func(ctx context.Context, shortUrl string) (*domain.Url, error)
	setFn    func(ctx context.Context, url *domain.Url) error
	deleteFn func(ctx context.Context, shortUrl string) error
}

func (m *mockUrlCache) Get(ctx context.Context, shortUrl string) (*domain.Url, error) {
	return m.getFn(ctx, shortUrl)
}

func (m *mockUrlCache) Set(ctx context.Context, url *domain.Url) error {
	return m.setFn(ctx, url)
}

func (m *mockUrlCache) Delete(ctx context.Context, shortUrl string) error {
	return m.deleteFn(ctx, shortUrl)
}

func noopCache() *mockUrlCache {
	return &mockUrlCache{
		getFn:    func(ctx context.Context, shortUrl string) (*domain.Url, error) { return nil, domain.ErrUrlNotFound },
		setFn:    func(ctx context.Context, url *domain.Url) error { return nil },
		deleteFn: func(ctx context.Context, shortUrl string) error { return nil },
	}
}

func TestInsertUrl_Success(t *testing.T) {
	var inserted *domain.Url
	repo := &mockUrlRepository{
		findFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return nil, domain.ErrUrlNotFound
		},
		insertFn: func(ctx context.Context, url *domain.Url) error {
			inserted = url
			return nil
		},
	}
	svc := NewService(repo, noopCache())

	result, err := svc.InsertUrl(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LongUrl != "https://example.com" {
		t.Errorf("expected longUrl https://example.com, got %s", result.LongUrl)
	}
	if inserted == nil {
		t.Error("Insert was not called on repo")
	}
}

func TestInsertUrl_HashCollision(t *testing.T) {
	callCount := 0
	repo := &mockUrlRepository{
		findFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			callCount++
			if callCount == 1 {
				return &domain.Url{ShortUrl: shortUrl, LongUrl: "other"}, nil
			}
			return nil, domain.ErrUrlNotFound
		},
		insertFn: func(ctx context.Context, url *domain.Url) error { return nil },
	}
	svc := NewService(repo, noopCache())

	result, err := svc.InsertUrl(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount < 2 {
		t.Error("expected Find to be called at least twice (collision + retry with salt)")
	}
	if result.ShortUrl == "" {
		t.Error("expected a non-empty short url")
	}
}

func TestInsertUrl_RepoInsertError(t *testing.T) {
	repo := &mockUrlRepository{
		findFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return nil, domain.ErrUrlNotFound
		},
		insertFn: func(ctx context.Context, url *domain.Url) error {
			return errors.New("db error")
		},
	}
	svc := NewService(repo, noopCache())

	_, err := svc.InsertUrl(context.Background(), "https://example.com")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetLongUrl_CacheHit(t *testing.T) {
	expected := &domain.Url{LongUrl: "https://example.com", ShortUrl: "abc12345"}
	repoCalled := false

	repo := &mockUrlRepository{
		findFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			repoCalled = true
			return nil, nil
		},
		insertFn: func(ctx context.Context, url *domain.Url) error { return nil },
	}
	cache := &mockUrlCache{
		getFn:    func(ctx context.Context, shortUrl string) (*domain.Url, error) { return expected, nil },
		setFn:    func(ctx context.Context, url *domain.Url) error { return nil },
		deleteFn: func(ctx context.Context, shortUrl string) error { return nil },
	}
	svc := NewService(repo, cache)

	result, err := svc.GetLongUrl(context.Background(), "abc12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LongUrl != expected.LongUrl {
		t.Errorf("expected %s, got %s", expected.LongUrl, result.LongUrl)
	}
	if repoCalled {
		t.Error("repo should not be called on cache hit")
	}
}

func TestGetLongUrl_CacheMiss_FoundInRepo(t *testing.T) {
	expected := &domain.Url{LongUrl: "https://example.com", ShortUrl: "abc12345"}
	setCalled := false

	repo := &mockUrlRepository{
		findFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return expected, nil
		},
		insertFn: func(ctx context.Context, url *domain.Url) error { return nil },
	}
	cache := &mockUrlCache{
		getFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return nil, domain.ErrUrlNotFound
		},
		setFn: func(ctx context.Context, url *domain.Url) error {
			setCalled = true
			return nil
		},
		deleteFn: func(ctx context.Context, shortUrl string) error { return nil },
	}
	svc := NewService(repo, cache)

	result, err := svc.GetLongUrl(context.Background(), "abc12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LongUrl != expected.LongUrl {
		t.Errorf("expected %s, got %s", expected.LongUrl, result.LongUrl)
	}
	if !setCalled {
		t.Error("cache.Set should be called after fetching from repo")
	}
}

func TestGetLongUrl_NotFound(t *testing.T) {
	repo := &mockUrlRepository{
		findFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return nil, domain.ErrUrlNotFound
		},
		insertFn: func(ctx context.Context, url *domain.Url) error { return nil },
	}
	svc := NewService(repo, noopCache())

	_, err := svc.GetLongUrl(context.Background(), "notfound")
	if !errors.Is(err, domain.ErrUrlNotFound) {
		t.Errorf("expected ErrUrlNotFound, got %v", err)
	}
}
