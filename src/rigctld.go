package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aherve/gopool"
	"gopkg.in/ini.v1"
)

func rigmain() {
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

	// Worker: rigctld (Frequency & Mode) -> https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
	// - Frequency? https://github.com/wavelog/WaveLogGate/blob/875e996cb49f368b97db4f76010ba05a66fd2402/renderer.js#L207-L240
	// - Mode? https://github.com/wavelog/WaveLogGate/blob/875e996cb49f368b97db4f76010ba05a66fd2402/main.js#L330-L334
	// - Wavelog (POST HTTP request) -> https://stackoverflow.com/a/24455606
	if Config.Rigctld.Enabled {
		pool := gopool.NewPool(1)
		for {
			pool.Add(1)
			go func(pool *gopool.GoPool) {
				defer pool.Done()
				conn, err := net.Dial("tcp", Config.Rigctld.Host+":"+Config.Rigctld.Port)
				if err != nil {
					log.Fatal(err.Error())
					time.Sleep(time.Millisecond * time.Duration(1000))
				} else {
					for {
						_, err := conn.Write([]byte("f\n"))
						if err != nil {
							log.Fatal(err.Error())
							break
						}

						freq, err = bufio.NewReader(conn).ReadString('\n')
						if err != nil {
							log.Fatal(err.Error())
							break
						}

						freq = strings.Trim(freq, "\n")

						_, err = conn.Write([]byte("m\n"))
						if err != nil {
							log.Fatal(err.Error())
							break
						}

						mode, err = bufio.NewReader(conn).ReadString('\n')
						if err != nil {
							log.Fatal(err.Error())
							break
						}

						mode = strings.Trim(mode, "\n")

						_, err = conn.Write([]byte("l RFPOWER\n"))
						if err != nil {
							log.Fatal(err.Error())
							break
						}

						power, err = bufio.NewReader(conn).ReadString('\n')
						if err != nil {
							log.Fatal(err.Error())
							break
						}

						power = strings.Trim(power, "\n")

						if power != "RPRT -1" {
							// power
						}

						current_time := time.Now()
						var jsonStr = []byte(`{"radio":"` + Config.Wavelog.Radio + `","key":"` + Config.Wavelog.Key + `","frequency":"` + freq + `","mode":"` + mode + `","timestamp":"` + current_time.Format("2006/01/02 15:04:05") + `"}`)
						req, err := http.NewRequest("POST", Config.Wavelog.URL+"/api/radio", bytes.NewBuffer(jsonStr))
						req.Header.Set("Content-Type", "application/json")

						client := &http.Client{}
						resp, err := client.Do(req)

						if err != nil {
							log.Fatal(err.Error())
						}

						defer resp.Body.Close()

						if resp.StatusCode != 200 {
							fmt.Println("Error("+resp.Status+"): ", resp.Body)
						}

						req = nil
						client = nil
						resp = nil

						time.Sleep(time.Millisecond * time.Duration(1000))
					}
				}
			}(pool)
		}
		pool.Wait()
	}
}
