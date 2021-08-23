package routes

import (
	"Cerberus/config"
	"Cerberus/serviceControl"
	"Cerberus/telegram"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
)

func ServiceActionRoute(auth string, botId string, chatId int64, gameServices config.ServiceConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		glog.Info("Received service action request\n")

		providedAuth := c.Query("auth")

		name := c.Param("name")
		actionName := c.Param("action")
		action, err := serviceControl.ActionFromString(actionName)
		if err != nil {
			glog.Error("Bad Service Action, bailing.\n")
			return
		}

		if providedAuth == auth && validateService(name, gameServices) {
			success := serviceControl.ExecuteServiceAction(name, action, botId, chatId)
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
	}
}

func validateService(name string, services config.ServiceConfig) bool {
	found := false

	for _, service := range services.Service {
		if name == service.Name {
			found = true
			break
		}
	}

	return found
}
