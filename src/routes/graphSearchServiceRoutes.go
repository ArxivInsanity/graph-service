package routes

import (
	"github.com/ArxivInsanity/graph-service/src/services"
	"github.com/gin-gonic/gin"
)

func GraphSearchRoutes(routerGroup *gin.RouterGroup, cachingMiddleware gin.HandlerFunc) {
	// defining graphBuilder Routes
	routerGroup.GET("/isSeed/:paperId", services.IsSeedPaperHandler())
	routerGroup.GET("/filteredGraph/:paperId", services.FilteredGraphHandler())
	routerGroup.GET("/graph/:paperId", cachingMiddleware, services.GraphHandler())
}
