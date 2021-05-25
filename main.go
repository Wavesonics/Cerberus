package main

import (
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

	router.GET("/service/:name/:action", func(c *gin.Context) {
		glog.Info("Received service action request\n")

		providedAuth := c.Query("auth")

		name := c.Param("name")
		action := c.Param("action")

		if providedAuth == auth && validateInput(name, services) && validateInput(action, actions) {
			success := executeServiceAction(name, action)
			if success {
				message := fmt.Sprintf("Service action %s on %s successfull", action, name)
				c.String(http.StatusOK, message)

				sendBotMessage(message, botId, chatId)
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

	router.GET("/incoming", func(c *gin.Context) {
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

		request := webhookReqBody{}

		err = json.Unmarshal(jsonData, &request)
		if err != nil {
			glog.Error(err)
			c.String(http.StatusBadRequest, "Failed to decode request body")
			return
		}

		if request.Message.Chat.ID != chatId {
			glog.Warningln("Bot message from incorrect channel: %s", strconv.FormatInt(request.Message.Chat.ID, 10))
			c.String(http.StatusNotAcceptable, "Message from invalid Chat")
			return
		}

		sendBotMessage("Hey I got your message", botId, chatId)
	})

	return router
}

func executeServiceAction(serviceName string, action string) bool {
	cmd, outBuff := exec.Command("/bin/sh", "service-action.sh", action, serviceName), new(strings.Builder)
	cmd.Stdout = outBuff
	err := cmd.Run()

	glog.Infof("Running command %s on %s\n", action, serviceName)

	if err != nil {
		glog.Error(err)
		return false
	} else {
		return true
	}
}
