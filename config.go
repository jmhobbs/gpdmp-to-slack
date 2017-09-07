package main

import "flag"

type Config struct {
	Emoji string
}

var config Config

func init() {
	flag.StringVar(&config.Emoji, "emoji", ":musical_note:", "Emoji used when music is playing.")
}
