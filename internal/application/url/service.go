package url

import (
	"context"
	"short_url/internal/domain"
)

type Service struct {
	repo domain.UrlRepository
}

func NewService(url_repo domain.UrlRepository) *Service {
	return &Service{repo: url_repo}
}

func (s *Service) InsertUrl(ctx context.Context, url domain.Url) error {
	// validar se nao existe antes de inserir
	err := s.repo.Insert(&url)
	if err != nil {
		return err
	}
	return nil
}
