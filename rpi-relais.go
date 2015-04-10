package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stianeikeland/go-rpio"
)

const LISTEN_ADDRESS = ":8080"

type GitHookJSON struct {
	Head string `json:"head"`
	Ref  string `json:"ref"`
	Size int    `json:"size"`
}

func main() {
	fmt.Printf("Hello, world.\n")

	err := rpio.Open()
	if err != nil {
		panic(err)
	}
}

func initServer() {
	router := gin.Default()

	router.POST("/git_hook", handleGitHook)
	router.GET("/demo_event", handleDemoEvent)

	router.Run(LISTEN_ADDRESS)
}

func handleGitHook(c *gin.Context) {
	// check branch for production
	// turn on light for 1 min
	var json GitHookJSON

	c.Bind(&json)

	if json.Ref == "ref/heads/production" {
		c.String(http.StatusOK, "OK")
	}
}

func handleDemoEvent(c *gin.Context) {
	// TODO
	// check demo calendar myself or by external source?
}
