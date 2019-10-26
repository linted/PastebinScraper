package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch"
)

var elasticURLs elasticConfig

type elasticConfig struct {
	endpoints []string
}

func (s *elasticConfig) String() string {
	return s.endpoints
}

func (s *elasticConfig) Set(arg string) error {
	if len(s.endpointURL) > 0 {
		return errors.New("elastic config flag already set")
	}

	configs, err := ioutil.ReadFile(arg)
	if err != nil {
		log.Panicf("failed while reading in elastic config file: %s", err)
	}

	urls := strings.Fields(configs)
	for _, elasticURL := range urls {
		tmpURL, err := url.Parse(string(elasticURL))
		if err != nil {
			return fmt.Errorf("failed while trying to parse elastic url from config file: %s", err)
		}
		s.endpoints = append(s.endpoints, tmpURL.String())
	}
	return nil
}

func registerElasticFlags() {
	flag.Var(&elasticURL, "elastic", "config file for elastic integration")
}

func validateElasticFlags() {
	if len(elasticURL.endpointURL) == 0 {
		log.Fatal("No Elastic config file supplied\n") // TODO: make not fatal
	}
}

func postToElastic(sendQueue chan pasteMatch) {
	log.Println("Starting Elastic logging")

	client, err := elasticsearch.NewClient()
	if err != nil {
		log.Fatalf("unable to create an elastic client:", err)
	}
	log.Printf("Configured urls = %q", elasticURLs.endpoints)
}
