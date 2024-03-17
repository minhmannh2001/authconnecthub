// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/config"
	router "github.com/minhmannh2001/authconnecthub/internal/controller/http"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/middleware"
	"github.com/minhmannh2001/authconnecthub/internal/usecase"
	"github.com/minhmannh2001/authconnecthub/internal/usecase/repo"
	"github.com/minhmannh2001/authconnecthub/pkg/httpserver"
	"github.com/minhmannh2001/authconnecthub/pkg/logger"
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	"github.com/minhmannh2001/authconnecthub/pkg/redis"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	_ = l

	pg, err := postgres.New(*cfg)

	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	redis, err := redis.New(*cfg)

	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}

	// Repo
	authRepo := repo.NewAuthRepo(pg, redis)
	userRepo := repo.NewUserRepo(pg)
	roleRepo := repo.NewRoleRepo(pg)

	// Use case
	roleUseCase := usecase.NewRoleUseCase(roleRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthUseCase(authRepo, *userUseCase, cfg.JwtPrivateKey)

	// HTTP Server
	handler := gin.New()

	_, err = helper.GetSwaggerInfo(cfg.App.SwaggerPath)
	if err != nil {
		panic(err)
	}

	// Middlewares
	handler.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})
	handler.Use(middleware.IsHtmxRequest)
	handler.Use(middleware.IsLoggedIn(authUseCase))
	handler.Use(middleware.IsAuthorized(authUseCase))

	router.NewRouter(handler, l, authUseCase, userUseCase, roleUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.App.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
