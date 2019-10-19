package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var queueSize = 20
var stopFlag = make(chan bool)

func waitForInevitableHeatDeathOfTheUniverse() { //Or atleast until we receive a signal
	log.Print("The heat death is coming. I can feel it in my bones!\n")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	log.Print("Received interupt signal, I guess we will stop waiting.\n")
	return
}

func main() {
	var (
		yaraRuleFiles rules
	)

	flag.Var(&yaraRuleFiles, "rule", "Add yara rule")
	registerSenderFlags()
	flag.Parse()

	if len(yaraRuleFiles) == 0 {
		log.Fatal("No rules provided\n")
	}
	validateSenderFlags()

	scanner := compileRules(yaraRuleFiles)
	inputStream := make(chan paste, queueSize)      //queueSize items from pastebin should probably be more then enough, right?
	matchStream := make(chan pasteMatch, queueSize) //should probably match the number of inputs

	log.Printf("Slack URL: %s\n", slackURL)

	go scrape(inputStream, stopFlag)

	go scanInputs(scanner, inputStream, matchStream)

	go startSending(matchStream)

	waitForInevitableHeatDeathOfTheUniverse()

	stopFlag <- false
	close(stopFlag)
	close(inputStream)
	close(matchStream)
	log.Print("Leaving main.\n")
}
