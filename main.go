package main

import (
	"fmt"
	"log"
	"strconv"

	memesapi "github.com/shitpostingio/randomapi/rest/client"

	limiter "github.com/shitpostingio/telegramrandombot/ratelimiter"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/telegramrandombot/groups"
	"github.com/shitpostingio/telegramrandombot/private"
)

var (
	//Build version of the bot, a compile-time value
	Build string

	//Version of the bot, a compile-time value
	Version string

	//configFilePath allows the user to pass a custom file path for the configuration file
	configFilePath string

	//debug allows the user to set the bot in debug mode
	debug bool
)

func main() {

	envSetup()

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := memesapi.New(apiEndpoint, apiPlatform)

	bot.Debug = debug
	log.Println(fmt.Sprintf("Shitposting tg-random-bot version %s, build %s", Version, Build))
	log.Println(fmt.Sprintf("Authorized on account @%s", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	privateActions, err := strconv.ParseUint(maxActionsPerMinute, 10, 0)
	if err != nil {
		log.Fatal("cannot parse max private actions")
	}

	groupActions, err := strconv.ParseUint(maxGroupActionsPerMinute, 10, 0)
	if err != nil {
		log.Fatal("cannot parse max group actions")
	}

	limiter.StartRateLimiter(bot, uint(privateActions))
	go groups.StartGroupRateLimiter(uint(groupActions))

	for update := range updates {
		if update.Message != nil {
			go mainBot(bot, update.Message, client)
		}
	}
}

func mainBot(bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *memesapi.Client) {

	switch {
	case message.Chat.IsPrivate():
		private.HandlePrivate(message, client, bot)
	case message.IsCommand():
		groups.HandleCommands(message, client, bot)
	}
}
