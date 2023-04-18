package routes

import (
	"github.com/ArxivInsanity/graph-service/src/services"
	"github.com/gin-gonic/gin"
)

func GraphSearchRoutes(routerGroup *gin.RouterGroup) {
	// defining graphBuilder Routes
	routerGroup.GET("/isSeed/:paperId", services.IsSeedPaperHandler())
	routerGroup.GET("/graph/:paperId", services.GraphHandler())
}
