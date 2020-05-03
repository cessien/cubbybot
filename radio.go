package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	// "github.com/cessien/cubbybot/gui"
	rf24 "github.com/cessien/gorf24"
)

type Message struct {
	Type string `json:"t"`
	Data string `json:"d"`
}

func (m *Message) Encode() []byte {
	e := make([]byte, 0, 64)
	e = append(e, []byte(m.Type)...)
	e = append(e, byte('|'))
	e = append(e, []byte(m.Data)...)
	if len(e) < 64 {
		e = append(e, byte('|'))
	} else {
		e[63] = '|'
	}

	return e
}

func (m *Message) Decode(b []byte) error {
	items := bytes.Split(b, []byte("|"))
	if len(items) < 2 {
		return fmt.Errorf("Could not parse %x", b)
	}
	m.Type = string(items[0])
	m.Data = string(items[1])

	return nil
}

const CE_PIN uint16 = 22      // RPi GPIO 22
const CS_PIN uint16 = 8       // RPi CE0 CSN, GPIO 08
const SPI_SPEED_HZ uint32 = 4 // 64Mhz

func initGUI() {
	// gui.InitializeHome()
}

func radioRoutine() {
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

	radio.StartListening()

	// forever loop
	for {
		if role == "ping_out" {
			Ping(&radio)
		}

		/*
		 * Pong back role. Receive each packet, dump it out, and send it back
		 */
		if role == "pong_back" {
			Pong(&radio)
		}
	}
}

func Ping(radio *rf24.R) {
	// first, stop listening so we can talk
	radio.StopListening()

	// take the time, and send it. will block until complete

	fmt.Println("Now sending...")
	m := Message{
		Type: "t",
		Data: fmt.Sprintf("%4d", time.Now().Unix()),
	}

	b := m.Encode()

	var size uint8 = uint8(len(b)) // I grimmace at this, but I'll live
	fmt.Printf("The message sending (size %d) is [%x]: %s\n", size, b, string(b))
	ok := radio.Write(b, size)

	if !ok {
		fmt.Println("failed.")
	}

	// now, continue listening
	radio.StartListening()

	// wait here until we get a response, or timeout (300ms)
	startedWaitingAt := time.Now()
	var timeout bool = false
	for !radio.Available() && !timeout {
		if time.Now().Sub(startedWaitingAt) > 300*time.Millisecond {
			timeout = true
		}
	}

	// describe the results
	if timeout {
		fmt.Println("Failed, response timed out")
	} else {
		// grab the response, compare, and send to debugging spew
		var data []byte = radio.Read(64)
		//data = bytes.Trim(data, "\x00")
		m := Message{}
		err := m.Decode(data)
		if err != nil {
			fmt.Printf("Error parsing response[%x]: %v\n", data, err)
		}

		// spew it
		fmt.Printf("Got response[%s]: %s\n", m.Type, m.Data)
	}
	time.Sleep(5 * time.Millisecond)
}

func Pong(radio *rf24.R) {
	// if there is data ready
	if radio.Available() {
		// dump the payloads until we've gotten everything
		var data []byte = radio.Read(64)
		m := Message{}
		err := m.Decode(data)
		if err != nil {
			fmt.Printf("Error parsing recieved [%x][%s]: %v\n", data, string(data), err)
		}
		radio.StopListening()

		m.Data = fmt.Sprintf("Ack: %s", m.Data)
		b := m.Encode()

		var size uint8 = uint8(len(b)) // I grimmace at this, but I'll live
		radio.Write(b, size)

		// now, resume listening so we can catch the next packets
		radio.StartListening()

		// spew it
		fmt.Printf("Got payload[%s] %s\n", m.Type, m.Data)
	}

	time.Sleep(5 * time.Millisecond) // delay after payload responded to, minimize RPi CPU time
}