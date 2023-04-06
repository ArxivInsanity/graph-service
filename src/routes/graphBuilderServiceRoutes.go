package routes

import (
	"github.com/gin-gonic/gin"
)
import "github.com/ArxivInsanity/graph-service/src/services"

func GraphBuilderRoutes(routerGroup *gin.RouterGroup) {
	// defining graphBuilder Routes
	routerGroup.GET("/getGraph/:paperId", services.GetGraph())
}
