package main

import (
	"time"
)

func main() {
	//go radioRoutine()

	go controllerRoutine()

	time.Sleep(1 * time.Hour)
}
