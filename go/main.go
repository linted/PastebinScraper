package main

import (
	"flag"
	"log"

	"github.com/hillu/go-yara"
)

type rule struct{ namespace, filename string }
type rules []rule

func main() {
	var (
		yaraRules rules
	)

	flag.Var(&yaraRules, "rule", "Add yara rule")
	flag.Parse()

	if len(yaraRules) == 0 {
		log.Fatal("No rules provided\n")
	}

	compiler, err := yara.NewCompiler()
	if err != nil {
		log.Fatal("Unable to instantiate yara compiler\n")
	}

}
