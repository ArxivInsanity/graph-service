package services

import (
	"log"
	"net/http"

	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

func IsSeedPaperHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		log.Printf("Seed Paper handler URL Params: %v", paperId)
		isSeed := getNode(paperId, ctx)
		ctx.JSON(http.StatusOK, isSeed)
	}
}

func IsSeedPaper(paperId string, ctx *gin.Context) bool {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	isSeedCypher := "RETURN EXISTS( (:SEED_PAPER)-[:SEED]-(:PAPER {paperId: $paperId})) as isSeed"
	cypherParam := map[string]any{
		"paperId": paperId,
	}
	result, err := dbSession.Run(dbContext, isSeedCypher, cypherParam)
	PanicOnErr(err)
	record, err := neo4j.CollectWithContext(dbContext, result, err)
	isSeed, _ := record[0].Get("isSeed")
	log.Printf("Paper: %s is Seed Paper: %v", paperId, isSeed)
	PanicOnErr(err)
	return isSeed.(bool)
}

type Node struct {
	paperId       string
	citationCount int
	year          int
	title         string
	citation      []*Node
	reference     []*Node
}

func getNode(rootPaperId string, ctx *gin.Context) *Node {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	cypher := "MATCH (p:PAPER {paperId: $rootPaperId}) RETURN p"
	// referenceCypher := "MATCH (p:PAPER {paperId: $rootPaperId})-[r:REFERENCE]-(n:PAPER) RETURN p, n"
	cypherParam := map[string]any{
		"rootPaperId": rootPaperId,
	}
	result, err := dbSession.Run(dbContext, cypher, cypherParam)
	PanicOnErr(err)
	citationRecord, err := neo4j.CollectWithContext(dbContext, result, err)
	PanicOnErr(err)
	citationNeo4jRecord, _ := citationRecord[0].Get("p")
	citationNeo4jProps := citationNeo4jRecord.(dbtype.Node).Props
	n := Node{paperId: rootPaperId,
		citationCount: citationNeo4jProps["citationCount"].(int),
		year:          citationNeo4jProps["year"].(int),
		title:         citationNeo4jProps["title"].(string)}
	return &n
}

// func bfs(n *Node, depth int, ctx *gin.Context) {
// 	queue := []*Node{n}
// 	visited := map[string]Node{}

// 	for len(queue) > 0 && depth > 0 {
// 		level_size := len(queue)
// 		cypher := "MATCH (p:PAPER {paperId: $paperId}) - [r:CITATION] -> (n:PAPER) RETURN n"
// 		cypherParam := map[string]any{
// 			"rootPaperId": rootPaperId,
// 		}
// 		for i := 0; i < level_size; i++ {
// 			current := queue[0]
// 			queue = queue[1:]
// 			visited[current.name] = *current
// 			for _, nghr := range n.neighbours {
// 				if _, exists := visited[nghr.name]; !exists {
// 					queue = append(queue, nghr)
// 				}
// 			}
// 		}
// 		depth -= 1
// 	}

// }
