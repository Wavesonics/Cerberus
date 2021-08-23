package serviceControl

import (
	"Cerberus/telegram"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"os/exec"
	"strings"
)

type ServiceAction int

const (
	Start   ServiceAction = iota // 0
	Stop                         // 1
	Restart                      // 2
	Invalid
)

var validActions = []string{"start", "stop", "restart"}

func (action ServiceAction) String() string {
	return validActions[action]
}

func ActionFromString(actionName string) (ServiceAction, error) {
	var action = Invalid

	for ii, validAction := range validActions {
		if actionName == validAction {
			action = ServiceAction(ii)
			break
		}
	}

	if action == Invalid {
		return Invalid, errors.New("invalid service name")
	} else {
		return action, nil
	}
}

func ExecuteServiceAction(serviceName string, action ServiceAction, botId string, chatId int64) bool {
	cmd, outBuff := exec.Command("/bin/sh", "service-action.sh", action.String(), serviceName), new(strings.Builder)
	cmd.Stdout = outBuff
	err := cmd.Run()

	glog.Infof("Running command %s on %s\n", action, serviceName)

	if err != nil {
		glog.Error(err)
		return false
	} else {
		var fmtStr string

		switch action {
		case Start:
			fmtStr = "I have brought %s to life."
		case Stop:
			fmtStr = "I have killed %s."
		case Restart:
			fmtStr = "Like a Phoenix %s is reborn."
		}

		message := fmt.Sprintf(fmtStr, serviceName)
		telegram.SendBotMessageSimple(message, botId, chatId)

		return true
	}
}
