package limiter

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

//rateLimiterChannels represents the channel to communicate with the rate limiter
type rateLimiterChannels struct {
	Send            chan actionRequest
	Request         chan actionRequest
	SendUrgent      chan actionRequest
	RequestUrgent   chan actionRequest
	FilePathRequest chan filePathRequest
}

//actionRequest represents a request for the rate limiter
type actionRequest struct {
	Action        tgbotapi.Chattable
	ResultChannel chan actionResult
}

//actionResult represents the result of a rate limiter action
type actionResult struct {
	RequestResult *tgbotapi.APIResponse
	SendResult    *tgbotapi.Message
	Error         error
}

//filePathRequest represents a rate-limited request for a Telegram File Path
type filePathRequest struct {
	FileID        string
	ResultChannel chan filePathResponse
}

//filePathResponse represents the result of a rate-limited request for a Telegram File Path
type filePathResponse struct {
	FilePath string
	FileSize int
	Error    error
}
