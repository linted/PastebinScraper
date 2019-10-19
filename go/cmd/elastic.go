package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

type elasticConfig struct {
	endpointURL string
}

func (s *elasticConfig) String() string {
	return s.endpointURL
}

func (s *elasticConfig) Set(arg string) error {
	if len(s.endpointURL) > 0 {
		return errors.New("elastic config flag already set")
	}

	elasticURL, err := ioutil.ReadFile(arg)
	if err != nil {
		log.Panicf("failed while reading in elastic config: %s", err)
	}

	tmpURL, err := url.Parse(string(elasticURL))
	if err != nil {
		return fmt.Errorf("failed while trying to parse elastic url from config file: %s", err)
	}
	s.endpointURL = tmpURL.String()
	return nil
}
