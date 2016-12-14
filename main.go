package main

import (
	"log"
)

// I've written the word "threshold" so many times it has ceased to have
// any meaning. It just looks like a jumble of letters now.
type Threshold struct {
	name         string
	responseTime float32
	soundFile    string
}

// TODO: Load these from the config file.
var thresholds = [...]Threshold{
	Threshold{"Normal", -1.0, ""},                        // No issues
	Threshold{"Yellow Alert", 170.0, "yellow_alert.mp3"}, // Bad but semi-functional
	Threshold{"Red Alert", 250.0, "red_alert.mp3"},       // We're probably boned
}

func main() {
	config := GetConfig()
	poller := NewResponseTimePoller(config)
	player := NewSoundPlayer(config)
	currentThreshold := thresholds[0]

	for responseTime := range poller {
		previousThreshold := currentThreshold
		currentThreshold = getThreshold(responseTime)

		// This happens when things are getting bad.
		if previousThreshold.responseTime < currentThreshold.responseTime && currentThreshold.soundFile != "" {
			log.Printf("Problem: %s! Response time: %.1f ms", currentThreshold.name, responseTime)
			player <- currentThreshold.soundFile
		} else if previousThreshold.responseTime > currentThreshold.responseTime {
			log.Printf("Recovery: %s. Response time: %.1f ms", currentThreshold.name, responseTime)
		}
	}
}

func getThreshold(responseTime float32) (threshold Threshold) {
	for i := 0; i < len(thresholds) && responseTime > thresholds[i].responseTime; i++ {
		threshold = thresholds[i]
	}
	return threshold
}
