package api

import (
	"mapper/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	store  *store.RedisStore
	logger *zap.SugaredLogger
}

func New(redis *store.RedisStore, logger *zap.SugaredLogger) *API {
	return &API{
		store:  redis,
		logger: logger,
	}
}

func (a *API) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "WORKING")
	})
	router.GET("/streams/:stream", a.GetStream)

	return router
}

func (a *API) GetStream(c *gin.Context) {
	term := c.Param("stream")

	host, err := a.store.GetStream(c.Request.Context(), term)
	if err == store.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "stream not found"})

		return
	}

	if err != nil {
		a.logger.Errorw("failed to get host", "error", err, "stream", term)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get host"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"host": host})
}
