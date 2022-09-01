package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type slackConfig struct {
	endpointURL string
}

type discordConfig struct {
	token string
	channel string
}

func (s *slackConfig) String() string {
	return strings.TrimSpace(s.endpointURL)
}

func (d *discordConfig) String() string {
	return strings.TrimSpace(d.token) + strings.TrimSpace(d.channel)
}

func (s *slackConfig) Set(arg string) error {
	if len(s.endpointURL) > 0 {
		return errors.New("config flag already set")
	}

	slackURL, err := ioutil.ReadFile(arg)
	if err != nil {
		log.Panicf("Error while reading in config: %s", err)
	}

	tmpURL, err := url.Parse(string(slackURL))
	if err != nil {
		return fmt.Errorf("error while trying to parse url from config: %s", err)
	}
	s.endpointURL = tmpURL.String()
	return nil
}

func (d *discordConfig) Set(arg string) error {
	if len(d.token) > 0 {
		return errors.New("config flag already set")
	}

	discordConfig, err := ioutil.ReadFile(arg)
	if err != nil {
		log.Panicf("Error while reading in config: %s", err)
	}

	var discordToken  = discordConfig[:72]
	var discordChannel = discordConfig[73:] // skip a newline

	if len(discordToken) != 72 {
		log.Panicf("Invalid discord token length!")
	}
	if len(discordChannel) != 18 {
		log.Panicf("Invalid discord channel length!")
	}

	d.token = string(discordToken)
	d.channel = string(discordChannel)
	return nil
}

func postToSlack(sendQueue chan pasteMatch, config slackConfig) {
	log.Print("Started slackbot!\n")

	payload := map[string]string{"text": ""}

	for next := range sendQueue {
		var matchingRules string
		for _, match := range next.matches {
			matchingRules += match.Rule + " "
			// for _, matchString := range match.Strings {
			// 	log.Print(string(matchString.Data))
			// }
		}
		payload["text"] = fmt.Sprintf("Pastebin Match\nURL: https://pastebin.com/%s\nTitle: %s\nMatches: %s", next.current.pasteID, next.current.title, matchingRules)

		contents, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error while marshaling contents: %s", err)
			continue
		}
		log.Printf("Sending message: %s", contents)
		resp, err := http.Post(config.endpointURL, "application/json", bytes.NewBuffer(contents))
		if err != nil {
			log.Printf("Error while sending! %s", err)
			continue
		}
		log.Printf("Resp = %s", resp.Status)
		resp.Body.Close()
	}

	log.Printf("Stopped slackbot\n")
	return
}


func postToDiscord(sendQueue chan pasteMatch, config discordConfig) {
	log.Print("Started Discord bot!\n")

	dg, err := discordgo.New("Bot " + config.token)
	if err != nil {
		log.Panicf("Unable to create discord bot: %s", err)
	}

	for next := range sendQueue {
		var matchingRules string

		for _,match := range next.matches {
			matchingRules += "[" + match.Rule + "]"
		}

		message = fmt.Sprintf("%s %s: https://pastebin.com/%s", matchingRules, next.current.title, next.current.pasteID)

		_, err = dg.ChannelMessageSend(config.channel, message)
		if err != nil {
			log.Printf("Error while sending! %s\n", err)
		}

	}

	log.Printf("Stopped Discord bot\n")
	return
}