package routes

import (
	"Cerberus/config"
	"Cerberus/messageHandlers"
	"Cerberus/telegram"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

func IncomingRoute(auth string, botId string, chatId int64, gameServiceConfig config.ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infoln("Received incoming Telegram Bot request")

		providedAuth := c.GetHeader("X-Telegram-Bot-Api-Secret-Token")
		if providedAuth != auth {
			glog.Infof("Provided auth (%s) did my match (%s)\n", providedAuth, auth)
			c.String(http.StatusUnauthorized, "Not Authorized")
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
			messageHandlers.HandleMessage(*request.Message, botId, chatId, gameServiceConfig)
		} else if request.InlineQuery != nil {
			messageHandlers.HandleInlineQuery(*request.InlineQuery, botId, chatId)
		} else if request.CallbackQuery != nil {
			messageHandlers.HandleCallbackQuery(*request.CallbackQuery, botId, chatId)
		} else {
			glog.Infoln("Unhandled message type. Throwing it away.")
			glog.Infoln(string(jsonData))
		}

		c.String(http.StatusOK, "Message consumed")
	}
}
