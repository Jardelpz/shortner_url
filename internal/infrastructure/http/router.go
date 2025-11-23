package http

import "github.com/gin-gonic/gin"

func NewRouter(urlHandler *UrlHandler) *gin.Engine {
	r := gin.New()
	// todo span_id, trace_id, parent_trace_id
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", urlHandler.HealthCheck)

	v1 := r.Group("/v1")
	{
		v1.GET("/url", urlHandler.GetLongUrlHandler)
		v1.POST("/url", urlHandler.CreateShortUrlHandler)
	}

	return r
}
