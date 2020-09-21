package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/ads1x15"
	"periph.io/x/periph/host"
)

type JoySticks struct {
	LeftX  int32
	LeftY  int32
	RightX int32
	RightY int32
}

func (j *JoySticks) ToData() string {
	return fmt.Sprintf("%d,%d,%d,%d", j.LeftX, j.LeftY, j.RightX, j.RightY)
}

func (j *JoySticks) FromData(d string) {
	dataPoints := strings.Split(d, ",")
	j.LeftX = unsafeConvertInt32(dataPoints[0])
	j.LeftY = unsafeConvertInt32(dataPoints[1])
	j.RightX = unsafeConvertInt32(dataPoints[2])
	j.RightY = unsafeConvertInt32(dataPoints[3])
}

func unsafeConvertInt32(s string) int32 {
	i64, _ := strconv.Atoi(s)
	return int32(i64)
}

var (
	joysticks JoySticks
)

func controllerRoutine() {
	joysticks = JoySticks{}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open default I²C bus.
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer bus.Close()

	// Create a new ADS1115 ADC.
	adc, err := ads1x15.NewADS1115(bus, &ads1x15.DefaultOpts)
	if err != nil {
		log.Fatalln(err)
	}

	// Obtain an analog pin from the ADC.
	rx, err := adc.PinForChannel(ads1x15.Channel0, 5*physic.Volt, 100*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		log.Fatalln(err)
	}
	defer rx.Halt()

	ry, err := adc.PinForChannel(ads1x15.Channel1, 5*physic.Volt, 100*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		log.Fatalln(err)
	}
	defer ry.Halt()

	lx, err := adc.PinForChannel(ads1x15.Channel2, 5*physic.Volt, 100*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		log.Fatalln(err)
	}
	defer lx.Halt()

	ly, err := adc.PinForChannel(ads1x15.Channel3, 5*physic.Volt, 100*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		log.Fatalln(err)
	}
	defer ly.Halt()

	// Read values continuously from ADC.
	fmt.Println("Continuous reading")
	crx := rx.ReadContinuous()
	cry := ry.ReadContinuous()
	clx := lx.ReadContinuous()
	cly := ly.ReadContinuous()

	for {
		select {
		case reading1 := <-crx:
			joysticks.LeftX = reading1.Raw
		case reading2 := <-cry:
			joysticks.LeftY = reading2.Raw
		case reading3 := <-clx:
			joysticks.RightX = reading3.Raw
		case reading4 := <-cly:
			joysticks.RightY = reading4.Raw
		}

		// fmt.Printf("left: {%d,%d}, right: {%d,%d}\n", joysticks.LeftX, joysticks.LeftY, joysticks.RightX, joysticks.RightY)
	}
}
