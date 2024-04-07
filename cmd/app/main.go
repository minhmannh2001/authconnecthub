package main

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/config"
	_ "github.com/minhmannh2001/authconnecthub/docs"
	router "github.com/minhmannh2001/authconnecthub/internal/controller/http"
	"github.com/minhmannh2001/authconnecthub/internal/middlewares"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/repos"
	"github.com/minhmannh2001/authconnecthub/pkg/httpserver"
	"github.com/minhmannh2001/authconnecthub/pkg/logger"
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	"github.com/minhmannh2001/authconnecthub/pkg/redis"
)

// @title 		  AuthConnect Hub
// @version       1.0
// @description   A centralized authentication hub for my home applications in Go using Gin framework.

// @contact.name  Nguyen Minh Manh
// @contact.email nguyenminhmannh2001@gmail.com

// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host          localhost:8080
// @BasePath      /v1
func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			logger.NewLogger,
			config.NewConfig,
			postgres.New,
			redis.New,
			fx.Annotate(
				repos.NewAuthRepo,
				fx.As(new(repos.IAuthRepo)),
			),
			fx.Annotate(
				repos.NewRoleRepo,
				fx.As(new(repos.IRoleRepo)),
			),
			fx.Annotate(
				repos.NewUserRepo,
				fx.As(new(repos.IUserRepo)),
			),
			fx.Annotate(
				usecases.NewRoleUseCase,
				fx.As(new(usecases.IRoleUC)),
			),
			fx.Annotate(
				usecases.NewAuthUseCase,
				fx.As(new(usecases.IAuthUC)),
			),
			fx.Annotate(
				usecases.NewUserUseCase,
				fx.As(new(usecases.IUserUC)),
			),
			func(cfg *config.Config, authUseCase usecases.IAuthUC) *gin.Engine {
				e := gin.New()
				// Middlewares
				e.Use(func(c *gin.Context) {
					c.Set("config", cfg)
					c.Next()
				})
				e.Use(middlewares.IsHtmxRequest)
				e.Use(middlewares.IsLoggedIn(authUseCase))
				e.Use(middlewares.IsAuthorized(authUseCase))
				e.Use(gin.Logger())
				e.Use(gin.Recovery())

				return e
			},
			router.New,
			func(e *gin.Engine, cfg *config.Config, l *slog.Logger) *httpserver.Server {
				httpServer := httpserver.New(e, httpserver.Port(cfg.App.Port))
				l.Info("Server was started", slog.String("host", cfg.App.Host), slog.String("port", cfg.App.Port))
				return httpServer
			},
		),
		fx.Invoke(
			setLifeCycle,
		),
	)

	app.Run()
}

func setLifeCycle(
	lc fx.Lifecycle,
	r *router.HTTP,
	s *httpserver.Server,
	e *gin.Engine,
	l *slog.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			r.Start(e)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := s.Shutdown()
			if err != nil {
				l.Error("app - Run - httpServer.Shutdown", slog.Any("error", err))
				return err
			}
			l.Info("Server was shutdown")
			return nil
		},
	})
}

// https://github.com/emarifer/goCMS/blob/main/cmd/gocms_admin/main.go
// clean template: https://github.com/evrone/go-clean-template/tree/master
// https://amitshekhar.me/blog/go-backend-clean-architecture
// https://github.com/amitshekhariitbhu/go-backend-clean-architecture
// go, gorm & gin crud example: https://github.com/herusdianto/gorm_crud_example/tree/master
// go fx: https://juejin.cn/post/7153582825399124005#heading-2
// RBAC:https://dev.to/bensonmacharia/role-based-access-control-in-golang-with-jwt-go-ijn
