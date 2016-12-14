package main

import (
	"log"
	"os/exec"
	"path"
	"strconv"
)

func NewSoundPlayer(config Config) chan string {
	player := make(chan string)
	volumeAdjustment := strconv.Itoa(config.VolumeAdjustment)

	go func() {
		for soundFile := range player {
			soundPath := path.Join(config.SoundsDir, soundFile)
			log.Printf("Playing %s", soundPath)
			command := exec.Command("play", soundPath, "vol", volumeAdjustment, "dB")
			if err := command.Run(); err != nil {
				log.Printf("Error playing sound: %s", err)
			}
		}
	}()

	return player
}
