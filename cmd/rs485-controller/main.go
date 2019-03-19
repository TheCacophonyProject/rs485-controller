package main

import (
	"log"
)

var version = "<not set>"

func main() {

	log.SetFlags(0) // Removes timestamp output
	log.Printf("running version: %s", version)

	//TODO make this an array of devices
	device, err := NewTrap("./trap-config.yaml", "/dev/ttyAMA0", 9600, 3000)
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

	for {
		err = device.Update()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(device.Actuators[0].Value)
		log.Println(device.Actuators[0].Retracted)
		//time.Sleep(1 * time.Second)
	}
}
