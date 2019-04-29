package main

import (
	"log"

	"github.com/hillu/go-yara"
)

var scrapeAmount := 30
var pastebinURL := "https://pastebin.com"
var scrapePath := fmt.Sprintf("%s/api_scraping.php?limit=%d", pastebinURL, scrapeAmount) //We don't change the limit, so compile it once and be done
var fetchPath := fmt.Sprintf("%s/api_scrape_item.php") //The query string on this one changes a lot so do it as needed

func getPaste(pasteID string, queue chan []byte) {

}

func scrape(pasteQueue chan []byte) {
	log.Print("Starting to scrape!")
	
}

func parse(matches chan []yara.MatchRule) {
	//test code start
	log.Print("Started parsing")
	for m := range matches {
		for _, match := range m {
			log.Printf("Matched rule: %s\n", match.Rule)
		}
	}

	//test code end
}
