package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
)

type Slack struct {
	Client       *slack.Client
	CurrentSong  Song
	Set          bool
	InitialText  string
	InitialEmoji string
}

func NewSlack(token string) *Slack {
	return &Slack{slack.New(token), Song{}, true, "", ""}
}

func (s *Slack) Init() {
	auth, err := s.Client.AuthTest()
	if err != nil {
		log.Fatal(err)
	}

	user, err := s.Client.GetUserInfo(auth.UserID)
	if err != nil {
		log.Fatal(err)
	}

	s.InitialText = user.Profile.StatusText
	s.InitialEmoji = user.Profile.StatusEmoji
	log.Printf("Initial status: %s %s", s.InitialEmoji, s.InitialText)
}

func (s *Slack) Sync(updates chan Song, revert_after time.Duration) {
	for {
		select {
		case song := <-updates:
			if !s.CurrentSong.Equal(song) {
				log.Printf("Sync: %s by %s\n", song.Title, song.Artist)
				s.Client.SetUserCustomStatus(fmt.Sprintf("%s by %s", song.Title, song.Artist), ":musical_note:")
				s.CurrentSong = song
			}
		case <-time.After(revert_after):
			if s.Set {
				log.Printf("Reverting Status: %s %s\n", s.InitialEmoji, s.InitialText)
				s.Client.SetUserCustomStatus(s.InitialText, s.InitialEmoji)
				s.CurrentSong = Song{}
				s.Set = false
			}
		}
	}
}
