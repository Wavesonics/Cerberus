package messageHandlers

import (
	"Cerberus/config"
	"Cerberus/serviceControl"
	"Cerberus/telegram"
	"fmt"
	"github.com/golang/glog"
	"os/exec"
	"strconv"
	"strings"
)

func HandleMessage(message telegram.Message, botId string, chatId int64, config config.ServiceConfig) {
	glog.Infof("Handling message: %s", message.Text)

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
		startCommand(message, botId, config)
	case "/status@CerberusTheGameServerBot":
		status(botId, chatId, config)
	case "/stopall@CerberusTheGameServerBot":
		stopAll(botId, chatId, config)
	default:
		telegram.SendBotMessageSimple("I don't know what you mean...", botId, chatId)
	}
}

func startCommand(message telegram.Message, botId string, config config.ServiceConfig) {
	var keyboardButtons []telegram.InlineKeyboardButton
	for _, service := range config.Services {
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

func stopAll(botId string, chatId int64, config config.ServiceConfig) {
	for _, service := range config.Services {
		serviceControl.ExecuteServiceAction(service.Service, serviceControl.Stop, botId, chatId)
	}
	telegram.SendBotMessageSimple("They have all been killed...", botId, chatId)
}

func status(botId string, chatId int64, config config.ServiceConfig) {
	resultMap := make(map[string] string)
	maxLength := 0

	for _, service := range config.Services {
		out, err := exec.Command("bash", "-c", fmt.Sprintf(`systemctl status %s | grep -Po "(?<=Active: )\w+"`, service.Service)).Output()
		logError("bash cmd: ", err)

		serviceStatus := string(out)

		resultMap[service.Service] = strings.TrimSpace(serviceStatus)
		if maxLength < len(service.Name) {
			maxLength = len(service.Name)
		}
	}

	memOut, err := exec.Command("bash", "-c", `free | grep -Po "(?<=Mem:)(\s+\d+){2}"`).Output()
	logError("Get memory", err)
	memOutString := strings.TrimSpace(string(memOut))
	memUsage := calcMemUsageFromFreeCommand(memOutString)

	result := generateTable(maxLength, config, resultMap, memUsage)
	telegram.SendBotMessage(telegram.SendMessageBody{
		ChatID:    chatId,
		Text:      result,
		ParseMode: "MarkdownV2",
	}, botId)
}

func logError(message string, err error) {
	if err != nil {
		glog.Errorln(message, err)
	}
}

func generateTable(maxLength int,
	config config.ServiceConfig,
	resultMap map[string]string,
	memoryUsage int,
) string {
	fmtString := fmt.Sprintf("%%-%ds : %%s\n", maxLength)
	builder := strings.Builder{}
	builder.WriteString("```\n")

	builder.WriteString(fmt.Sprintf(fmtString, "Memory", fmt.Sprintf("%02d%%", memoryUsage)))
	builder.WriteString("\n")

	for _, service := range config.Services {
		builder.WriteString(fmt.Sprintf(fmtString, service.Name, resultMap[service.Service]))
	}
	builder.WriteString("```")
	return builder.String()
}

func calcMemUsageFromFreeCommand(memOutString string) int {
	parts := strings.Fields(memOutString)
	totalMem, _ := strconv.ParseFloat(parts[0], 32)
	usedMem, _ := strconv.ParseFloat(parts[1], 32)

	memUsage := int(usedMem / (totalMem + 1.0) * 100)
	return memUsage
}