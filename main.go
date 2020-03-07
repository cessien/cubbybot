package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	rf24 "github.com/cessien/gorf24"
	//"github.com/cessien/cubbybot/gui"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

const CE_PIN uint16 = 22     // RPi GPIO 22
const CS_PIN uint16 = 8      // RPi CE0 CSN, GPIO 08
const SPI_SPEED_HZ = 8000000 // 8Mhz

func main() {
	//gui.InitializeHome()
	var role string = "pong_back"
	reader := bufio.NewReader(os.Stdin)

	var pipes []uint64 = []uint64{0x7171717171, 0x3434343434}

	fmt.Println("example getting started")

	// create the Radio
	radio := rf24.New(CE_PIN, CS_PIN, SPI_SPEED_HZ)

	// clean up c++ references
	defer radio.Delete()

	// setup and configure rf radio
	radio.Begin()

	// optional - increase the delay of retries, and the total number of retries
	radio.SetRetries(15, 15)
	// dump the configuration of the rf unit for debugging
	radio.PrintDetails()

	/*** choose a role ***/
	fmt.Println("Choose a role: Enter 0 for pong_back, 1 for ping_out (CTRL+C to exit)")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if !strings.Contains(text, "0") && !strings.Contains(text, "1") {
		panic(errors.New("Need to enter 0 or 1"))
	}

	if strings.Contains(text, "1") {
		role = "ping_out"
		radio.OpenWritingPipe(pipes[0])
		radio.OpenReadingPipe(1, pipes[1])
	} else {
		radio.OpenWritingPipe(pipes[1])
		radio.OpenReadingPipe(1, pipes[0])
	}

	//radio.SetChannel(1)

	radio.StartListening()

	// forever loop
	for {
		if role == "ping_out" {
			// first, stop listening so we can talk
			radio.StopListening()

			// take the time, and send it. will block until complete

			fmt.Println("Now sending...")
			m := Message{
				Type: "test",
				Data: fmt.Sprintf("Pi time: %d", time.Now().Unix()),
			}

			b, _ := json.Marshal(m)

			ok := radio.Write(b, 8)

			if !ok {
				fmt.Println("failed.")
			}

			// now, continue listening
			radio.StartListening()

			// wait here until we get a response, or timeout (250ms)
			startedWaitingAt := time.Now()
			var timeout bool = false
			for !radio.Available() && !timeout {
				if time.Now().Sub(startedWaitingAt) > 200*time.Millisecond {
					timeout = true
				}
			}

			// describe the results
			if timeout {
				fmt.Println("Failed, response timed out")
			} else {
				// grab the response, compare, and send to debugging spew
				var data []byte = radio.Read(255)
				m := Message{}
				err = json.Unmarshal(data, &m)
				if err != nil {
					fmt.Printf("Error parsing message: %v\n", err)
				}

				// spew it
				fmt.Printf("Got response[%s]: %s\n", m.Type, m.Data)
			}
			time.Sleep(1 * time.Second)
		}

		/*
		 * Pong back role. Receive each packet, dump it out, and send it back
		 */
		if role == "pong_back" {
			// if there is data ready
			if radio.Available() {
				// dump the payloads until we've gotten everything
				var data []byte = radio.Read(255)
				m := Message{}
				err = json.Unmarshal(data, &m)
				if err != nil {
					fmt.Printf("Error parsing message: %v\n", err)
				}
				radio.StopListening()

				m.Data = "Heck yeah!"
				b, _ := json.Marshal(m)

				radio.Write(b, 255)

				// now, resume listening so we can catch the next packets
				radio.StartListening()

				// spew it
				fmt.Printf("Got payload[%s] %s\n", m.Type, m.Data)

				time.Sleep(925 * time.Millisecond) // delay after payload responded to, minimize RPi CPU time
			}

		}
	}
}
