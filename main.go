package main

import (
	"os"
	"time"
)

// API watches file -> Changes sent to channel -> Slack reads channel -> Posts to API

func main() {
	api := NewSlack(os.Getenv("SLACK_TOKEN"))
	gpdmp := &GPDMPAPI{os.Getenv("GPDMPAPI_PATH")}

	updates := make(chan Song)
	done := make(chan bool)

	go gpdmp.Watch(updates, done, time.Duration(5))
	go api.Sync(updates)
	<-done
}
