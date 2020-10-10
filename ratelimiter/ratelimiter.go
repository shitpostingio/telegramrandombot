package limiter

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot                 *tgbotapi.BotAPI
	actions             uint
	maxActionsPerSecond uint
	mutex               sync.Mutex
	channels            rateLimiterChannels
	rateResetChannel    chan bool
)

//StartRateLimiter starts the rate limiter
func StartRateLimiter(botAPI *tgbotapi.BotAPI, maxActionsToPerformEverySecond uint) {

	/* SET VARIABLES */
	bot = botAPI
	maxActionsPerSecond = maxActionsToPerformEverySecond

	/* MAKE CHANNELS */
	rateResetChannel = make(chan bool)
	channels.Send = make(chan actionRequest)
	channels.SendUrgent = make(chan actionRequest)
	channels.Request = make(chan actionRequest)
	channels.RequestUrgent = make(chan actionRequest)
	channels.FilePathRequest = make(chan filePathRequest)

	/* START RATE LIMITER */
	go limitRates()

	/* HANDLE REQUESTS */
	go handleRequests()
}

func limitRates() {
	timeToWait := 1 * time.Second
	for {
		time.Sleep(timeToWait)
		mutex.Lock()
		actions = 0
		select {
		case rateResetChannel <- true:
		default:
		}
		mutex.Unlock()
	}
}

func increaseActions() {
	mutex.Lock()
	actions++
	mutex.Unlock()
}

func canExecuteAction() bool {
	mutex.Lock()
	result := actions < maxActionsPerSecond
	mutex.Unlock()
	return result
}

func handleRequests() {
	for {
		if !canExecuteAction() {
			<-rateResetChannel
		}

		increaseActions()
		select {
		case urgentSend := <-channels.SendUrgent:
			go send(urgentSend)
		case urgentRequest := <-channels.RequestUrgent:
			go request(urgentRequest)
		default:
			waitForNextAction()
		}
	}
}

func waitForNextAction() {
	select {
	case urgentSend := <-channels.SendUrgent:
		go send(urgentSend)
	case urgentRequest := <-channels.RequestUrgent:
		go request(urgentRequest)
	case normalSend := <-channels.Send:
		go send(normalSend)
	case normalRequest := <-channels.Request:
		go request(normalRequest)
	case filepathRequest := <-channels.FilePathRequest:
		go getFilePath(filepathRequest)
	}
}

func request(actionRequest actionRequest) {
	result, err := bot.Request(actionRequest.Action)
	actionRequest.ResultChannel <- actionResult{RequestResult: &result, Error: err}
}

func send(actionRequest actionRequest) {
	result, err := bot.Send(actionRequest.Action)
	actionRequest.ResultChannel <- actionResult{SendResult: &result, Error: err}
}

//noinspection GoNilness
func getFilePath(filepathRequest filePathRequest) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: filepathRequest.FileID})
	filepathRequest.ResultChannel <- filePathResponse{FilePath: file.FilePath, FileSize: file.FileSize, Error: err}
}
