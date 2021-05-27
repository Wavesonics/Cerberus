package main

import (
	"Cerberus/telegram"
	"github.com/golang/glog"
)

func handleMessage(message telegram.Message, botId string, chatId int64) {
	glog.Infoln("Handling message")

	if message.Chat.ID != chatId {
		glog.Warningf("Bot message from incorrect channel: %d Must be from: %d\n", message.Chat.ID, chatId)
		return
	}

	/**
		Current command list:

	startfactorio - Start the Factorio server
	startminecraft - Start the Minecraft server
	start7d2d - Start the 7 Days to Die server
	stopfactorio - Stop the Factorio server
	stopminecraft - Stop the Minecraft server
	stop7d2d - Stop the 7 Days to Die server
	*/

	success := false

	// Handle the actual message text
	switch message.Text {
	case "/startfactorio@CerberusTheGameServerBot":
		success = executeServiceAction("factorio", "start", botId, chatId)
	case "/start7d2d@CerberusTheGameServerBot":
		success = executeServiceAction("7daystodie", "start", botId, chatId)
	case "/startminecraft@CerberusTheGameServerBot":
		success = executeServiceAction("minecraft", "start", botId, chatId)
	case "/stopfactorio@CerberusTheGameServerBot":
		success = executeServiceAction("factorio", "stop", botId, chatId)
	case "/stop7d2d@CerberusTheGameServerBot":
		success = executeServiceAction("7daystodie", "stop", botId, chatId)
	case "/stopminecraft@CerberusTheGameServerBot":
		success = executeServiceAction("minecraft", "stop", botId, chatId)
	case "/command@CerberusTheGameServerBot":
		startCommand(message, botId)
		success = true
	default:
		success = false
	}

	if !success {
		telegram.SendBotMessageSimple("I don't know what you mean...", botId, chatId)
	}
}

func startCommand(message telegram.Message, botId string) {
	// Send a new message with the keyboard
	sendMessageReq := telegram.SendMessageBody{
		ChatID: message.Chat.ID,
		Text:   "Which Game Server?",
		ReplyMarkup: &telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					telegram.InlineKeyboardButton{
						Text:         "factorio",
						CallbackData: callbackData1("factorio"),
					},
					telegram.InlineKeyboardButton{
						Text:         "7daystodie",
						CallbackData: callbackData1("7daystodie"),
					},
					telegram.InlineKeyboardButton{
						Text:         "minecraft",
						CallbackData: callbackData1("minecraft"),
					},
				},
			},
		},
	}

	telegram.SendBotMessage(sendMessageReq, botId)
}
