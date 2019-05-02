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

type listing struct {
	Scrape_url string `json:"scrape_url"`
	Full_url string `json:"full_url"`
	Date string `json:"date"`
	Key string `json:"key"`
	Size string `json:"size"`
	Expire string `json:"expire"`
	Title string `json:"title"`
	Syntax string `json:"syntax"`
	User string `json:"user"`
}

type listings []listing

var scrapeAmount := 30
var pastebinURL := "https://pastebin.com"
var scrapePath := fmt.Sprintf("%s/api_scraping.php?limit=%d", pastebinURL, scrapeAmount) //We don't change the limit, so compile it once and be done
var fetchPath, fetechPathError := url.Parse(fmt.Sprintf("%s/api_scrape_item.php")) //The query string on this one changes a lot so do it as needed

func filter()

func getPaste(pasteID string, queue chan paste) {
	log.Printf("Fetching paste: %s", pasteID)
	u := url.Values{}
	u.Add("i", pasteID)
	fetchPath.RawQuery = u.Encode()

	resp, err := http.Get(fetchPath.string())
	if err != nil {
		log.Printf("Error while fetching %s: %s\n", pasteID, err)
		return
	}

	//always clean up after yourself!
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error while trying to read response of %s: %s\n", pasteID, err)
		return
	} 
	
	queue <- new paste{pasteID, contents}
	
	return
}

func filterRecent(recent listings, previous map[string]listing) error {
	
	for _, newPaste := range recent {
		
	}

}


func scrape(pasteQueue chan paste, stop chan bool) {
	if fetechPathError != nil {
		log.Panicf("Oh no! your pastebin url is messed up. check it! %s", fetechPathError)
	}
	log.Print("Starting to scrape!\n")
	
	var recentPastes := make(map[string]listing)

	foreverLoop: for {
		select {
		case <- stop:
			break foreverLoop //this breaks out of the for loop
		default:
			resp, err := http.Get(scrapePath)
			if err != nil {
				log.Print("Error while scraping: %s\n", err)
			}
			
			newListing := new(listings)
			err := json.Unmarshal(resp, &newListing)
			if err != nil {
				log.Printf("Error while parsing the json: %s", err)
			}

			recentPastes, err := filterRecent(newListing, recentPastes)
			if err != nil {
				log.Printf("Error while compairing results\n")
			}

			for key, val := range recentPastes {
				go getPaste(key, pasteQueue)
			}
		}
	}	

	log.Print("Shutting down scrapper.\n")
}
