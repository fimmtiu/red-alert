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
			command := exec.Command("play", path.Join(config.SoundsDir, soundFile), "vol", volumeAdjustment, "dB")
			if err := command.Run(); err != nil {
				log.Printf("Error playing sound: %s", err)
			}
		}
	}()

	return player
}
