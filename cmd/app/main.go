package main

import (
	"log"

	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/app"

	_ "github.com/minhmannh2001/authconnecthub/docs"
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
	// Configuration
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
