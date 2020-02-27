package main

import (
	rf24 "github.com/Draal/gorf24"
	"github.com/cessien/cubbybot/gui"
)

func main() {
	gui.InitializeHome()

	radio := rf24.New("/dev/spidev0.0", 8000000, 25)
	defer r.Delete()

	radio.Begin()
	radio.SetRetries(15, 15)

}
