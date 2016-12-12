package main

import (
	"log"
)

func NewSoundPlayer() chan string {
	player := make(chan string)

	go func() {
		for soundFile := range player {
			log.Printf("FIXME: Played %s", soundFile)
		}
	}()

	return player
}
