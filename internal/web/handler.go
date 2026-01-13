package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/likaia/nginxpulse/internal/analytics"
	"github.com/likaia/nginxpulse/internal/config"
	"github.com/likaia/nginxpulse/internal/ingest"
	"github.com/sirupsen/logrus"
)

// 初始化Web路由
func SetupRoutes(
	router *gin.Engine,
	statsFactory *analytics.StatsFactory) {

	// 获取所有网站列表
	router.GET("/api/websites", func(c *gin.Context) {
		websiteIDs := config.GetAllWebsiteIDs()

		websites := make([]map[string]string, 0, len(websiteIDs))
		for _, id := range websiteIDs {
			website, ok := config.GetWebsiteByID(id)
			if !ok {
				continue
			}

			websites = append(websites, map[string]string{
				"id":   id,
				"name": website.Name,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"websites": websites,
		})
	})

	router.GET("/api/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"log_parsing":          ingest.IsIPParsing(),
			"log_parsing_progress": ingest.GetIPParsingProgress(),
		})
	})

	// 查询接口
	router.GET("/api/stats/:type", func(c *gin.Context) {
		statsType := c.Param("type")
		params := make(map[string]string)
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		query, err := statsFactory.BuildQueryFromRequest(statsType, params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 执行查询
		result, err := statsFactory.QueryStats(statsType, query)
		if err != nil {
			logrus.WithError(err).Errorf("查询统计数据[%s]失败", statsType)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("查询失败: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

}
