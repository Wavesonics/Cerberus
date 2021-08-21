package main

import (
	"Cerberus/telegram"
	"fmt"
	"github.com/golang/glog"
	"strings"
)

func handleCallbackQuery(callbackQuery telegram.CallbackQuery, botId string, chatId int64) {
	glog.Infoln("Handling callback query")

	telegram.AckCallbackQuery(callbackQuery, botId)

	if callbackQuery.Message == nil || callbackQuery.Message.Chat.ID != chatId {
		glog.Warningln("Bot message from incorrect channel\n")
		return
	}

	if callbackQuery.Data == nil {
		glog.Warningln("No data, dropping callback query\n")
		return
	} else {
		glog.Infof("Data: %s\n", *callbackQuery.Data)
	}

	data := *callbackQuery.Data
	parts := strings.Split(data, ":")
	if len(parts) < 2 || (parts[0] != "1" && parts[0] != "2") {
		glog.Error("Bad callback data, bailing.\n")
		return
	}

	if parts[0] == "1" {
		service := parts[1]

		// Pop up the second keyboard
		sendMessageReq := telegram.EditMessageTextBody{
			ChatID:    &callbackQuery.Message.Chat.ID,
			MessageID: &callbackQuery.Message.MessageId,
			Text:      fmt.Sprintf("What action to preform on %s", service),
			ReplyMarkup: &telegram.InlineKeyboardMarkup{
				InlineKeyboard: [][]telegram.InlineKeyboardButton{
					{
						telegram.InlineKeyboardButton{
							Text:         "start",
							CallbackData: callbackData2(service, "start"),
						},
						telegram.InlineKeyboardButton{
							Text:         "stop",
							CallbackData: callbackData2(service, "stop"),
						},
						telegram.InlineKeyboardButton{
							Text:         "restart",
							CallbackData: callbackData2(service, "restart"),
						},
					},
				},
			},
		}

		telegram.SendBotEditMessageText(sendMessageReq, botId)
	} else if parts[0] == "2" {
		service := parts[1]
		action := parts[2]

		// Delete the callback message
		deleteMessageBody := telegram.DeleteMessageBody{
			MessageID: callbackQuery.Message.MessageId,
			ChatID:    callbackQuery.Message.Chat.ID,
		}
		telegram.SendBotDeleteMessage(deleteMessageBody, botId)

		// Actually execute the action finally
		executeServiceAction(service, action, botId, chatId)
	} else {
		glog.Error("Bad callback sequence, bailing.\n")
	}
}