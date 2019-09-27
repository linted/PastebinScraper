package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var queueSize = 20

func waitForInevitableHeatDeathOfTheUniverse() { //Or atleast until we receive a signal
	log.Print("The heat death is coming. I can feel it in my bones!\n")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	log.Print("Received interupt signal, I guess we will stop waiting.\n")
	return
}

func parse(matches chan pasteMatch) {
	//test code start
	log.Print("Started parsing")
	for m := range matches {
		for _, match := range m.matches {
			log.Printf("Matched rule: %s\n%s matches %s\n", m.current.pasteID, m.current.title, match.Rule)
		}
	}

	//test code end
}

func main() {
	var (
		yaraRuleFiles rules
		slackURL      slackConfig
	)

	flag.Var(&yaraRuleFiles, "rule", "Add yara rule")
	flag.Var(&slackURL, "config", "config file for slack integration")
	flag.Parse()

	if len(yaraRuleFiles) == 0 {
		log.Fatal("No rules provided\n")
	} else if len(slackURL.endpointURL) == 0 {
		log.Fatal("No config file supplied\n")
	}

	scanner := compileRules(yaraRuleFiles)
	inputStream := make(chan paste, queueSize)      //queueSize items from pastebin should probably be more then enough, right?
	matchStream := make(chan pasteMatch, queueSize) //should probably match the number of inputs
	stopFlag := make(chan bool)

	log.Printf("Slack URL: %s\n", slackURL)

	go scrape(inputStream, stopFlag)

	go scanInputs(scanner, inputStream, matchStream)

	go postToSlack(matchStream, slackURL)

	waitForInevitableHeatDeathOfTheUniverse()

	stopFlag <- false
	close(stopFlag)
	close(inputStream)
	close(matchStream)
	log.Print("Leaving main.\n")
}
