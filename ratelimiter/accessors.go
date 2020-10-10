package limiter

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendUrgent wraps `tgbotapi.Send` and uses the rate limiter to send something urgently
func SendUrgent(chattable tgbotapi.Chattable) (requestResult *tgbotapi.APIResponse, sendResult *tgbotapi.Message, err error) {
	return doSend(chattable, true)
}

//Send wraps `tgbotapi.Send` and uses the rate limiter to send something
func Send(chattable tgbotapi.Chattable) (requestResult *tgbotapi.APIResponse, sendResult *tgbotapi.Message, err error) {
	return doSend(chattable, false)
}

//RequestUrgent wraps `tgbotapi.Request` and uses the rate limiter to request something urgently
func RequestUrgent(chattable tgbotapi.Chattable) (requestResult *tgbotapi.APIResponse, sendResult *tgbotapi.Message, err error) {
	return doRequest(chattable, true)
}

//Request wraps `tgbotapi.Request` and uses the rate limiter to request something
func Request(chattable tgbotapi.Chattable) (requestResult *tgbotapi.APIResponse, sendResult *tgbotapi.Message, err error) {
	return doRequest(chattable, true)
}

//GetTelegramFile gets the FilePath of a specific fileID
func GetTelegramFile(fileID string) (telegramFilePath string, fileSize int, err error) {
	return doGetTelegramFilePath(fileID)
}

/* RATE LIMITED FUNCTIONS */

func doSend(chattable tgbotapi.Chattable, urgent bool) (requestResult *tgbotapi.APIResponse, sendResult *tgbotapi.Message, err error) {

	/* CREATE ACTION REQUEST */
	sendAction := actionRequest{Action: chattable, ResultChannel: make(chan actionResult)}

	/* SEND IT TO APPROPRIATE CHANNEL */
	if urgent {
		channels.SendUrgent <- sendAction
	} else {
		channels.Send <- sendAction
	}

	/* HANDLE RESULT */
	result := <-sendAction.ResultChannel
	return result.RequestResult, result.SendResult, result.Error
}

func doRequest(chattable tgbotapi.Chattable, urgent bool) (requestResult *tgbotapi.APIResponse, sendResult *tgbotapi.Message, err error) {

	/* CREATE ACTION REQUEST */
	requestAction := actionRequest{Action: chattable, ResultChannel: make(chan actionResult)}

	/* SEND IT TO APPROPRIATE CHANNEL */
	if urgent {
		channels.RequestUrgent <- requestAction
	} else {
		channels.Request <- requestAction
	}

	/* HANDLE RESULT */
	result := <-requestAction.ResultChannel
	return result.RequestResult, result.SendResult, result.Error
}

func doGetTelegramFilePath(fileID string) (telegramFilePath string, fileSize int, err error) {

	/* CREATE FILEPATH REQUEST */
	filepathRequest := filePathRequest{FileID: fileID, ResultChannel: make(chan filePathResponse)}

	/* SEND REQUEST TO THE RATE LIMITER */
	channels.FilePathRequest <- filepathRequest
	result := <-filepathRequest.ResultChannel
	return result.FilePath, result.FileSize, result.Error
}
