package helper

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/config"
)

func GetConfig(c *gin.Context) *config.Config {
	return c.MustGet("config").(*config.Config)
}

var sizes = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}

func FormatFileSize(s float64, base float64) string {
	unitsLimit := len(sizes)
	i := 0
	for s >= base && i < unitsLimit {
		s = s / base
		i++
	}

	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}

	return fmt.Sprintf(f, s, sizes[i])
}
