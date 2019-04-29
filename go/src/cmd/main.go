package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var (
		yaraRuleFiles rules
	)

	flag.Var(&yaraRuleFiles, "rule", "Add yara rule")
	flag.Parse()

	if len(yaraRuleFiles) == 0 {
		log.Fatal("No rules provided\n")
	}

	scanner := compileRules(yaraRuleFiles)

	fmt.Print("Everything works up to here!\n")

	matches, err := scanner.ScanProc(10342, 0, 0)
	if err != nil {
		log.Fatalf("Something went wrong while scanning the proc: %s", err)
	}
	log.Printf("Results = %d\n", len(matches))
	for _, match := range matches {
		log.Printf("+ [%s] %s", match.Namespace, match.Rule)
	}

}
