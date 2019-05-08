package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
)

type slackConfig struct {
	endpointURL *url.URL
}

func (s *slackConfig) String() string {
	return s.endpointURL.String()
}

func (s *slackConfig) Set(arg string) error {
	if s.endpointURL != nil {
		return errors.New("config flag already set")
	}
	fin, err := os.Open(arg)
	if err != nil {
		return fmt.Errorf("unable to open config file: %s", err)
	}

	var slackURL [80]byte
	len, err := fin.Read(slackURL)
	if err != nil {
		return fmt.Errorf("error while getting url from config: %s", err)
	}

	s.endpointURL, err = url.Parse(slackURL)
	if err != nil {
		return fmt.Errorf("error while trying to parse url from config: %s", err)
	}
	return nil
}

func postToSlack() {
	log.Print("Started slackbot!\n")
	return
}
