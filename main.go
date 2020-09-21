package main

import (
	"time"

	"github.com/cessien/cubbybot/gui"
)

func main() {
	gui.InitializeHome()
	//go radioRoutine()

	//go controllerRoutine()

	time.Sleep(1 * time.Hour)
}
