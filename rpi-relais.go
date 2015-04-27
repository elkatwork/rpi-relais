package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stianeikeland/go-rpio"
)

const listenAddress = ":8080"
const pinYellow = 4

type GitHookJSON struct {
	Head string `json:"head"`
	Ref  string `json:"ref"`
	Size int    `json:"size"`
}

func main() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()

	initServer()
}

func initServer() {
	router := gin.Default()

	router.POST("/git_hook", handleGitHook)
	router.GET("/demo_event", handleDemoEvent)
	router.GET("/test", handleTest)

	router.Run(listenAddress)
}

func handleGitHook(c *gin.Context) {
	// check branch for production
	// turn on light for 1 min
	var json GitHookJSON

	c.Bind(&json)
	c.String(http.StatusOK, "OK")

	if json.Ref == "refs/heads/production" {
		go activatePin(pinYellow, time.Minute)
	}
}

func handleDemoEvent(c *gin.Context) {
	// TODO
	// check demo calendar myself or by external source?
}

func handleTest(c *gin.Context) {
	c.Request.ParseForm()
	if c.Request.Form.Get("p") == "1234" {
		go activatePin(pinYellow, time.Minute)
	}
}

func activatePin(p int, d time.Duration) {
	pin := rpio.Pin(p)
	pin.Output()
	pin.Low()
	time.Sleep(d)
	pin.High()
}
