package main

import (
	"log"
	"os"

	memesapi "gitlab.com/shitposting/memesapi/rest/client"
)

var (
	//Telegram
	telegramToken string

	//Memesapi
	apiEndpoint string
	apiPlatform string

	//Ratelimiter
	maxActionsPerMinute      string
	maxGroupActionsPerMinute string

	//err is declared here for functions that return an error as the second value
	err error

	mClient *memesapi.Client
)

func envSetup() error {
	var ok bool

	telegramToken, ok = os.LookupEnv("TELEGRAM_TOKEN")
	if telegramToken == "" || !ok {
		log.Fatalf("telegram token bot is not optional!")
	}

	apiEndpoint, ok = os.LookupEnv("API_ENDPOINT")
	if apiEndpoint == "" || !ok {
		apiEndpoint = "http://127.0.0.1:34378"
	}

	apiPlatform, ok = os.LookupEnv("API_PLATFORM")
	if apiPlatform == "" || !ok {
		apiPlatform = "tgrandombot"
	}

	maxActionsPerMinute, ok = os.LookupEnv("MAX_PRIVATE_ACTIONS")
	if maxActionsPerMinute == "" || !ok {
		maxActionsPerMinute = "15"
	}

	maxGroupActionsPerMinute, ok = os.LookupEnv("MAX_GROUP_ACTIONS")
	if maxGroupActionsPerMinute == "" || !ok {
		maxGroupActionsPerMinute = "10"
	}

	return nil
}
