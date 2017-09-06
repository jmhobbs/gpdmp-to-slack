package main

import (
	"os"
	"time"
)

// API watches file -> Changes sent to channel -> Slack reads channel -> Posts to API

func main() {
	api := NewSlack(os.Getenv("SLACK_TOKEN"))
	gpdmp := &GPDMPAPI{os.Getenv("GPDMPAPI_PATH")}

	api.Init()

	updates := make(chan Song)
	done := make(chan bool)

	go gpdmp.Watch(updates, done, 5*time.Second)
	go api.Sync(updates, 15*time.Second)
	<-done
}
