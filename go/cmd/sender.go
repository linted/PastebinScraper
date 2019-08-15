package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type slackConfig struct {
	endpointURL string
}

func (s *slackConfig) String() string {
	return s.endpointURL
}

func (s *slackConfig) Set(arg string) error {
	if len(s.endpointURL) > 0 {
		return errors.New("config flag already set")
	}

	slackURL, err := ioutil.ReadFile(arg)
	if err != nil {
		log.Panicf("Error while reading in config: %s", err)
	}

	tmpURL, err := url.Parse(string(slackURL))
	if err != nil {
		return fmt.Errorf("error while trying to parse url from config: %s", err)
	}
	s.endpointURL = tmpURL.String()
	return nil
}

func postToSlack(sendQueue chan pasteMatch, config slackConfig) {
	log.Print("Started slackbot!\n")

	payload := map[string]string{"text": ""}

	for next := range sendQueue {
		var matchingRules string
		for _, match := range next.matches {
			matchingRules += match.Rule + " "
			// for _, matchString := range match.Strings {
			// 	log.Print(string(matchString.Data))
			// }
		}
		payload["text"] = fmt.Sprintf("Pastebin Match\nURL: https://pastebin.com/%s\nTitle: %s\nMatches: %s", next.current.pasteID, next.current.title, matchingRules)

		contents, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error while marshaling contents: %s", err)
			continue
		}
		log.Printf("Sending message: %s", contents)
		resp, err := http.Post(config.endpointURL, "application/json", bytes.NewBuffer(contents))
		if err != nil {
			log.Printf("Error while sending! %s", err)
			continue
		}
		log.Printf("Resp = %s", resp.Status)
	}

	log.Printf("Stopped slackbot\n")
	return
}
