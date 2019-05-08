package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type paste struct {
	pasteID  string
	title    string
	contents []byte
}

type listing struct {
	Scrape_url string `json:"scrape_url"`
	Full_url   string `json:"full_url"`
	Date       string `json:"date"`
	Key        string `json:"key"`
	Size       string `json:"size"`
	Expire     string `json:"expire"`
	Title      string `json:"title"`
	Syntax     string `json:"syntax"`
	User       string `json:"user"`
}

type listings []listing

var scrapeAmount = 30
var pastebinURL = "https://scrape.pastebin.com"
var scrapePath = fmt.Sprintf("%s/api_scraping.php?limit=%d", pastebinURL, scrapeAmount)        //We don't change the limit, so compile it once and be done
var fetchPath, fetechPathError = url.Parse(fmt.Sprintf("%s/api_scrape_item.php", pastebinURL)) //The query string on this one changes a lot so do it as needed

func getPaste(currentPaste listing, queue chan paste) {
	log.Printf("Fetching paste: %s", currentPaste.Key)
	u := url.Values{}
	u.Add("i", currentPaste.Key)
	fetchPath.RawQuery = u.Encode()

	resp, err := http.Get(fetchPath.String())
	if err != nil {
		log.Printf("Error while fetching %s: %s\n", currentPaste.Key, err)
		return
	}

	//always clean up after yourself!
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while trying to read response of %s: %s\n", currentPaste.Key, err)
		return
	}

	queue <- paste{currentPaste.Key, currentPaste.Title, contents}
	return
}

func filterRecent(recent *listings, previous *map[string]listing) *map[string]listing {
	newListings := make(map[string]listing)
	for _, newPaste := range *recent {
		//only add values that were not in the previous one
		if _, ok := (*previous)[newPaste.Key]; !ok {
			newListings[newPaste.Key] = newPaste
		}
	}
	return &newListings
}

func scrape(pasteQueue chan paste, stop chan bool) {
	if fetechPathError != nil {
		log.Panicf("Oh no! your pastebin url is messed up. check it! %s", fetechPathError)
	}
	log.Print("Starting to scrape!\n")

	recentPastes := make(map[string]listing)

foreverLoop:
	for {
		select {
		case <-stop:
			break foreverLoop //this breaks out of the for loop
		default:
			resp, err := http.Get(scrapePath)
			if err != nil {
				log.Print("Error while scraping: %s\n", err)
				continue
			}
			defer resp.Body.Close()

			unparsedListing, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error while tyring to read response body: %s\n", err)
				continue
			}

			newListing := new(listings)
			err = json.Unmarshal(unparsedListing, &newListing)
			if err != nil {
				log.Panicf("Error while parsing the json: %s\nurl = %s\ndata = %s", err, scrapePath, unparsedListing)
				continue
			}

			recentPastes = *filterRecent(newListing, &recentPastes)
			log.Println("new listings = ", recentPastes)
			for _, val := range recentPastes {
				go getPaste(val, pasteQueue)
			}

			//sleeeeeeeeep
			select {
			case <-stop:
				log.Print("Shutting down scrapper!\n")
				break foreverLoop //get out of this... lovely loop
			case <-time.After(10000 * time.Millisecond):
				continue
			}
		}
	}

	log.Print("Shutting down scrapper.\n")
}
