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

func main() {
	var (
		yaraRuleFiles rules
		slackURL      slackConfig
		discordConf   discordConfig
	)

	flag.Var(&yaraRuleFiles, "rule", "Add yara rule")
	flag.Var(&slackURL, "slack", "config file for slack integration")
	flag.Var(&discordConf, "discord", "config file for discord")
	flag.Parse()

	if len(yaraRuleFiles) == 0 {
		log.Fatal("No rules provided\n")
	} else if len(slackURL.endpointURL) == 0 && (len(discordConf.token) == 0 || len(discordConf.channel) == 0) {
		log.Fatal("No config file supplied\n")
	} else if len(slackURL.endpointURL) != 0 && ( len(discordConf.token) != 0 || len(discordConf.channel) != 0){
		log.Fatal("Slack or Discord. Not both, sorry!")
	}

	scanner := compileRules(yaraRuleFiles)
	inputStream := make(chan paste, queueSize)      //queueSize items from pastebin should probably be more then enough, right?
	matchStream := make(chan pasteMatch, queueSize) //should probably match the number of inputs
	stopFlag := make(chan bool)

	go scrape(inputStream, stopFlag)

	go scanInputs(scanner, inputStream, matchStream)

	if len(slackURL.endpointURL) != 0 {
		log.Printf("Slack URL: %s\n", slackURL)
		go postToSlack(matchStream, slackURL)
	} else {
		go postToDiscord(matchStream, discordConf)
	}


	waitForInevitableHeatDeathOfTheUniverse()

	stopFlag <- false
	close(stopFlag)
	close(inputStream)
	close(matchStream)
	log.Print("Leaving main.\n")
}
