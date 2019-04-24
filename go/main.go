package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hillu/go-yara"
)

type rule struct{ namespace, filename string }
type rules []rule

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

	for _, ruleFile := range yaraRules {
		f, err := os.Open(ruleFile.filename)
		if err != nil {
			log.Fatalf("Could not open rule file %s: %s", ruleFile.filename, err)
		}
		err = compiler.AddFile(f, ruleFile.namespace)
		if err != nil {
			log.Fatalf("Could not parse rule file %s: %s", ruleFile.filename, err)
		}
	}

	fmt.Print("Everything works up to here!\n")
}
