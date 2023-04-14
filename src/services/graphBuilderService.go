package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func BuildGraphHndler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		BuildGraph(paperId, ctx)
	}
}

func BuildGraph(seedPaperId string, ctx *gin.Context) {
	node := GetPaperDetails((seedPaperId))
	log.Printf("Build graph Paper details : %v", node)
	ctx.IndentedJSON(http.StatusOK, node)
}

func GetPaperDetails(paperId string) Node {
	url := viper.Get("s2ag.urlRoot").(string) + paperId + viper.Get("s2ag.paperUrlFields").(string)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", os.Getenv("S2AG_KEY"))
	resp, err := client.Do(req)
	PanicOnErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	PanicOnErr(err)
	var node Node
	err = json.Unmarshal(body, &node)
	PanicOnErr(err)
	return node
}
