package messageHandlers

import (
	"Cerberus/telegram"
	"github.com/golang/glog"
)

func HandleMessage(message telegram.Message, botId string, chatId int64, services []string) {
	glog.Infoln("Handling message")

	if message.Chat.ID != chatId {
		glog.Warningf("Bot message from incorrect channel: %d Must be from: %d\n", message.Chat.ID, chatId)
		return
	}

	/**
		Current command list:

	command - Start a command sequence
	status - Whats up with cerberus
	stopall - Stop all running game servers
	*/

	success := false

	// Handle the actual message text
	switch message.Text {
	case "/command@CerberusTheGameServerBot":
		startCommand(message, botId, services)
		success = true
	case "/status@CerberusTheGameServerBot":
		telegram.SendBotMessageSimple("Status is not yet implemented.", botId, chatId)
		success = true
	case "/stopall@CerberusTheGameServerBot":
		telegram.SendBotMessageSimple("stopall is not yet implemented.", botId, chatId)
		success = true
	default:
		success = false
	}

	if !success {
		telegram.SendBotMessageSimple("I don't know what you mean...", botId, chatId)
	}
}

func startCommand(message telegram.Message, botId string, services []string) {
	var keyboardButtons []telegram.InlineKeyboardButton
	for _, service := range services {
		keyboardButton := telegram.InlineKeyboardButton{
			Text:         service,
			CallbackData: callbackData1(service),
		}
		keyboardButtons = append(keyboardButtons, keyboardButton)
	}

	// Send a new message with the keyboard
	sendMessageReq := telegram.SendMessageBody{
		ChatID: message.Chat.ID,
		Text:   "Which Game Server?",
		ReplyMarkup: &telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{keyboardButtons},
		},
	}

	telegram.SendBotMessage(sendMessageReq, botId)
}
