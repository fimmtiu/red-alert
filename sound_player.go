package main

import (
	"log"
	"os/exec"
)

func NewSoundPlayer() chan string {
	player := make(chan string)

	go func() {
		for soundFile := range player {
			command := exec.Command("play", soundFile)
			if err := command.Run(); err != nil {
				log.Printf("Error playing sound: %s", err)
			}
		}
	}()

	return player
}
