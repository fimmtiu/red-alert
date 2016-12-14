package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

const usageString = "Usage: red-alert [config-file]"

type Config struct {
	PollingInterval    uint   // Poll New Relic every N seconds
	DurationForAverage uint   // Calculate average response time based on the last N minutes
	VolumeAdjustment   int    // Volume change in decibels
	ApiKey             string // The New Relic API key
	ApplicationId      string // The New Relic application ID
	SoundsDir          string // The directory containing the sound files
}

var defaultConfig Config = Config{
	30,  // PollingInterval
	600, // DurationForAverage
	0,   // VolumeLevel
	"",  // ApiKey
	"",  // ApplicationId
	".", // SoundsDir
}

// TODO: Allow setting the ApiKey and AppId, at least, with command-line flags.
// Also, it would be nice to allow setting the VolumeAdjustment on a per-sound basis.
func GetConfig() Config {
	config := defaultConfig
	var logfile string

	flag.Parse()
	switch len(flag.Args()) {
	case 0:
		logfile = os.Getenv("HOME") + "/.red-alert.config"
	case 1:
		logfile = flag.Arg(0)
	default:
		log.Fatal(usageString)
	}

	if _, err := toml.DecodeFile(logfile, &config); err != nil {
		log.Fatalf("Couldn't read config file %s: %s", logfile, err)
	}

	if config.ApiKey == "" {
		log.Fatal("No ApiKey found in the config file!")
	}
	if config.ApplicationId == "" {
		log.Fatal("No ApplicationId found in the config file!")
	}
	return config
}
