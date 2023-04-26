package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// BuildGraphHandler handler func to build graph for new seed paper
func BuildGraphHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		graph := BuildGraph(paperId, ctx)
		ctx.IndentedJSON(http.StatusOK, graph)
	}
}

// BuildGraph builds a graph using bfs algo for a new seed paper and returns the graph
func BuildGraph(seedPaperId string, ctx *gin.Context) map[string]Node {
	graph := BfsBuilder(seedPaperId)
	log.Printf("Build graph Paper details : %v", graph)
	persistGraphToDB(graph, ctx)
	return graph
}

// BfsBuilder implementation of the bfs algo that builds the graph
func BfsBuilder(seedPaperId string) map[string]Node {
	depth := viper.GetInt("graph.depth")
	breadth := viper.GetInt("graph.refBreadth")
	var nodeRef []NodeReferences
	var nodes []Node
	n := GetPaperNode(seedPaperId)
	queue := []Node{n}
	visited := map[string]Node{}

	for len(queue) > 0 && depth > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			current := queue[0]
			queue = queue[1:]
			nodeRef = GetReferences(current.PaperId, true, breadth)
			nodes = GetNodeReferences(nodeRef, true)
			for _, child := range nodes {
				if child.PaperId != "" {
					current.Reference = append(current.Reference, child)
				}
			}
			for _, child := range nodes {
				if _, exists := visited[child.PaperId]; !exists && child.PaperId != "" {
					queue = append(queue, child)
				}
			}
			nodeRef = GetReferences(current.PaperId, false, breadth)
			nodes = GetNodeReferences(nodeRef, false)
			for _, child := range nodes {
				if _, exists := visited[child.PaperId]; !exists && child.PaperId != "" {
					child.Reference = append(child.Reference, current)
					queue = append(queue, child)
				}
			}
			visited[current.PaperId] = current
		}
		depth -= 1
		// if breadth > 1 {
		// 	breadth -= 1
		// }

	}
	return visited
}

// GetPaperNode function to seed paper node by getting details from S2AG API
func GetPaperNode(paperId string) Node {
	url := viper.Get("s2ag.urlRoot").(string) + paperId + viper.Get("s2ag.paperUrlFields").(string)
	var node Node
	err := json.Unmarshal(FetchFromS2ag(url), &node)
	PanicOnErr(err)
	return node
}

// GetReferences function to get all the references and citations for a paperId
func GetReferences(paperId string, reference bool, breadth int) []NodeReferences {
	url := viper.GetString("s2ag.urlRoot") + paperId
	if reference {
		url += viper.GetString("s2ag.referenceUrlFields")
	} else {
		url += viper.GetString("s2ag.citationUrlFields")
	}
	url += viper.GetString("s2ag.buffer")
	var resp ResponseReferences
	err := json.Unmarshal(FetchFromS2ag(url), &resp)
	PanicOnErr(err)
	// return []NodeReferences{}
	if len(resp.Data) <= breadth {
		return resp.Data
	} else {
		sort.Slice(resp.Data, func(i, j int) bool {
			if resp.Data[i].IsInfluential {
				return true
			} else {
				return false
			}
		})
		return resp.Data[:breadth]
	}
}

// GetNodeReferences function to get all the citation or reference node
func GetNodeReferences(nodeReferences []NodeReferences, reference bool) []Node {
	var nodes []Node
	for _, child := range nodeReferences {
		paperId := ""
		if reference {
			paperId = child.CitedPaper.RefPaperId
		} else {
			paperId = child.CitingPaper.CitPaperId
		}
		nodes = append(nodes, GetPaperNode(paperId))
	}
	return nodes
}

// persistGraphToDB writes the graph to neo4j db
func persistGraphToDB(graph map[string]Node, ctx *gin.Context) {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	for paperId, node := range graph {
		log.Printf("Persisting graph for seed paper : %s", paperId)
		checkAndCreateNode(node, dbContext, dbSession)

		// iterate through references and attach the relations
		if node.Reference != nil {
			for _, childNode := range node.Reference {
				checkAndCreateNode(childNode, dbContext, dbSession)
				checkAndCreateRelation(node, childNode, dbContext, dbSession)
			}
		}
	}
}

// checkAndCreateNode create neo4j paper node if it does not exist in db
func checkAndCreateNode(node Node, dbContext context.Context, dbSession neo4j.SessionWithContext) {
	cypher := "MERGE (p: PAPER {paperId: $paperId}) SET p = {paperId: $paperId, citationCount: $citationCount, title: $title, year: $year} return p"
	cypherParam := map[string]any{
		"paperId":       node.PaperId,
		"citationCount": node.CitationCount,
		"title":         node.Title,
		"year":          node.Year,
	}
	_, err := dbSession.Run(dbContext, cypher, cypherParam)
	PanicOnErr(err)
	log.Printf("Created node for Paper: %s", node.PaperId)
}

// checkAndCreateRelation adds reference relation between the new nodes created
func checkAndCreateRelation(node Node, childNode Node, dbContext context.Context, dbSession neo4j.SessionWithContext) {
	cypher := "MATCH (p:PAPER {paperId: $paperId}), (q:PAPER {paperId: $childPaperId}) MERGE (p)-[r: REFERENCES]->(q)"
	cypherParam := map[string]any{
		"paperId":      node.PaperId,
		"childPaperId": childNode.PaperId,
	}
	_, err := dbSession.Run(dbContext, cypher, cypherParam)
	PanicOnErr(err)
	log.Printf("Attached reference relationship paper{%s} -> paper{%s}", node.PaperId, childNode.PaperId)
}
