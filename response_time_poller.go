package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const urlTemplate = "https://api.newrelic.com/v2/applications/<app_id>/metrics/data.json"

type ResponseTimePoller struct {
	channel chan float32
	client  http.Client
	config  Config
	url     string
}

// Good god, what a gnarled, barnacled, convoluted response format!
type NewRelicResponse struct {
	Metric_data struct {
		Metrics []struct {
			Timeslices []struct {
				Values struct {
					Average_call_time float32
				}
			}
		}
	}
}

func NewResponseTimePoller(config Config) chan float32 {
	url := strings.Replace(urlTemplate, "<app_id>", config.ApplicationId, 1)
	rtp := &ResponseTimePoller{make(chan float32), http.Client{}, config, url}

	go func() {
		rtp.pollNewRelic()

		ticker := time.Tick(time.Duration(config.PollingInterval) * time.Second)
		for _ = range ticker { // Go 1.3-compatible syntax (sigh)
			rtp.pollNewRelic()
		}
	}()

	return rtp.channel
}

// For the time being, average_call_time is the only number we're measuring.
// If we decide to add request queueing time to our NR stats, we'll have to
// do some calculation here.
func (rtp *ResponseTimePoller) pollNewRelic() {
	averageCall, success := rtp.getAverageCallTime()
	if success {
		rtp.channel <- averageCall
	}
}

func (rtp *ResponseTimePoller) getAverageCallTime() (float32, bool) {
	response, err := rtp.makeNewRelicRequest("names[]=HttpDispatcher&values[]=average_call_time")
	if err != nil {
		log.Printf("Couldn't get data from New Relic: %s", err)
		return 0.0, false
	}

	var decodedResponse NewRelicResponse
	err = json.Unmarshal(response, &decodedResponse)
	if err != nil {
		log.Printf("Unexpected JSON from New Relic: %s\n%s", err, decodedResponse)
		return 0.0, false
	}

	return decodedResponse.Metric_data.Metrics[0].Timeslices[0].Values.Average_call_time, true
}

func (rtp *ResponseTimePoller) makeNewRelicRequest(requestBody string) ([]byte, error) {
	duration := time.Duration(rtp.config.DurationForAverage) * time.Minute
	now := time.Now()
	earlier := now.Add(-duration)
	requestBody += "&from=" + earlier.Format(time.RFC3339) + "&to=" + now.Format(time.RFC3339) + "&summarize=true"

	request, err := http.NewRequest("GET", rtp.url, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Api-Key", rtp.config.ApiKey)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	log.Printf("request: %s", request)
	response, err := rtp.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("contents: %s", contents)
	return contents, nil
}
