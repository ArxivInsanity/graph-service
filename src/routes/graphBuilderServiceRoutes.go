package routes

import (
	"github.com/ArxivInsanity/graph-service/src/services"
	"github.com/gin-gonic/gin"
)

func GraphBuilderRoutes(routerGroup *gin.RouterGroup) {
	// defining graphBuilder Routes
	routerGroup.GET("/:paperId", services.BuildGraphHandler())
}
