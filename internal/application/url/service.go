package url

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"short_url/internal/domain"
)

type Service struct {
	repo domain.UrlRepository
}

func NewService(url_repo domain.UrlRepository) *Service {
	return &Service{repo: url_repo}
}

func (s *Service) InsertUrl(ctx context.Context, longUrl string) (*domain.Url, error) {
	shortUrl, err := s.GenerateHashValue(ctx, longUrl)

	err = s.repo.Insert(ctx, &domain.Url{LongUrl: longUrl, ShortUrl: shortUrl})

	if err != nil {
		return &domain.Url{}, err
	}
	return &domain.Url{ShortUrl: shortUrl, LongUrl: longUrl}, nil
}

func (s *Service) GetLongUrl(ctx context.Context, shortUrl string) (*domain.Url, error) {
	url, err := s.repo.Find(ctx, shortUrl)
	if err != nil {
		return nil, err
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
