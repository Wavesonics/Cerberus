package main

import (
	"Cerberus/config"
	"Cerberus/routes"
	"Cerberus/telegram"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	ginglog "github.com/szuecs/gin-glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var ipAddr string
	var portNum int

	const nullArg = ""

	var auth, botId, certFile, keyFile string
	var webhookSecret string
	var chatId int64

	flag.StringVar(&auth, "auth", nullArg, "Authentication password")
	flag.StringVar(&botId, "botid", nullArg, "Telegram BotId")
	flag.Int64Var(&chatId, "chatid", -1, "Telegram ChatId")
	flag.StringVar(&certFile, "cert", nullArg, "TLS certificate filename")
	flag.StringVar(&keyFile, "key", nullArg, "TLS key filename")

	flag.StringVar(&webhookSecret, "secret", nullArg, "Webhook secret")

	flag.StringVar(&ipAddr, "a", "0.0.0.0", "IP address for repository  to listen on")
	flag.IntVar(&portNum, "p", 8080, "TCP port for repository to listen on")
	flag.Parse()

	glog.Infoln("Starting cerberus bot...")
	glog.Infof("auth: %s\n", auth)
	glog.Infof("port: %d\n", portNum)

	if auth == nullArg {
		glog.Error("auth not provided")
		return
	}

	serveAddr := net.JoinHostPort(ipAddr, strconv.Itoa(portNum))
	router := initApp(auth, botId, chatId, webhookSecret)

	var err error
	postWakeupMessage(botId, chatId)
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

func loadServicesConfig() config.ServiceConfig {
	data, err := ioutil.ReadFile("services.yaml")
	if err != nil {
		log.Fatal(err)
	}

	gameServices := config.ServiceConfig{}

	err = yaml.Unmarshal(data, &gameServices)
	if err != nil {
		log.Fatalf("error loading services.yaml: %v", err)
	}

	return gameServices
}

func initApp(auth string, botId string, chatId int64, webhookSecret string) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(ginglog.Logger(3 * time.Second))

	gameServices := loadServicesConfig()

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusTeapot, "The Teapot is intact")
	})

	router.GET("/service/:name/:action", routes.ServiceActionRoute(auth, botId, chatId, gameServices))

	router.POST("/incoming", routes.IncomingRoute(auth, botId, chatId, gameServices))

	router.POST("/rebuild", routes.RebuildRoute(webhookSecret))

	return router
}

var wakeupMessages = []string{
	"I have awakened!",
	"I have been summoned yet again.",
	"I am here to guard the gates.",
}

func postWakeupMessage(botId string, chatId int64) {
	index := rand.Intn(len(wakeupMessages))
	telegram.SendBotMessageSimple(wakeupMessages[index], botId, chatId)
}
