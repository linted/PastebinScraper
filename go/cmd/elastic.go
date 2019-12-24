package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

var elasticURLs elasticConfig

type elasticConfig struct {
	endpoints []string
	apikey    string
	index     string
}

func (s *elasticConfig) String() string {
	return fmt.Sprintf("%q", s.endpoints)
}

func (s *elasticConfig) Set(arg string) error {
	// if len(s.endpointURL) > 0 {
	// 	return errors.New("elastic config flag already set")
	// }

	configs, err := ioutil.ReadFile(arg)
	if err != nil {
		log.Panicf("failed while reading in elastic config file: %s", err)
	}

	urls := strings.Fields(string(configs))
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
	flag.Var(&elasticURLs, "elastic", "config file for elastic integration")
}

func validateElasticFlags() {
	if len(elasticURLs.endpoints) == 0 {
		log.Fatal("No Elastic config file supplied\n") // TODO: make not fatal
	}
}

func postToElastic(sendQueue chan pasteMatch) {
	log.Println("Starting Elastic logging")

	log.Printf("Elastic urls = %q", elasticURLs.endpoints)
	cfg := elasticsearch.Config{Addresses: elasticURLs.endpoints}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("unable to create an elastic client:", err)
	}

	for next := range sendQueue {
		contents, err := json.Marshal(next)
		if err != nil {
			log.Printf("unable to marshal this match. ID = %s", next.current.pasteID)
		} else {
			request := esapi.IndexRequest{
				Index: "pastebin",
				Body:  bytes.NewReader(contents),
			}

			func() {
				response, err := request.Do(context.Background(), client)
				if err != nil {
					log.Printf("Error while getting elastic response for id %s\n", next.current.pasteID)
				}

				defer response.Body.Close()

				if response.IsError() {
					log.Printf("Error %s while indexing document ID=%d", response.Status(), next.current.pasteID)
				} else {
					log.Printf("Response %s while adding %s", response.Status(), next.current.pasteID)
				}
			}()

		}

	}
}
