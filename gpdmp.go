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

type GPDMPAPI struct {
	Path string
}

func (gp *GPDMPAPI) Watch(updates chan Song, done chan bool, debounce_seconds time.Duration) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		var lastRead time.Time

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if time.Now().After(lastRead.Add(time.Second * debounce_seconds)) {
						lastRead = time.Now()
						f, err := os.Open(event.Name)
						if err != nil {
							log.Println(err)
							continue
						}

						dec := json.NewDecoder(f)
						pb := PlaybackJSON{}

						err = dec.Decode(&pb)
						if err != nil {
							log.Println(err)
							continue
						}

						updates <- pb.Song
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(gp.Path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
