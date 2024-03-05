package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/config"
)

func GetConfig(c *gin.Context) *config.Config {
	return c.MustGet("config").(*config.Config)
}
