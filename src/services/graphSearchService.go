package services

import (
	"encoding/json"
	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strconv"
)

// IsSeedPaperHandler handler func to check if you give paperId is seed or not
func IsSeedPaperHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		log.Printf("Seed Paper handler URL Params: %v", paperId)
		isSeed := isSeedPaper(paperId, ctx)
		ctx.IndentedJSON(http.StatusOK, isSeed)
	}
}

// GraphHandler handler func for getting the graph for a given seed paper Id
func GraphHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		log.Printf("graph Paper handler URL Params: %v", paperId)

		// check if it is a seed paper
		isSeed := isSeedPaper(paperId, ctx)

		if !isSeed {
			BuildGraph(paperId, ctx)
			addSeedRelation(paperId, ctx)
		}
		seedNode := getNode(paperId, ctx)
		visited := bfs(seedNode, 3, ctx)
		transformedGraph := transformGraphToGS(visited, paperId)
		ctx.IndentedJSON(http.StatusOK, transformedGraph)
	}
}

// FilteredGraphHandler handler func that returns filtered graph to ui
func FilteredGraphHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		authorFilter := ctx.QueryArray("authors")
		yearFilter := ctx.QueryArray("year")
		minCitationFilter, err := strconv.ParseInt(ctx.Query("minCitation"), 10, 64)
		PanicOnErr(err)

		// get complete graph
		var graph map[string]any
		graphRespBody := FetchFromGraphService("/graphSearch/graph/" + paperId)
		err = json.Unmarshal(graphRespBody, &graph)
		PanicOnErr(err)

		// filter
		filteredNodes := map[string]any{}

		for _, node := range graph["nodes"].([]interface{}) {
			keepNode := true
			// author filter
			keepAuthorNode := false
			node := node.(map[string]any)
			for _, author := range authorFilter {
				if slices.Contains(GetStringList(node["authorList"]), author) {
					keepAuthorNode = true
					break
				}
			}
			keepNode = keepAuthorNode

			// year, minCitationCount
			paperYear := int64(node["year"].(float64))
			paperCitationCount := int64(node["citationCount"].(float64))
			minYear, _ := strconv.ParseInt(yearFilter[0], 10, 64)
			maxYear, _ := strconv.ParseInt(yearFilter[1], 10, 64)
			if paperYear <= minYear && paperYear >= maxYear {
				keepNode = false
			} else if paperCitationCount < minCitationFilter {
				keepNode = false
			}

			if keepNode {
				filteredNodes[node["id"].(string)] = node
			}

		}

		// edge filter
		filteredEdges := []map[string]any{}
		for _, edges := range graph["edges"].([]interface{}) {
			edges := edges.(map[string]any)
			_, existSource := filteredNodes[edges["source"].(string)]
			_, existTarget := filteredNodes[edges["target"].(string)]
			if existSource && existTarget {
				filteredEdges = append(filteredEdges, edges)
			}
		}

		nodeList := []any{}
		for _, node := range filteredNodes {
			nodeList = append(nodeList, node)
		}
		filterGraph := map[string]any{
			"nodes": nodeList,
			"edges": filteredEdges,
		}

		ctx.IndentedJSON(http.StatusOK, filterGraph)
	}
}

// isSeedPaper checks is given paper-id has been already a seed
func isSeedPaper(paperId string, ctx *gin.Context) bool {
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

// addSeedRelation adds seed relation to the new seed paper node
func addSeedRelation(seedPaperId string, ctx *gin.Context) {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	cypher := "MATCH (n: PAPER {paperId: $paperId}), (s: SEED_PAPER) MERGE (s) - [r:SEED] -> (n) RETURN r "
	cypherParam := map[string]any{
		"paperId": seedPaperId,
	}
	_, err := dbSession.Run(dbContext, cypher, cypherParam)
	PanicOnErr(err)
}

// getNode gets the seed node with the paperId
func getNode(rootPaperId string, ctx *gin.Context) Node {
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
		Title:         nodeNeo4jProps["title"].(string),
		AuthorsList:   GetStringList(nodeNeo4jProps["authorList"]),
	}
	return n
}

// bfs performs bfs on the seed paper node and returns the graph
func bfs(n Node, depth int, ctx *gin.Context) map[string]Node {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	queue := []Node{n}
	visited := map[string]Node{}

	for len(queue) > 0 && depth > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
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
					Title:         linkNeo4jProps["title"].(string),
					AuthorsList:   GetStringList(linkNeo4jProps["authorList"]),
				}
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
					Title:         linkNeo4jProps["title"].(string),
					AuthorsList:   GetStringList(linkNeo4jProps["authorList"]),
				}
				if _, exists := visited[parent.PaperId]; !exists {
					queue = append(queue, parent)
				}
			}
		}

		depth -= 1
	}
	return visited
}

func createG6Node(paperDetails Node, seedPaperId string) map[string]any {
	donutAttr := map[string]int{}

	if paperDetails.PaperId == seedPaperId {
		donutAttr["seed"] = 100
		donutAttr["notSeed"] = 0
	} else {
		donutAttr["seed"] = 0
		donutAttr["notSeed"] = 100
	}
	return map[string]any{
		"id":            paperDetails.PaperId,
		"label":         paperDetails.Title,
		"donutAttrs":    donutAttr,
		"year":          paperDetails.Year,
		"authorList":    paperDetails.AuthorsList,
		"citationCount": paperDetails.CitationCount,
	}
}
func transformGraphToGS(graph map[string]Node, seedPaperId string) map[string]any {
	nodes := map[string]any{}
	var edges []map[string]any

	for paperId, paper := range graph {
		if _, exists := nodes[paperId]; !exists {
			nodes[paperId] = createG6Node(paper, seedPaperId)
		}
		for _, childPaper := range paper.Reference {
			if _, exists := nodes[childPaper.PaperId]; !exists {
				nodes[childPaper.PaperId] = createG6Node(childPaper, seedPaperId)
			}
			edgeDetails := map[string]any{
				"source": paperId,
				"target": childPaper.PaperId,
				"size":   2,
			}
			edges = append(edges, edgeDetails)
		}
	}

	// add to node slice
	var nodeList []any
	for _, node := range nodes {
		nodeList = append(nodeList, node)
	}
	return map[string]any{
		"nodes": nodeList,
		"edges": edges,
	}
}
