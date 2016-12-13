package main

import (
	"log"
	"os/exec"
	"path"
)

func NewSoundPlayer(config Config) chan string {
	player := make(chan string)

	go func() {
		for soundFile := range player {
			command := exec.Command("play", path.Join(config.SoundsDir, soundFile))
			if err := command.Run(); err != nil {
				log.Printf("Error playing sound: %s", err)
			}
		}
	}()

	return player
}
