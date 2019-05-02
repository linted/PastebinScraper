package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/hillu/go-yara"
)

var QUEUESIZE := 20

func waitForInevitableHeatDeathOfTheUniverse() { //Or atleast until we receive a signal
	log.Print("The heat death is coming. I can feel it in my bones!\n")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	log.Print("Received interupt signal, I guess we will stop waiting.\n")
	return
}

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
	inputStream := make(chan []byte, QUEUESIZE)           //queuesize items from pastebin should probably be more then enough, right?
	matchStream := make(chan []yara.MatchRule, QUEUESIZE) //should probably match the number of inputs

	log.Print("Everything works up to here!\n")

	go scanInputs(scanner, inputStream, matchStream)

	testInput1 := []byte("email@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\nemail@email.email:password\r\n")
	inputStream <- testInput1
	testInput2 := []byte("yolo swag bro\n")
	inputStream <- testInput2

	go parse(matchStream)

	waitForInevitableHeatDeathOfTheUniverse()

	close(inputStream)
	close(matchStream)
	log.Print("Leaving main.\n")
}
