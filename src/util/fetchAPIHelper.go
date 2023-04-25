package util

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func FetchFromS2ag(url string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", os.Getenv("S2AG_KEY"))
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	PanicOnErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	PanicOnErr(err)
	log.Printf("Fetching %s in %v", url, duration)
	return body
}

func FetchFromGraphService(url string) []byte {
	client := &http.Client{}

	// append url to baseurl
	absoluteUrl := viper.GetString("service.baseUrl") + url

	req, _ := http.NewRequest("GET", absoluteUrl, nil)
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	PanicOnErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	PanicOnErr(err)
	log.Printf("Fetching %s in %v", url, duration)
	return body
}
