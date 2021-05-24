package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	ginglog "github.com/szuecs/gin-glog"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	var ipAddr string
	var portNum int

	nullAuth := "null"
	var auth string

	flag.StringVar(&auth, "auth", nullAuth, "Authentication password")

	flag.StringVar(&ipAddr, "a", "0.0.0.0", "IP address for repository  to listen on")
	flag.IntVar(&portNum, "p", 8080, "TCP port for repository to listen on")
	flag.Parse()

	glog.Infof("auth: %s\n", auth)

	if auth == nullAuth {
		glog.Error("auth not provided")
		return
	}

	serveAddr := net.JoinHostPort(ipAddr, strconv.Itoa(portNum))
	router := initApp(auth)
	http.ListenAndServe(serveAddr, router)
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

func initApp(auth string) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(ginglog.Logger(3 * time.Second))

	services := []string{"7daystodie"}
	actions := []string{"start", "stop", "restart"}

	router.GET("/service/:name/*action", func(c *gin.Context) {
		providedAuth := c.Query("auth")

		name := c.Param("name")
		action := c.Param("action")

		if providedAuth == auth && validateInput(name, services) && validateInput(action, actions) {
			executeCommand(fmt.Sprintf("systemctl %s %s", name, action))

			c.String(http.StatusOK, "Service")
		} else {
			c.String(http.StatusUnauthorized, "Incorrect authorization")
		}
	})

	return router
}

func executeCommand(command string) {
	cmd := exec.Command("/bin/sh", command)

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
