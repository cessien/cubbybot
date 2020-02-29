package main

import (
	"fmt"
	"encoding/binary"
	rf24 "github.com/Draal/gorf24"
	//"github.com/cessien/cubbybot/gui"
)

func main() {
	//gui.InitializeHome()

	var pipe uint64 = 0xF0F0F0F0E1

	radio := rf24.New(17, 25, 8000000)
	defer radio.Delete()

	radio.Begin()
	radio.SetRetries(15, 15)
	radio.SetAutoAck(true)
	radio.OpenReadingPipe(1, pipe)
	radio.StartListening()
	radio.PrintDetails()

	for {
		if radio.Available() {
			data := radio.Read(4)
			fmt.Printf("data: %v\n", data)

			payload := binary.LittleEndian.Uint32(data)
			fmt.Printf("Received %v\n", payload)
		}
	}

}
