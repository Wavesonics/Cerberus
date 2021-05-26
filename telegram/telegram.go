package telegram

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

func SendBotMessageSimple(message string, botId string, chatId int64) {
	url := makeUrl(botId, "sendMessage")

	request := SendMessageBody{
		ChatID: chatId,
		Text:   message,
	}

	sendJsonRequest(url, request)
}

func SendBotMessage(message SendMessageBody, botId string) {
	url := makeUrl(botId, "sendMessage")
	sendJsonRequest(url, message)
}

func SendBotAnswerCallback(body AnswerCallbackQueryBody, botId string) {
	url := makeUrl(botId, "answerCallbackQuery")
	sendJsonRequest(url, body)
}

func SendBotEditMessageText(body EditMessageTextBody, botId string) {
	url := makeUrl(botId, "editMessageText")
	sendJsonRequest(url, body)
}

func SendBotDeleteMessage(body DeleteMessageBody, botId string) {
	url := makeUrl(botId, "deleteMessage")
	sendJsonRequest(url, body)
}

func sendJsonRequest(url string, requestBody interface{}) bool {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		glog.Error(err)
		return false
	}

	//glog.Infof("Request URL: %s\n", url)
	//glog.Infof("Request Body: %s\n", string(jsonData))

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		glog.Error(err)
		return false
	}

	glog.Infof("Response Code: %d\n", response.StatusCode)

	/*
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			glog.Error("Failed to read response body\n")
		} else {
			bodyStr := string(bodyBytes)
			glog.Infof("Response Body: %s\n", bodyStr)
		}
	*/

	return true
}
