// go:build (darwin && cgo) || linux

package main

import (
	. "github.com/ArxivInsanity/graph-service/src/common"
	. "github.com/ArxivInsanity/graph-service/src/routes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {

	// loading env vars
	loadEnv()

	// init engine
	r := gin.Default()

	// load routes
	graphBuilderGroup := r.Group("/graphBuilderGroup")
	{
		GraphBuilderRoutes(graphBuilderGroup)
	}

	// default route
	r.GET("/", func(c *gin.Context) {
		//c.Header("Content-Type", "application/json")
		c.IndentedJSON(http.StatusOK, "Hello Welcome to ArxivInsanity GraphBuilder Service")
	})

	err := r.Run()
	PanicOnErr(err)
}

func loadEnv() {
	// loading environment variables
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resources")
	err := viper.ReadInConfig()
	PanicOnErr(err)
}
