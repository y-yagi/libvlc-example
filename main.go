package main

import (
	"log"
	"os"

	vlc "github.com/adrg/libvlc-go/v3"
)

func main() {
	if err := vlc.Init("--quiet"); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()

	player, err := vlc.NewListPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	err = player.SetPlaybackMode(vlc.Loop)
	if err != nil {
		log.Fatal(err)
	}

	list, err := vlc.NewMediaList()
	if err != nil {
		log.Fatal(err)
	}
	defer list.Release()

	err = list.AddMediaFromPath(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if err = player.SetMediaList(list); err != nil {
		log.Fatal(err)
	}

	manager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		close(quit)
	}

	eventID, err := manager.Attach(vlc.MediaListPlayerPlayed, eventCallback, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Detach(eventID)

	if err = player.Play(); err != nil {
		log.Fatal(err)
	}

	<-quit
}
