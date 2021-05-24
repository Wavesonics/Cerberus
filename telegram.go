package main

import (
	"fmt"
	"github.com/golang/glog"
	"net/http"
)

// https://api.telegram.org/bot1873709271:AAFV9RrNIEC6PI7adTv9-3ydxfcbCM8yriY/sendMessage?chat_id=-1001249096766&text=test

func sendBotMessage(message string, botId string, chatId string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botId, chatId, message)

	_, err := http.Get(url)
	if err != nil {
		glog.Error(err)
	}
}
