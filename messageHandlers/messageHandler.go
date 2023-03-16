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
		telegram.SendBotMessageSimple("🤷 I don't know what you mean...", botId, chatId)
	}
}

func getEmojiForStatus(status string) string {
	var emoji string
	if status == "active" {
		emoji = "▶️"
	} else if status == "inactive" {
		emoji = "🛑"
	} else if status == "failed" {
		emoji = "😵"
	} else {
		emoji = "🤔"
	}

	return emoji
}

func createServiceString(service config.Service, serviceStatuses map[string]string) string {
	status := serviceStatuses[service.Name]
	emoji := getEmojiForStatus(status)

	return service.Name + " " + emoji
}

func startCommand(message telegram.Message, botId string, config config.ServiceConfig) {
	serviceStatuses := getStatuses(config)

	numColumns := 3
	var keyboardButtons [][]telegram.InlineKeyboardButton
	var buttonRow []telegram.InlineKeyboardButton = nil
	for ii, service := range config.Services {
		if ii%numColumns == 0 {
			if buttonRow != nil {
				keyboardButtons = append(keyboardButtons, buttonRow)
				buttonRow = nil
			}
			buttonRow = []telegram.InlineKeyboardButton{}
		}

		serviceButtonText := createServiceString(service, serviceStatuses)

		keyboardButton := telegram.InlineKeyboardButton{
			Text:         serviceButtonText,
			CallbackData: callbackData1(service.Service),
		}
		buttonRow = append(buttonRow, keyboardButton)
	}

	if buttonRow != nil {
		keyboardButtons = append(keyboardButtons, buttonRow)
		buttonRow = nil
	}

	// Send a new message with the keyboard
	sendMessageReq := telegram.SendMessageBody{
		ChatID: message.Chat.ID,
		Text:   "🎮 Which Game Server❓",
		ReplyMarkup: &telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboardButtons,
		},
	}

	telegram.SendBotMessage(sendMessageReq, botId)
}

func stopAll(botId string, chatId int64, config config.ServiceConfig) {
	for _, service := range config.Services {
		serviceControl.ExecuteServiceAction(service.Service, serviceControl.Stop, botId, chatId)
	}
	telegram.SendBotMessageSimple("💀 They have all been killed...", botId, chatId)
}

func getStatuses(config config.ServiceConfig) map[string]string {
	resultMap := make(map[string]string)

	for _, service := range config.Services {
		out, err := exec.Command("bash", "-c", fmt.Sprintf(`systemctl status %s | grep -Po "(?<=Active: )\w+"`, service.Service)).Output()
		logError("bash cmd: ", err)

		serviceStatus := string(out)

		resultMap[service.Name] = strings.TrimSpace(serviceStatus)
	}

	return resultMap
}

func status(botId string, chatId int64, config config.ServiceConfig) {
	serviceStatuses := getStatuses(config)
	maxLength := 0

	for service := range serviceStatuses {
		if maxLength < len(service) {
			maxLength = len(service)
		}
	}

	memOut, err := exec.Command("bash", "-c", `free | grep -Po "(?<=Mem:)(\s+\d+){2}"`).Output()
	logError("Get memory", err)
	memOutString := strings.TrimSpace(string(memOut))
	memUsage := calcMemUsageFromFreeCommand(memOutString)

	cpuOut, cpuErr := exec.Command("bash", "-c", "mpstat | grep all").Output()
	logError("Get CPU", cpuErr)
	cpuOutString := strings.TrimSpace(string(cpuOut))
	cpuUsage := calcCpuUsageFromMpstatCommand(cpuOutString)

	result := generateTable(maxLength, config, serviceStatuses, memUsage, cpuUsage)
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
	memoryUsage string,
	cpuUsage string,
) string {
	fmtString := fmt.Sprintf("%%-%ds : %%s\n", maxLength)
	builder := strings.Builder{}
	builder.WriteString("```\n")

	builder.WriteString(fmt.Sprintf(fmtString, "🧠 Memory", memoryUsage))

	builder.WriteString(fmt.Sprintf(fmtString, "💻 CPU", cpuUsage))
	builder.WriteString("\n")

	for _, service := range config.Services {
		var result = resultMap[service.Name]
		var emoji = getEmojiForStatus(result)
		result = emoji + "️ " + result

		builder.WriteString(fmt.Sprintf(fmtString, service.Name, result))
	}
	builder.WriteString("```")
	return builder.String()
}

func calcMemUsageFromFreeCommand(memOutString string) string {
	parts := strings.Fields(memOutString)
	totalMem, _ := strconv.ParseFloat(parts[0], 32)
	usedMem, _ := strconv.ParseFloat(parts[1], 32)

	memoryUsage := usedMem / (totalMem + 1.0) * 100
	return fmt.Sprintf("%03.1f%%", memoryUsage)
}

func calcCpuUsageFromMpstatCommand(cpuOutString string) string {
	parts := strings.Fields(cpuOutString)
	idle, _ := strconv.ParseFloat(parts[len(parts)-1], 64)
	usage := 100 - idle
	return fmt.Sprintf("%03.1f%%", usage)
}
