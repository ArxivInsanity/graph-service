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
		isSeed := IsSeedPaper(paperId, ctx)
		ctx.JSON(http.StatusOK, isSeed)
	}
}

func GraphHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		log.Printf("graph Paper handler URL Params: %v", paperId)
		isSeed := IsSeedPaper(paperId, ctx)
		if !isSeed {
			BuildGraph(paperId, ctx)
		}
		seedNode := GetNode(paperId, ctx)
		visited := bfs(seedNode, 3, true, ctx)
		ctx.JSON(http.StatusOK, visited)
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

func GetNode(rootPaperId string, ctx *gin.Context) Node {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	cypher := "MATCH (p:PAPER {paperId: $rootPaperId}) RETURN p"
	// referenceCypher := "MATCH (p:PAPER {paperId: $rootPaperId})-[r:REFERENCE]-(n:PAPER) RETURN p, n"
	cypherParam := map[string]any{
		"rootPaperId": rootPaperId,
	}
	result, err := dbSession.Run(dbContext, cypher, cypherParam)
	PanicOnErr(err)
	nodeRecord, err := neo4j.CollectWithContext(dbContext, result, err)
	PanicOnErr(err)
	nodeNeo4jRecord, _ := nodeRecord[0].Get("p")
	nodeNeo4jProps := nodeNeo4jRecord.(dbtype.Node).Props
	n := Node{PaperId: rootPaperId,
		CitationCount: nodeNeo4jProps["citationCount"].(int64),
		Year:          nodeNeo4jProps["year"].(int64),
		Title:         nodeNeo4jProps["title"].(string)}
	return n
}

func bfs(n Node, depth int, citation bool, ctx *gin.Context) map[string]Node {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	queue := []Node{n}
	visited := map[string]Node{}

	for len(queue) > 0 && depth > 0 {
		level_size := len(queue)
		for i := 0; i < level_size; i++ {
			current := queue[0]
			queue = queue[1:]
			cypher := "MATCH (p:PAPER {paperId: $paperId}) - [r:REFERENCES] -> (n:PAPER) RETURN n"
			cypherParam := map[string]any{
				"paperId": current.PaperId,
			}
			result, err := dbSession.Run(dbContext, cypher, cypherParam)
			PanicOnErr(err)
			linkRecord, err := neo4j.CollectWithContext(dbContext, result, err)
			PanicOnErr(err)
			for i := 0; i < len(linkRecord); i++ {
				linkNeo4jRecord, _ := linkRecord[i].Get("n")
				linkNeo4jProps := linkNeo4jRecord.(dbtype.Node).Props
				child := Node{PaperId: linkNeo4jProps["paperId"].(string),
					CitationCount: linkNeo4jProps["citationCount"].(int64),
					Year:          linkNeo4jProps["year"].(int64),
					Title:         linkNeo4jProps["title"].(string)}
				current.Reference = append(current.Reference, child)
			}
			for _, child := range current.Reference {
				if _, exists := visited[child.PaperId]; !exists {
					queue = append(queue, child)
				}
			}
			visited[current.PaperId] = current
			cypher = "MATCH (p:PAPER {paperId: $paperId}) <- [r:REFERENCES] - (n:PAPER) RETURN n"
			cypherParam = map[string]any{
				"paperId": current.PaperId,
			}
			result, err = dbSession.Run(dbContext, cypher, cypherParam)
			PanicOnErr(err)
			linkRecord, err = neo4j.CollectWithContext(dbContext, result, err)
			PanicOnErr(err)
			for i := 0; i < len(linkRecord); i++ {
				linkNeo4jRecord, _ := linkRecord[i].Get("n")
				linkNeo4jProps := linkNeo4jRecord.(dbtype.Node).Props
				parent := Node{PaperId: linkNeo4jProps["paperId"].(string),
					CitationCount: linkNeo4jProps["citationCount"].(int64),
					Year:          linkNeo4jProps["year"].(int64),
					Title:         linkNeo4jProps["title"].(string)}
				if _, exists := visited[parent.PaperId]; !exists {
					queue = append(queue, parent)
				}
			}
		}

		depth -= 1
	}
	return visited
}
