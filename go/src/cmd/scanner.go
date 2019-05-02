package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/hillu/go-yara"
)

type rule struct{ namespace, filename string }
type rules []rule
type pasteMatch struct{ pasteID string, matches []yara.MatchRule}

func (r *rules) Set(arg string) error {
	if len(arg) == 0 {
		return errors.New("empty rule specification")
	}
	a := strings.SplitN(arg, ":", 2)

	switch len(a) {
	case 1:
		*r = append(*r, rule{filename: a[0]})
	case 2:
		*r = append(*r, rule{namespace: a[0], filename: a[1]})
	}

	return nil
}

func (r *rules) String() string {
	var s string
	for _, rule := range *r {
		if len(s) > 0 {
			s += " "
		}
		if rule.namespace != "" {
			s += rule.namespace + ":"
		}
		s += rule.filename
	}
	return s
}

func compileRules(files rules) *yara.Rules {
	compiler, err := yara.NewCompiler()
	if err != nil {
		log.Fatal("Unable to instantiate yara compiler\n")
	}

	for _, ruleFile := range files {
		f, err := os.Open(ruleFile.filename)
		if err != nil {
			log.Fatalf("Could not open rule file %s: %s", ruleFile.filename, err)
		}
		err = compiler.AddFile(f, ruleFile.namespace)
		if err != nil {
			log.Fatalf("Could not parse rule file %s: %s", ruleFile.filename, err)
		}
	}

	scanner, err := compiler.GetRules()
	if err != nil {
		log.Fatalf("Unable to compile rules: %s", err)
	}

	return scanner
}

func scanInputs(ruleSet *yara.Rules, inputs chan paste, results chan pasteMatch) {
	log.Print("Starting to scan inputs\n")
	for target := range inputs {
		matches, err := ruleSet.ScanMem(target.contents, 0, 5) //TODO: figure out what the flags do
		if err != nil {
			log.Printf("Got error while scanning: %s", err)
		} else if len(matches) > 0 {
			log.Print("Found a match\n")
			results <- new pasteMatch{ target.pasteID, matches }
		} else {
			log.Print("Not a match\n")
		}
	}
}
