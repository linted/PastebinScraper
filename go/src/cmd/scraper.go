package main

import (
	"log"
	"net/http"
	"net/url"
	"github.com/hillu/go-yara"
	"fmt"
	"io/ioutil"
)

type paste struct { 
	pasteID string,
	contents []byte
 }

var scrapeAmount := 30
var pastebinURL := "https://pastebin.com"
var scrapePath := fmt.Sprintf("%s/api_scraping.php?limit=%d", pastebinURL, scrapeAmount) //We don't change the limit, so compile it once and be done
var fetchPath, fetechPathError := url.Parse(fmt.Sprintf("%s/api_scrape_item.php")) //The query string on this one changes a lot so do it as needed

func getPaste(pasteID string, queue chan paste) {
	log.Printf("Fetching paste: %s", pasteID)
	u := url.Values{}
	u.Add("i", pasteID)
	fetchPath.RawQuery = u.Encode()

	resp, err := http.Get(fetchPath.string())
	if err != nil {
		log.Printf('Error while fetching %s: %s', pasteID, err)
		return
	}

	//always clean up after yourself!
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error while trying to read response of %s: %s", pasteID, err)
		return
	} 
	
	queue <- new paste{pasteID, contents}
	
	return
}

func scrape(pasteQueue chan []byte) {
	if fetechPathError != nil {
		log.Panicf("Oh no! your pastebin url is messed up. check it! %s", fetechPathError)
	}
	log.Print("Starting to scrape!")
	
}
