package utility

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	limiter "gitlab.com/shitposting/tg-random-bot/ratelimiter"

	memesapi "gitlab.com/shitposting/memesapi/rest/client"
	"gitlab.com/shitposting/telegram-markdown-processor/telegramCaption"
)

const (
	//MaxRetries is the maximum amount of attempts to
	//forward that the bot should make
	MaxRetries = 3
)

//TrySending tries to send a meme up to 3 times
func TrySending(recent bool, chatID int64, userid int, client *memesapi.Client, bot *tgbotapi.BotAPI) {

	var err error
	var startDate string

	if recent {
		startDate = strconv.FormatInt(time.Now().AddDate(0, 0, -14).Unix(), 10)
	}

	for i := 0; i < MaxRetries; i++ {

		resp, err := client.Random("", startDate, "", strconv.Itoa(userid))
		if err != nil {
			log.Println(err)
		}

		switch resp.Meme.Type {
		case "image":
			bot.Send(tgbotapi.NewChatAction(chatID, "upload_photo"))
			photoConfig := createPhotoConfig(resp.Meme.URL, chatID)
			photoConfig.Caption = fmt.Sprintf("%s\n\n🔎 source: t.me/shitpost/%d", telegramCaption.PrepareCaptionToSend(resp.Meme.Caption, ""), resp.Meme.MessageID)
			photoConfig.ParseMode = "HTML"
			_, _, err = limiter.Send(photoConfig)
		case "video":
			bot.Send(tgbotapi.NewChatAction(chatID, "upload_video"))

			videoConfig := createVideoConfig(resp.Meme.URL, chatID)
			videoConfig.Caption = fmt.Sprintf("%s\n\n🔎 source: t.me/shitpost/%d", telegramCaption.PrepareCaptionToSend(resp.Meme.Caption, ""), resp.Meme.MessageID)
			videoConfig.ParseMode = "HTML"
			_, _, err = limiter.Send(videoConfig)
		default:
			bot.Send(tgbotapi.NewChatAction(chatID, "upload_video"))

			animationConfig := createAnimationConfig(resp.Meme.URL, chatID)
			animationConfig.Caption = fmt.Sprintf("%s\n\n🔎 source: t.me/shitpost/%d", telegramCaption.PrepareCaptionToSend(resp.Meme.Caption, ""), resp.Meme.MessageID)
			animationConfig.ParseMode = "HTML"
			_, _, err = limiter.Send(animationConfig)
		}

		if err == nil {
			return
		}
	}

	log.Println("Unable to send after 3 attempts. Giving up")
	_, _, err = limiter.Send(tgbotapi.NewMessage(chatID, "👷‍♂️ An error has occurred 👷‍♂️\n\nPlease try again in a few minutes!"))
	if err != nil {
		log.Println(fmt.Sprintf("Unable to send error feedback message: %s", err.Error()))
	}

	return
}

func createPhotoConfig(url string, chatID int64) (photoConfig tgbotapi.PhotoConfig) {
	if strings.HasPrefix(url, "/") {
		return tgbotapi.NewPhotoUpload(chatID, url)
	}

	return tgbotapi.NewPhotoShare(chatID, url)
}

func createVideoConfig(url string, chatID int64) (videoConfig tgbotapi.VideoConfig) {
	if strings.HasPrefix(url, "/") {
		return tgbotapi.NewVideoUpload(chatID, url)
	}

	return tgbotapi.NewVideoShare(chatID, url)
}

func createAnimationConfig(url string, chatID int64) (animationConfig tgbotapi.AnimationConfig) {
	if strings.HasPrefix(url, "/") {
		return tgbotapi.NewAnimationUpload(chatID, url)
	}

	return tgbotapi.NewAnimationShare(chatID, url)
}
