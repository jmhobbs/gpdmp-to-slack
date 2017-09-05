package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Song struct {
	Title    string
	Artist   string
	Album    string
	AlbumArt string
}

func (a Song) Equal(b Song) bool {
	return a.Title == b.Title && a.Artist == b.Artist && a.Album == b.Album
}

type PlaybackJSON struct {
	Playing bool
	Song    Song
	Rating  struct {
		Liked    bool
		Disliked bool
	}
	Time struct {
		Current int
		Total   int
	}
	SongLyrics string
	Shuffle    string
	Repeat     string
	Volume     int
}

var lastRead time.Time
var currentSong Song

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if time.Now().After(lastRead.Add(time.Second * 5)) {
						lastRead = time.Now()
						go readFile(event.Name)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/Users/johnhobbs/Library/Application Support/Google Play Music Desktop Player/json_store/playback.json")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func readFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}

	dec := json.NewDecoder(f)
	pb := PlaybackJSON{}

	err = dec.Decode(&pb)
	if err != nil {
		log.Println(err)
		return
	}

	if !currentSong.Equal(pb.Song) {
		log.Printf("%s by %s\n", pb.Song.Title, pb.Song.Artist)
		currentSong = pb.Song
	}
}
