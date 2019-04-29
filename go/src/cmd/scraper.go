package main

import (
	"log"

	"github.com/hillu/go-yara"
)

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
