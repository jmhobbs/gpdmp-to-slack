package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
)

type Slack struct {
	Client      *slack.Client
	CurrentSong Song
	Set         bool
}

func NewSlack(token string) *Slack {
	return &Slack{slack.New(token), Song{}, true}
}

func (s *Slack) Sync(updates chan Song) {
	for {
		select {
		case song := <-updates:
			if !s.CurrentSong.Equal(song) {
				log.Printf("Sync: %s by %s\n", song.Title, song.Artist)
				s.Client.SetUserCustomStatus(fmt.Sprintf("%s by %s", song.Title, song.Artist), ":musical_note:")
				s.CurrentSong = song
			}
		case <-time.After(time.Second * 10):
			if s.Set {
				log.Println("Clearing Status")
				s.Client.SetUserCustomStatus("", "")
				s.CurrentSong = Song{}
				s.Set = false
			}
		}
	}
}
