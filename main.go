package main

import (
	"Cerberus/telegram"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	ginglog "github.com/szuecs/gin-glog"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	var ipAddr string
	var portNum int

	const nullArg = ""

	var auth, botId, certFile, keyFile string

	var chatId int64

	flag.StringVar(&auth, "auth", nullArg, "Authentication password")
	flag.StringVar(&botId, "botid", nullArg, "Telegram BotId")
	flag.Int64Var(&chatId, "chatid", -1, "Telegram ChatId")
	flag.StringVar(&certFile, "cert", nullArg, "TLS certificate filename")
	flag.StringVar(&keyFile, "key", nullArg, "TLS key filename")

	flag.StringVar(&ipAddr, "a", "0.0.0.0", "IP address for repository  to listen on")
	flag.IntVar(&portNum, "p", 8080, "TCP port for repository to listen on")
	flag.Parse()

	glog.Infof("auth: %s\n", auth)
	glog.Infof("port: %d\n", portNum)

	if auth == nullArg {
		glog.Error("auth not provided")
		return
	}

	serveAddr := net.JoinHostPort(ipAddr, strconv.Itoa(portNum))
	router := initApp(auth, botId, chatId)

	var err error
	if certFile != nullArg && keyFile != nullArg {
		glog.Infof("Listening on port %d via TLS\n", portNum)
		err = http.ListenAndServeTLS(serveAddr, certFile, keyFile, router)
	} else {
		glog.Infof("Listening on port %d\n", portNum)
		err = http.ListenAndServe(serveAddr, router)
	}
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("Finished.\n")
}

func validateInput(name string, available []string) bool {
	found := false

	for _, service := range available {
		if name == service {
			found = true
			break
		}
	}

	return found
}

func initApp(auth string, botId string, chatId int64) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(ginglog.Logger(3 * time.Second))

	services := []string{"7daystodie", "factorio", "minecraft", "minetally"}
	actions := []string{"start", "stop", "restart"}

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusTeapot, "The Teapot is intact")
	})

	router.GET("/service/:name/:action", func(c *gin.Context) {
		glog.Info("Received service action request\n")

		providedAuth := c.Query("auth")

		name := c.Param("name")
		action := c.Param("action")

		if providedAuth == auth && validateInput(name, services) && validateInput(action, actions) {
			success := executeServiceAction(name, action, botId, chatId)
			if success {
				message := fmt.Sprintf("Service action %s on %s successfull", action, name)
				c.String(http.StatusOK, message)

				telegram.SendBotMessageSimple(message, botId, chatId)
				glog.Info("Action performed\n")
			} else {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Service action %s on %s FAILED", action, name))
				glog.Info("Action failed\n")
			}
		} else {
			glog.Infof("Service '%s' Action '%s'\n", name, action)
			glog.Info("Bad arguments.\n")
			c.String(http.StatusUnauthorized, "Bad Arguments")
		}
	})

	router.POST("/incoming", func(c *gin.Context) {
		glog.Info("Received incoming Telegram Bot request\n")

		providedAuth := c.Query("auth")
		if providedAuth != auth {
			c.String(http.StatusUnauthorized, "Bad Arguments")
			return
		}

		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			glog.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		request := telegram.WebhookUpdateBody{}

		err = json.Unmarshal(jsonData, &request)
		if err != nil {
			glog.Error(err)
			c.String(http.StatusFailedDependency, "Failed to decode request body")
			return
		}

		if request.Message != nil {
			handleMessage(*request.Message, botId, chatId)
		} else if request.InlineQuery != nil {
			handleInlineQuery(*request.InlineQuery, botId, chatId)
		} else if request.CallbackQuery != nil {
			handleCallbackQuery(*request.CallbackQuery, botId, chatId)
		} else {
			glog.Infoln("Unhandled message type. Throwing it away.")
			glog.Infoln(string(jsonData))
		}

		c.String(http.StatusOK, "Message consumed")
	})

	return router
}

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

func handleInlineQuery(inlineQuery telegram.InlineQuery, botId string, chatId int64) {
	glog.Infoln("Handling inline query")
	glog.Infof("Got inline Query: %s\n", inlineQuery.Query)
}

func handleCallbackQuery(callbackQuery telegram.CallbackQuery, botId string, chatId int64) {
	glog.Infoln("Handling callback query")

	ackCallbackQuery(callbackQuery, botId)

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
			Text:      fmt.Sprintf("What action to preform on %s", *callbackQuery.Data),
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

func ackCallbackQuery(callbackQuery telegram.CallbackQuery, botId string) {
	answerCallbackBody := telegram.AnswerCallbackQueryBody{
		CallbackQueryId: callbackQuery.Id,
	}
	telegram.SendBotAnswerCallback(answerCallbackBody, botId)
}

func callbackData1(service string) *string {
	data := fmt.Sprintf("%d:%s", 1, service)
	return &data
}

func callbackData2(service string, action string) *string {
	data := fmt.Sprintf("%d:%s:%s", 2, service, action)
	return &data
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

func executeServiceAction(serviceName string, action string, botId string, chatId int64) bool {
	cmd, outBuff := exec.Command("/bin/sh", "service-action.sh", action, serviceName), new(strings.Builder)
	cmd.Stdout = outBuff
	err := cmd.Run()

	glog.Infof("Running command %s on %s\n", action, serviceName)

	if err != nil {
		glog.Error(err)
		return false
	} else {
		var fmtStr string

		switch action {
		case "start":
			fmtStr = "I have brought %s to life."
		case "stop":
			fmtStr = "I have killed %s."
		case "restart":
			fmtStr = "Like a Phoenix %s is reborn."
		}

		message := fmt.Sprintf(fmtStr, serviceName)
		telegram.SendBotMessageSimple(message, botId, chatId)

		return true
	}
}
