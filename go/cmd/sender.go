package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
	fin, err := os.Open(arg)
	if err != nil {
		return fmt.Errorf("unable to open config file: %s", err)
	}

	var slackURL [80]byte
	len, err := fin.Read(slackURL[:])
	if err != nil && len > 0 {
		return fmt.Errorf("error while getting url from config: %s", err)
	}

	tmpURL, err := url.Parse(string(slackURL[:]))
	if err != nil {
		return fmt.Errorf("error while trying to parse url from config: %s", err)
	}
	s.endpointURL = tmpURL.String()
	return nil
}

func postToSlack(sendQueue chan pasteMatch, config slackConfig) {
	log.Print("Started slackbot!\n")

	payload := map[string]string{"text": nil}

	for next := range sendQueue {
		var matchingRules []string
		for _, match := range next.matches {
			matchingRules.Append(match.Rule)
		}
		payload["text"] = fmt.Sprintf("Pastebin Match\nURL: https://pastebin.com/%s\nMatches: %s", next.current)
		http.Post(config.endpointURL, "application/json")

	}

	log.Printf("Stopped slackbot\n")
	return
}
