package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net/http"
)

const BaseUrlFormat = "https://api.telegram.org/bot%s/%s"

func makeUrl(botId string, path string) string {
	return fmt.Sprintf(BaseUrlFormat, botId, path)
}

func sendBotMessage(message string, botId string, chatId int64) {
	url := makeUrl(botId, "sendMessage")

	request := sendMessageReqBody{
		ChatID: chatId,
		Text:   message,
	}

	sendJsonRequest(url, request)
}

// https://core.telegram.org/bots/api#update
type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// https://core.telegram.org/bots/api#sendmessage
type sendMessageReqBody struct {
	ChatID              int64  `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode"`
	DisableNotification bool   `json:"disable_notification"`
}

func sendJsonRequest(url string, requestBody interface{}) bool {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		glog.Error(err)
		return false
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		glog.Error(err)
		return false
	}

	return true
}
