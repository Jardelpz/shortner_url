package http

import (
	"net/http"
	"short_url/internal/application/url"

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

func (u *UrlHandler) CreateShortUrl(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (u *UrlHandler) GetLongUrl(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}
