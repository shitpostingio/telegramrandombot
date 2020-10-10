package private

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/shitposting/loglog-ng"

	limiter "gitlab.com/shitposting/tg-random-bot/ratelimiter"
	"gitlab.com/shitposting/tg-random-bot/utility"

	memesapi "gitlab.com/shitposting/memesapi/rest/client"
)

const welcomeMessage = "Tap the buttons to get a recent or a random meme from @shitpost\n\nUse /start or /help if it disappears"

//HandlePrivate handles commands in private chats
func HandlePrivate(message *tgbotapi.Message, client *memesapi.Client, bot *tgbotapi.BotAPI) {

	if message.IsCommand() {

		command := strings.ToLower(message.Command())

		switch command {
		case "start", "help":
			handleStart(message)
		case "random":
			requestMeme(false, message, client, bot)
		case "recent":
			requestMeme(true, message, client, bot)
		}

		return
	}

	switch message.Text {
	case "Random Meme 🔀":
		requestMeme(false, message, client, bot)
	case "Recent Meme ⏰":
		requestMeme(true, message, client, bot)
	}
}

func handleStart(message *tgbotapi.Message) {

	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Recent Meme ⏰"), tgbotapi.NewKeyboardButton("Random Meme 🔀")))
	keyboard.OneTimeKeyboard = false

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeMessage)
	msg.ReplyMarkup = keyboard
	if _, _, err := limiter.Send(msg); err != nil {
		loglog.Err(err.Error())
	}
}

func requestMeme(recent bool, message *tgbotapi.Message, client *memesapi.Client, bot *tgbotapi.BotAPI) {
	utility.TrySending(recent, message.Chat.ID, message.From.ID, client, bot)
}
