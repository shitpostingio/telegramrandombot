package groups

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/shitposting/loglog-ng"

	limiter "gitlab.com/shitposting/tg-random-bot/ratelimiter"
	"gitlab.com/shitposting/tg-random-bot/utility"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	memesapi "gitlab.com/shitposting/memesapi/rest/client"

	"github.com/patrickmn/go-cache"
)

const (
	groupRoutineLifespan  = time.Minute
	groupEntryExpiration  = 5 * time.Minute
	pleaseSlowDownMessage = "🚦Please ask for memes more slowly, so I don't get limited by Telegram!🚦\n\nTry again in a minute or write me in private!"
	groupSlowDownReport   = "The group %s (handle %s, id %d) has requested over %d memes in the last minute and has been asked to slow down"
)

var (
	requests                  *cache.Cache
	maxGroupRequestsPerMinute uint
)

// StartGroupRateLimiter limits the requests per minute a group can perform
func StartGroupRateLimiter(maxRequestsPerMinute uint) {
	maxGroupRequestsPerMinute = maxRequestsPerMinute
	requests = cache.New(groupRoutineLifespan, groupEntryExpiration)
}

// HandleCommands handles commands in groups
func HandleCommands(message *tgbotapi.Message, client *memesapi.Client, bot *tgbotapi.BotAPI) {
	command := strings.ToLower(message.Command())
	switch command {
	case "random":
		authorizeMeme(false, message, client, bot)
	case "recent":
		authorizeMeme(true, message, client, bot)
	}
}

func authorizeMeme(recent bool, message *tgbotapi.Message, client *memesapi.Client, bot *tgbotapi.BotAPI) {

	groupKey := strconv.FormatInt(message.Chat.ID, 10)
	_, found := requests.Get(groupKey)
	if !found {
		err := requests.Add(groupKey, uint(0), groupRoutineLifespan)
		if err != nil {
			loglog.Warn(fmt.Sprintf("Unable to add group with ID %s to the request cache", groupKey))
		}
	}

	groupRequests, err := requests.IncrementUint(groupKey, 1)
	if err != nil {
		loglog.Warn(fmt.Sprintf("Unable to increment request count for group with ID %s", groupKey))
	}

	if groupRequests > maxGroupRequestsPerMinute {

		if groupRequests == maxGroupRequestsPerMinute+1 {
			askToSlowDown(message)
		}

		return
	}

	sendMemeToGroup(recent, message, client, bot)
}

func sendMemeToGroup(recent bool, message *tgbotapi.Message, client *memesapi.Client, bot *tgbotapi.BotAPI) {
	utility.TrySending(recent, message.Chat.ID, message.From.ID, client, bot)
}

func askToSlowDown(message *tgbotapi.Message) {
	loglog.Warn(fmt.Sprintf(groupSlowDownReport, message.Chat.Title, message.Chat.UserName, message.Chat.ID, maxGroupRequestsPerMinute))

	_, _, err := limiter.Send(tgbotapi.NewMessage(message.Chat.ID, pleaseSlowDownMessage))
	if err != nil {
		loglog.Err(fmt.Sprintf("Unable to send slow down message to group with ID %d", message.Chat.ID))
	}
}
