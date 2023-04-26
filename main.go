// go:build (darwin && cgo) || linux

package main

import (
	. "github.com/ArxivInsanity/graph-service/src/db"
	. "github.com/ArxivInsanity/graph-service/src/routes"
	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func init() {
	// Loading ENV vars
	loadEnv()
}

func main() {
	// Init Engine
	r := gin.Default()

	// Get Neo4j DB connection
	neo4jContext, neo4jSession := GetNeo4jContextAndSession()

	// Get Redis connection
	redisStore := GetRedisClient()

	// Load routes
	graphBuilderGroup := r.Group("/graphBuilder")
	{
		// Inject neo4j ctx and session in middleware
		graphBuilderGroup.Use(InjectNeo4jContextAndSession(neo4jContext, neo4jSession))
		GraphBuilderRoutes(graphBuilderGroup)
	}
	graphSearchGroup := r.Group("/graphSearch")
	{
		// Inject neo4j ctx and session in middleware
		graphSearchGroup.Use(InjectNeo4jContextAndSession(neo4jContext, neo4jSession))
		cachingMiddleware := GetCachingMiddleware(redisStore)
		GraphSearchRoutes(graphSearchGroup, cachingMiddleware)
	}

	// Default Route
	r.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, "Hello Welcome to ArxivInsanity GraphBuilder Service")
	})

	err := r.Run()
	defer PanicOnClosureError(err, neo4jContext, neo4jSession)
}

func loadEnv() {
	// Loading Environment Variables
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resources")
	err := viper.ReadInConfig()
	PanicOnErr(err)
}
