package main

import (
	"fmt"
	// "encoding/binary"
	rf24 "github.com/Draal/gorf24"
	//"github.com/cessien/cubbybot/gui"
)

const CE_PIN uint16 = 22
const CS_PIN uint16 = 24

func main() {
	//gui.InitializeHome()

	var pipe uint64 = 0x7878787878
	fmt.Println("1")

	radio := rf24.New(CE_PIN, CS_PIN, 8000000)
	fmt.Println("1")
	defer radio.Delete()
	fmt.Println("1")

	radio.Begin()
	fmt.Println("1")
	radio.SetRetries(15, 15)
	fmt.Println("1")
	radio.SetDataRate(rf24.RATE_1MBPS)
	fmt.Println("1")
	radio.SetAutoAck(false)
	fmt.Println("1")
	radio.EnableDynamicPayloads()
	fmt.Println("1")
	radio.SetChannel(1)
	fmt.Println("1")

	radio.OpenReadingPipe(0, pipe)
	fmt.Println("1")
	
	radio.PrintDetails()
	fmt.Println("1")

	//radio.StartListening()
		text := "hello world"
		if(radio.Write([]byte(text),4)) {
			fmt.Println("delivered")
		} else {
			fmt.Println("could not deliver")
		}

	
		// if radio.Available() {
		// 	data := radio.Read(4)
		// 	fmt.Printf("data: %v\n", data)

		// 	payload := binary.LittleEndian.Uint32(data)
		// 	fmt.Printf("Received %v\n", payload)
		// }

}
