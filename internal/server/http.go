package server

import (
	"net/http"
	"strings"
	"time"
	"github.com/likaia/nginxpulse/internal/analytics"
	"github.com/likaia/nginxpulse/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// StartHTTPServer configures and starts the HTTP server in a goroutine.
func StartHTTPServer(statsFactory *analytics.StatsFactory, addr string) *http.Server {
	router := buildRouter(statsFactory)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("Failed to start the server")
		}
	}()

	logrus.Infof("服务器已启动，监听地址: %s", addr)
	return server
}

func buildRouter(statsFactory *analytics.StatsFactory) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(requestLogger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	web.SetupRoutes(router, statsFactory)

	return router
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()

		if status >= 400 {
			logrus.Warnf("HTTP %d %s %s %s %v",
				status, c.Request.Method, path, c.ClientIP(), duration)
			return
		}

		if strings.HasPrefix(path, "/api/") && duration > 100*time.Millisecond {
			logrus.Warnf("高延迟 %s %s %d %s %v",
				c.Request.Method, path, status, c.ClientIP(), duration)
		}
	}
}
