package server

import (
	"github.com/gin-gonic/gin"
)

func router(router gin.IRouter) {
	api := router.Group("/api/v1")
	{
		api.GET("/minted", getMints)
		api.GET("/count", getTotalCount)
	}
}
