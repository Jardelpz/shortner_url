package url

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"short_url/internal/domain"
)

type Service struct {
	repo  domain.UrlRepository
	cache domain.UrlCache
}

func NewService(urlRepo domain.UrlRepository, cache domain.UrlCache) *Service {
	return &Service{repo: urlRepo, cache: cache}
}

func (s *Service) InsertUrl(ctx context.Context, longUrl string) (*domain.Url, error) {
	shortUrl, err := s.GenerateHashValue(ctx, longUrl)
	if err != nil {
		return &domain.Url{}, err
	}

	url := &domain.Url{LongUrl: longUrl, ShortUrl: shortUrl}

	err = s.repo.Insert(ctx, url)
	if err != nil {
		return &domain.Url{}, err
	}

	if err := s.cache.Set(ctx, url); err != nil {
		log.Printf("cache set error: %v", err)
	}

	return url, nil
}

func (s *Service) GetLongUrl(ctx context.Context, shortUrl string) (*domain.Url, error) {
	url, err := s.cache.Get(ctx, shortUrl)
	if err == nil {
		log.Printf("value recovered from cache: %v", shortUrl)
		return url, nil
	}

	url, err = s.repo.Find(ctx, shortUrl)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, url); err != nil {
		log.Printf("cache set error: %v", err)
	}

	return url, nil
}

func (s *Service) GenerateHashValue(ctx context.Context, longUrl string) (string, error) {
	h := md5.Sum([]byte(longUrl))
	hash := hex.EncodeToString(h[:])[:8]

	url, err := s.repo.Find(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrUrlNotFound) {
			return hash, nil
		}
		return "", err
	}

	if url != nil {
		fmt.Print("Using salt")
		longUrl = longUrl + "p" // salt
		return s.GenerateHashValue(ctx, longUrl)
	}

	return hash, nil
}
