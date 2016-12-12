package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

const usageString = "Usage: red-alert [config-file]"

type Config struct {
	PollingInterval    int    // Poll New Relic every N seconds
	DurationForAverage int    // Calculate average response time based on the last N minutes
	ApiKey             string // The New Relic API key
	ApplicationId      string // The New Relic application ID
}

var defaultConfig Config = Config{
	30, // PollingInterval
	15, // DurationForAverage
	"", // ApiKey
	"", // ApplicationId
}

// TODO: Allow setting the ApiKey and AppId, at least, with command-line flags.
func GetConfig() Config {
	config := defaultConfig
	var logfile string

	switch len(flag.Args()) {
	case 0:
		logfile = os.Getenv("HOME") + "/.red-alert.config"
	case 1:
		logfile = flag.Arg(0)
	default:
		log.Fatal(usageString)
	}

	if _, err := toml.DecodeFile(logfile, &config); err != nil {
		log.Fatal("Couldn't read config file %s: %s", logfile, err)
	}

	if config.ApiKey == "" {
		log.Fatal("No ApiKey found in the config file!")
	}
	if config.ApplicationId == "" {
		log.Fatal("No ApplicationId found in the config file!")
	}
	return config
}
