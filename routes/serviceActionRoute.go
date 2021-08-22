package routes

import (
	"Cerberus/serviceControl"
	"Cerberus/telegram"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
)

func ServiceActionRoute(auth string, botId string, chatId int64, gameServices []string, actions []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		glog.Info("Received service action request\n")

		providedAuth := c.Query("auth")

		name := c.Param("name")
		action := c.Param("action")

		if providedAuth == auth && validateInput(name, gameServices) && validateInput(action, actions) {
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