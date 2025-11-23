package http

import (
	"context"
	"errors"
	"net/http"
	"short_url/internal/application/url"
	"short_url/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

type UrlHandler struct {
	service *url.Service
}

func NewUrlHandler(svc *url.Service) *UrlHandler {
	return &UrlHandler{service: svc}
}

func (u *UrlHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive and kicking"})
}

func (u *UrlHandler) CreateShortUrlHandler(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (u *UrlHandler) GetLongUrlHandler(c *gin.Context) {
	longURL := c.Query("longUrl")
	if longURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query param: longUrl"})
		return
	}

	// if _, err := netUrl.ParseRequestURI(longURL); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longUrl"})
	// 	return
	// }

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	url, err := u.service.GetLongUrl(ctx, longURL)
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
