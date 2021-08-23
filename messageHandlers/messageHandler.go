package messageHandlers

import (
	"Cerberus/config"
	"Cerberus/serviceControl"
	"Cerberus/telegram"
	"github.com/golang/glog"
)

func HandleMessage(message telegram.Message, botId string, chatId int64, services config.ServiceConfig) {
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

	// Handle the actual message text
	switch message.Text {
	case "/command@CerberusTheGameServerBot":
		startCommand(message, botId, services)
	case "/status@CerberusTheGameServerBot":
		telegram.SendBotMessageSimple("Status is not yet implemented.", botId, chatId)
	case "/stopall@CerberusTheGameServerBot":
		stopAll(botId, chatId, services)
	default:
		telegram.SendBotMessageSimple("I don't know what you mean...", botId, chatId)
	}
}

func startCommand(message telegram.Message, botId string, services config.ServiceConfig) {
	var keyboardButtons []telegram.InlineKeyboardButton
	for _, service := range services.Service {
		keyboardButton := telegram.InlineKeyboardButton{
			Text:         service.Name,
			CallbackData: callbackData1(service.Service),
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

func stopAll(botId string, chatId int64, services config.ServiceConfig) {
	for _, service := range services.Service {
		serviceControl.ExecuteServiceAction(service.Service, serviceControl.Stop, botId, chatId)
	}
}
