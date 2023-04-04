// go:build (darwin && cgo) || linux

package main

import (
	. "github.com/ArxivInsanity/graph-service/src/common"
	. "github.com/ArxivInsanity/graph-service/src/db"
	. "github.com/ArxivInsanity/graph-service/src/routes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func init() {
	// loading env vars
	loadEnv()
}

func main() {
	// init engine
	r := gin.Default()

	// get Neo4j DB connection
	neo4jContext, neo4jSession := GetNeo4jContextAndSession()

	// load routes
	graphBuilderGroup := r.Group("/graphBuilder")
	{
		// inject neo4j ctx and session in middleware
		graphBuilderGroup.Use(InjectNeo4jContextAndSession(neo4jContext, neo4jSession))
		GraphBuilderRoutes(graphBuilderGroup)
	}

	// default route
	r.GET("/", func(c *gin.Context) {
		//c.Header("Content-Type", "application/json")
		c.IndentedJSON(http.StatusOK, "Hello Welcome to ArxivInsanity GraphBuilder Service")
	})

	err := r.Run()
	PanicOnClosureError(err, neo4jContext, neo4jSession)
}

func loadEnv() {
	// loading environment variables
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resources")
	err := viper.ReadInConfig()
	PanicOnErr(err)
}
