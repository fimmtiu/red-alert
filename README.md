# red-alert

This is a little daemon which monitors average site response time in New
Relic, then plays sounds if it goes above predefined thresholds. This is
just a simple first draft where the thresholds are hard-coded in main.go,
but it works and is being used in production.

## Installation

Just `go install` like any other go program!

For Linux users, I've included a sample systemd service script which you
can use to run it at startup.

For sound output, it requires that the `sox` library be installed such that
the `play` command exists in its PATH. (On some systems, you may also need
to install the `libsox-fmt-mp3` library to get MP3 playback.)

## Configuration

The config file looks like:

```
# No, these are not real credentials.
ApiKey = "e49524050d4b8e04f4d0e886b82921d74e58f05169105de"
ApplicationId = "2031337"
SoundsDir = "/home/fimmtiu/go/src/github.com/fimmtiu/red-alert/sounds"
```

It knows the following variables:

**ApiKey** (string)

Your New Relic API key. No default, obviously; must be provided.

**ApplicationId** (string)

Your New Relic application ID. No default, obviously; must be provided.

**DurationForAverage** (unsigned integer)

The "average response time" value will be calculated based on the response
times from the past `DurationForAverage` seconds. Defaults to 600 seconds
(10 minutes).

**PollingInterval** (unsigned integer)

The time between calls, in seconds, to the New Relic API to fetch the response time.
Defaults to 30 seconds.

**SoundsDir** (string)

The directory containing the sound files. Defaults to the directory the
daemon was started in.

**VolumeAdjustment** (integer)

The adjustment, in positive or negative decibels, to apply to the sound
files during playback. Defaults to 0 (no adjustment).
