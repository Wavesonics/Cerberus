package messageHandlers

import (
	"Cerberus/telegram"
	"github.com/golang/glog"
)

func HandleInlineQuery(inlineQuery telegram.InlineQuery, botId string, chatId int64) {
	glog.Infoln("Handling inline query")
	glog.Infof("Got inline Query: %s\n", inlineQuery.Query)
}
