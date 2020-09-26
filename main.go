package main

import (
	"flag"

	"github.com/cessien/cubbybot/gui"
)

var (
	radioEnabled = flag.Bool("radio", false, "enable radio subsystem")
	controllerEnabled = flag.Bool("controller", false, "enable controller subsystem")
	guiEnabled = flag.Bool("gui", false, "enable gui subsystem")
)

func main() {
	flag.Parse()

	if *guiEnabled {
		gui.InitializeHome()
	}

	if *radioEnabled {
		go radioRoutine()
	}

	if *controllerEnabled {
		go controllerRoutine()
	}

	for {
		select {
		}
	}
}
