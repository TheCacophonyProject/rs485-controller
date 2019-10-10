package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	trapController "github.com/TheCacophonyProject/rs485-controller/trapController"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var version = "<not set>"

const powerPin = "GPIO27"

func main() {

	log.SetFlags(0) // Removes timestamp output
	log.Printf("running version: %s", version)

	// Turn power on to traps
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	restartTrapPower(powerPin)

	//TODO make this an array of devices
	device, err := trapController.NewTrap("/etc/cacophony/rs485-controller.yaml", "/dev/ttyAMA0", 9600, 3000)
	if err != nil {
		log.Println("trap fail")
		log.Fatal(err)
	}
	log.Print(device)

	err = device.Test()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("device test passed")

	log.Println("starting DBUS service")
	if err := startDbusService(*device); err != nil {
		log.Fatal(err)
	}
	log.Println("started DBUS service")

	runtime.Goexit() // Stay on for dbus service
	/*
		select {}

		for {
			err = device.Update()
			if err != nil {
				log.Fatal(err)
			}
			//log.Println(device.Actuators[0].Value)
			//log.Println(device.Actuators[0].Retracted)
			time.Sleep(60 * time.Second)
		}
	*/
}

func restartTrapPower(powerPin string) error {
	log.Println("restarting trap power")
	pin := gpioreg.ByName(powerPin)
	if pin == nil {
		fmt.Errorf("no '%s' pin found", powerPin)
	}
	pin.Out(gpio.Low)
	time.Sleep(time.Second * 4)
	pin.Out(gpio.High)
	time.Sleep(time.Second * 2)
	return nil
}
