package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"short_url/internal/domain"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockUrlService struct {
	insertUrlFn  func(ctx context.Context, longUrl string) (*domain.Url, error)
	getLongUrlFn func(ctx context.Context, shortUrl string) (*domain.Url, error)
}

func (m *mockUrlService) InsertUrl(ctx context.Context, longUrl string) (*domain.Url, error) {
	return m.insertUrlFn(ctx, longUrl)
}

func (m *mockUrlService) GetLongUrl(ctx context.Context, shortUrl string) (*domain.Url, error) {
	return m.getLongUrlFn(ctx, shortUrl)
}

func newTestRouter(svc UrlService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(NewUrlHandler(svc))
}

func TestHealthCheck(t *testing.T) {
	router := newTestRouter(&mockUrlService{})

	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCreateShortUrlHandler_Success(t *testing.T) {
	svc := &mockUrlService{
		insertUrlFn: func(ctx context.Context, longUrl string) (*domain.Url, error) {
			return &domain.Url{LongUrl: longUrl, ShortUrl: "abc12345"}, nil
		},
	}
	router := newTestRouter(svc)

	body, _ := json.Marshal(map[string]string{"longUrl": "https://example.com"})
	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodPost, "/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["shorUrl"] != "abc12345" {
		t.Errorf("expected shorUrl abc12345, got %s", res["shorUrl"])
	}
}

func TestCreateShortUrlHandler_InvalidPayload(t *testing.T) {
	router := newTestRouter(&mockUrlService{})

	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodPost, "/v1/url", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestCreateShortUrlHandler_ServiceError(t *testing.T) {
	svc := &mockUrlService{
		insertUrlFn: func(ctx context.Context, longUrl string) (*domain.Url, error) {
			return nil, errors.New("internal error")
		},
	}
	router := newTestRouter(svc)

	body, _ := json.Marshal(map[string]string{"longUrl": "https://example.com"})
	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodPost, "/v1/url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetLongUrlHandler_Success(t *testing.T) {
	svc := &mockUrlService{
		getLongUrlFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return &domain.Url{LongUrl: "https://example.com", ShortUrl: shortUrl}, nil
		},
	}
	router := newTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodGet, "/v1/url?shortUrl=abc12345", nil)
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["longUrl"] != "https://example.com" {
		t.Errorf("expected longUrl https://example.com, got %s", res["longUrl"])
	}
}

func TestGetLongUrlHandler_MissingParam(t *testing.T) {
	router := newTestRouter(&mockUrlService{})

	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodGet, "/v1/url", nil)
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetLongUrlHandler_NotFound(t *testing.T) {
	svc := &mockUrlService{
		getLongUrlFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return nil, domain.ErrUrlNotFound
		},
	}
	router := newTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodGet, "/v1/url?shortUrl=notfound", nil)
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetLongUrlHandler_Timeout(t *testing.T) {
	svc := &mockUrlService{
		getLongUrlFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return nil, context.DeadlineExceeded
		},
	}
	router := newTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodGet, "/v1/url?shortUrl=abc12345", nil)
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusGatewayTimeout {
		t.Errorf("expected 504, got %d", w.Code)
	}
}

func TestGetLongUrlHandler_InternalError(t *testing.T) {
	svc := &mockUrlService{
		getLongUrlFn: func(ctx context.Context, shortUrl string) (*domain.Url, error) {
			return nil, errors.New("unexpected error")
		},
	}
	router := newTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := nethttp.NewRequest(nethttp.MethodGet, "/v1/url?shortUrl=abc12345", nil)
	router.ServeHTTP(w, req)

	if w.Code != nethttp.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
