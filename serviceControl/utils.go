package serviceControl

import (
	"Cerberus/telegram"
	"fmt"
	"github.com/golang/glog"
	"os/exec"
	"strings"
)

func ExecuteServiceAction(serviceName string, action string, botId string, chatId int64) bool {
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
