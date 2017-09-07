package main

import (
	"flag"
	"os"
	"time"
)

func main() {
	flag.Parse()

	api := NewSlack(os.Getenv("SLACK_TOKEN"))
	gpdmp := &GPDMPAPI{os.Getenv("GPDMPAPI_PATH")}

	api.Init()

	updates := make(chan Song)
	done := make(chan bool)

	go gpdmp.Watch(updates, done, 5*time.Second)
	go api.Sync(config.Emoji, updates, 15*time.Second)
	<-done
}
