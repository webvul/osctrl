package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/jmpsec/osctrl/pkg/types"
	"github.com/jmpsec/osctrl/pkg/utils"
)

const (
	// GELF spec version
	graylogVersion = "1.1"
	// Host as source for GELF data
	graylogHost = "osctrl"
	// Log Level (informational)
	graylogLevel = 6
	// Method to send
	graylogMethod = "POST"
)

// GraylogMessage to handle log format to be sent to Graylog
type GraylogMessage struct {
	Version      string `json:"version"`
	Host         string `json:"host"`
	ShortMessage string `json:"short_message"`
	Timestamp    int64  `json:"timestamp"`
	Level        uint   `json:"level"`
	Environment  string `json:"_environment"`
	Type         string `json:"_type"`
	UUID         string `json:"_uuid"`
}

// GraylogSend - Function that sends JSON logs to Graylog
func GraylogSend(logType string, data []byte, environment, uuid, url string, debug bool) {
	// Prepare headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	// Convert the array in an array of multiple message
	var logs []interface{}
	if logType == types.QueryLog {
		// For on-demand queries, just a JSON blob with results and statuses
		var result interface{}
		err := json.Unmarshal(data, &result)
		if err != nil {
			log.Printf("error parsing data %s %v", string(data), err)
		}
		logs = append(logs, result)
	} else {
		err := json.Unmarshal(data, &logs)
		if err != nil {
			log.Printf("error parsing logs %s %v", string(data), err)
		}
	}
	// Prepare data to send
	for _, l := range logs {
		logMessage, err := json.Marshal(l)
		if err != nil {
			log.Printf("error parsing log %s", err)
			continue
		}
		messsageData := GraylogMessage{
			Version:      graylogVersion,
			Host:         graylogHost,
			ShortMessage: string(logMessage),
			Timestamp:    time.Now().Unix(),
			Level:        graylogLevel,
			Environment:  environment,
			Type:         logType,
			UUID:         uuid,
		}
		// Serialize data using GELF
		jsonMessage, err := json.Marshal(messsageData)
		if err != nil {
			log.Printf("error marshaling data %s", err)
		}
		jsonParam := strings.NewReader(string(jsonMessage))
		if debug {
			log.Printf("Sending %d bytes to Graylog for %s - %s", len(data), environment, uuid)
		}
		// Send log with a POST to the Graylog URL
		resp, body, err := utils.SendRequest(graylogMethod, url, jsonParam, headers)
		if err != nil {
			log.Printf("error sending request %s", err)
			return
		}
		if debug {
			log.Printf("Graylog: HTTP %d %s", resp, body)
		}
	}
}
