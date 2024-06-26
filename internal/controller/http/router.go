package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	v1 "github.com/minhmannh2001/authconnecthub/internal/controller/http/v1"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// HTTP controller
type HTTP struct {
	logger *slog.Logger
	authUC usecases.IAuthUC
	userUC usecases.IUserUC
	roleUC usecases.IRoleUC
}

// New returns a new HTTP controller
func New(l *slog.Logger, a usecases.IAuthUC, u usecases.IUserUC, r usecases.IRoleUC) *HTTP {
	return &HTTP{
		logger: l,
		authUC: a,
		userUC: u,
		roleUC: r,
	}
}

// Start starts the HTTP controller
func (h *HTTP) Start(e *gin.Engine) {
	e.Static("/static", "./static")
	e.LoadHTMLGlob("templates/*")
	// Prometheus metrics
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))

	e.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.registerRoutes(e)
}

func (h *HTTP) registerRoutes(e *gin.Engine) {
	e.GET("/500.html", handleInternalServerError)
	e.NoRoute(handleNotFound)
	e.NoMethod(handleNoMethod)

	e.GET("/", homeHandler)

	e.PUT("/show-toast", func(c *gin.Context) {
		c.HTML(http.StatusOK, "toast-section", gin.H{
			"hidden": false,
		})
	})

	e.PUT("/close-toast", func(c *gin.Context) {
		url := strings.Split(c.Request.Header.Get("HX-Current-URL"), "?")[0]
		c.Header("HX-Replace-Url", url)
		c.HTML(http.StatusOK, "toast-section", gin.H{
			"hidden": true,
		})
	})

	e.GET("/private", privateHandler)

	// Routers
	groupRouter := e.Group("/v1")
	{
		v1.NewAuthenRoutes(groupRouter, h.logger, h.authUC, h.userUC, h.roleUC)
		e.GET("/dashboard", dashboardHandler)
	}
}

// @Summary Home Page
// @Description This endpoint renders the home page of the application.
// It accepts optional query parameters for toast notifications and validates them with a hash value.
// @Tags home
// @Produce html
// @Param toast-message query string false "Toast message to display"
// @Param toast-type query string false "Type of toast notification (e.g., success, error)"
// @Param hash-value query string false "Hash value for validation"
// @Success 200 {object} object Response object containing HTML data
// @Router / [GET]
func homeHandler(c *gin.Context) {
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

	userInfo := map[string]interface{}{}
	username, ok := c.Get("username")

	if ok {
		userInfo["username"] = username
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":         "Personal Hub",
		"toastSettings": toastSettings,
		"reload":        c.GetHeader("HX-Reload"),
		"userInfo":      userInfo,
	})
}

// @Summary Access a private resource
// @Description This endpoint is accessible only to authorized users and returns a greeting message.
// @Tags private
// @Security JWT
// @Produce json
// @Success 200 {string} Hello message
// @router /private [GET]
func privateHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello. You are in private path")
}

func handleInternalServerError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.html", gin.H{
		"title": "Personal Hub",
		"toastSettings": map[string]interface{}{
			"hidden": true,
		},
	})
}

func handleNotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"title": "Personal Hub",
		"toastSettings": map[string]interface{}{
			"hidden": true,
		},
	})
}

func handleNoMethod(c *gin.Context) {
	c.HTML(http.StatusNotFound, "405.html", gin.H{
		"title": "Personal Hub",
		"toastSettings": map[string]interface{}{
			"hidden": true,
		},
	})
}

func dashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Personal Hub",
		"toastSettings": map[string]interface{}{
			"hidden": true,
		},
		"reload": c.GetHeader("HX-Reload"),
	})
}
