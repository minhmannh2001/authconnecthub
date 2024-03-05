package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	v1 "github.com/minhmannh2001/authconnecthub/internal/controller/http/v1"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecase"
	"github.com/minhmannh2001/authconnecthub/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, l logger.Interface, a usecase.Auth, u usecase.User, r usecase.Role) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.Static("/static", "./static")
	handler.LoadHTMLGlob("templates/*")
	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	handler.GET("/", func(c *gin.Context) {
		queryParams := c.Request.URL.Query()

		toastMessage := helper.ExtractQueryParam(queryParams, "toast-message", "")
		toastType := helper.ExtractQueryParam(queryParams, "toast-type", "")
		hashValue := helper.ExtractQueryParam(queryParams, "hash-value", "")

		isValid := helper.IsMapValid(map[string]interface{}{
			"toast-message": toastMessage,
			"toast-type":    toastType,
		}, hashValue)

		toastSettings := map[string]interface{}{
			"hidden":  !isValid, // Toggle based on validity
			"type":    toastType,
			"message": helper.FormatToastMessage(toastMessage),
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":         "AuthConnect Hub",
			"toastSettings": toastSettings,
		})
	})

	handler.PUT("/show-toast", func(c *gin.Context) {
		c.HTML(http.StatusOK, "toast-section", gin.H{
			"hidden": false,
		})
	})

	handler.PUT("/close-toast", func(c *gin.Context) {
		url := strings.Split(c.Request.Header.Get("HX-Current-URL"), "?")[0]
		c.Header("HX-Replace-Url", url)
		c.HTML(http.StatusOK, "toast-section", gin.H{
			"hidden": true,
		})
	})

	handler.GET("/private", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello")
	})

	// Routers
	h := handler.Group("/v1")
	{
		v1.NewAuthenticationRoutes(h, l, a, u, r)
	}
}
