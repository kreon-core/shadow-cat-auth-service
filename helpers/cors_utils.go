package helpers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StandardCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",

			"Authorization",
			"X-Client-ID",
			"X-Client-Signature",

			"X-Device-ID",
			"X-Device-OS-Type",
			"X-Device-OS-Version",
			"X-Device-Model",
			"X-App-Version",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
