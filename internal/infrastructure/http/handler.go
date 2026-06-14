package http

import (
	"context"
	"errors"
	"net/http"
	"short_url/internal/domain"
	"short_url/internal/infrastructure/http/dto"
	"time"

	"github.com/gin-gonic/gin"
)

type UrlHandler struct {
	service UrlService
}

func NewUrlHandler(svc UrlService) *UrlHandler {
	return &UrlHandler{service: svc}
}

func (u *UrlHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive and kicking"})
}

func (u *UrlHandler) CreateShortUrlHandler(c *gin.Context) {
	// validar se nao existe antes de inserir

	var req dto.CreateShortUrlInput
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "missing payload data: longUrl"})
	}

	ctx := c.Request.Context()
	url, err := u.service.InsertUrl(ctx, req.LongUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusCreated, gin.H{"longUrl": url.LongUrl, "shorUrl": url.ShortUrl})
}

func (u *UrlHandler) GetLongUrlHandler(c *gin.Context) {
	shortUrl := c.Query("shortUrl")
	if shortUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query param: shortUrl"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	url, err := u.service.GetLongUrl(ctx, shortUrl)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUrlNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"longUrl": url.LongUrl, "shorUrl": url.ShortUrl})
}
