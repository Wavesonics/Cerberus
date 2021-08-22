package routes

import (
	"Cerberus/github"
	"crypto/hmac"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

// RebuildRoute Github webhook handler to restart/rebuild cerberus.
// Github uses an HMAC hex digest to compute the payload hash using the provided secret.
// The hash is in the header X-Hub-Signature-256
func RebuildRoute(webhookSecret string) gin.HandlerFunc {
	webhookSecretBytes := []byte(webhookSecret)
	return func(c *gin.Context) {
		// get provided hash
		providedHash := c.Request.Header.Get("X-Hub-Signature-256")

		if providedHash == "" {
			c.String(http.StatusForbidden, "Missing signature")
			return
		}

		providedHashBytes, hexErr := github.DecodeHex(providedHash)
		if hexErr != nil {
			c.String(http.StatusForbidden, "Invalid signature")
			return
		}

		// calculate hash from body using secret
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			glog.Error(err)
			c.String(http.StatusBadRequest, "Bad request body")
			return
		}

		computedHash := github.ComputeHmac256(jsonData, webhookSecretBytes)

		// compare hashes using a constant time comparer for security
		if !hmac.Equal(providedHashBytes, computedHash) {
			glog.Errorln("Webhook signature does not match")
			c.String(http.StatusForbidden, "Invalid signature")
			return
		}

		var payload github.GithubWebhookPayload
		payloadErr := json.Unmarshal(jsonData, &payload)
		if payloadErr != nil {
			glog.Errorln("Error deserializing webhook payload")
			c.String(http.StatusBadRequest, "Error deserializing webhook payload")
		}

		// we only want to restart when master is updated
		if payload.Ref != "refs/heads/master" {
			return
		}

		// restart cerberus
		cmd, outBuff := exec.Command("/bin/sh", "service-action.sh", "restart", "cerberus"), new(strings.Builder)
		cmd.Stdout = outBuff
		commandError := cmd.Run()
		if commandError != nil {
			c.String(http.StatusInternalServerError, "Unexpected error")
			return
		}

		glog.Infoln("New commit pushed. Restarting cerberus...")
	}
}
