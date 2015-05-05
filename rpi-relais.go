package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stianeikeland/go-rpio"
)

const listenAddress = ":8080"
const pinRed = 4

//const pinX = 17
const pinYellow = 27
const pinHorn = 22

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

	d_param, err := strconv.Atoi(c.Request.Form.Get("d"))
	if err != nil {
		d_param = 0
	}
	duration := time.Duration(d_param) * time.Second

	switch c.Request.Form.Get("c") {
	case "red":
		go activatePin(pinRed, duration)
	case "yellow":
		go activatePin(pinYellow, duration)
	case "horn":
		go activatePin(pinHorn, duration)
	}
}

func activatePin(p int, d time.Duration) {
	pin := rpio.Pin(p)
	pin.Output()
	pin.Low()
	time.Sleep(d)
	pin.High()
}
