package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/k0swe/wsjtx-go/v4"
	"gopkg.in/ini.v1"
)

func wsjtxmain() {
	if _, err := os.Stat(os.Getenv("WAVELOG_AGENT_COFNIG")); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Stat("config.ini"); errors.Is(err, os.ErrNotExist) {
			log.Fatal("Configuration file not found.")
		} else {
			configFile = "config.ini"
		}
	} else {
		configFile = os.Getenv("WAVELOG_AGENT_COFNIG")
	}

	// Configuration
	inidata, err := ini.Load(configFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	err = inidata.MapTo(&Config)
	if err != nil {
		fmt.Printf("Fail to map file: %v", err)
		os.Exit(1)
	}

	log.Println("Listening for WSJT-X...")
	wsjtxServer, err := wsjtx.MakeServer()

	if err != nil {
		log.Fatalf("%v", err)
	}

	wsjtxChannel := make(chan interface{}, 5)
	errChannel := make(chan error, 5)
	go wsjtxServer.ListenToWsjtx(wsjtxChannel, errChannel)

	for {
		select {
		case err := <-errChannel:
			log.Printf("error: %v", err)
		case message := <-wsjtxChannel:
			handleServerMessage(message)
		}
	}
}

func removeLBR(text string) string {
	re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	return re.ReplaceAllString(text, ` `)
}

// When we receive WSJT-X messages, display them.
func handleServerMessage(message interface{}) {
	switch message.(type) {
	case wsjtx.LoggedAdifMessage:
		log.Println("Logged ADIF:", message)
		current_time := time.Now()
		var jsonData = []byte(`{"key":"` + Config.Wavelog.Key + `", "timestamp":"` + current_time.Format("2006/01/02 15:04:05") + `", "station_profile_id":"` + Config.Wavelog.Profile + `", "type":"adif", "string":"` + removeLBR(message.(wsjtx.LoggedAdifMessage).Adif) + `"}`)

		req, err := http.NewRequest("POST", Config.Wavelog.URL+"/api/qso", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()

		if resp.StatusCode != 201 {
			log.Println("Error("+resp.Status+"): ", resp.Body)
		} else {
			log.Println("Logged Submitted["+resp.Status+"]: ", resp.Body)
		}

		req = nil
		client = nil
		resp = nil

	}
}
