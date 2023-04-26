package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"

	// . "github.com/ArxivInsanity/graph-service/src/util"
	// "github.com/gin-gonic/gin"
	// "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Paper struct {
	ID                      string `json:"id"`
	URL                     string `json:"url"`
	Year                    int    `json:"year"`
	Authors                 []struct {
		Name string `json:"name"`
	} `json:"authors"`
	InfluentialCitationCount int `json:"influentialCitationCount"`
}

func nodeDetails(ctx context.Context, paperTitle string) (string, int, string, bool) {
	var paper Paper
	url := fmt.Sprintf("https://api.semanticscholar.org/graph/v1/paper/autocomplete?query=%s", strings.ReplaceAll(paperTitle, " ", "%20"))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error retrieving paper metadata.")
		return "", 0, "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Error retrieving paper metadata.")
		return "", 0, "", false
	}

	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)
	matches := respData["matches"].([]interface{})
	if len(matches) > 0 {
		match := matches[0].(map[string]interface{})
		paperID := match["id"].(string)
		url := fmt.Sprintf("https://api.semanticscholar.org/graph/v1/paper/%s?fields=url,year,authors,influentialCitationCount", paperID)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error retrieving paper metadata.")
			return "", 0, "", false
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println("Error retrieving paper metadata.")
			return "", 0, "", false
		}

		json.NewDecoder(resp.Body).Decode(&paper)

		var authorList []string
		for _, author := range paper.Authors {
			authorList = append(authorList, author.Name)
		}

		return paper.ID, paper.Year, strings.Join(authorList, ", "), paper.InfluentialCitationCount > 100
	}

	return "", 0, "", false
}
