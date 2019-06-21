package main

import (
	"log"
	"time"
)

var (
	version = "<not set>"
)

func main() {
	err := runMain()
	if err != nil {
		log.Fatal(err)
	}
}

func runMain() error {
	log.SetFlags(0) // Removes timestamp output
	log.Printf("running version: %s", version)

	log.Println("starting service")
	if err := startDbusService(); err != nil {
		return err
	}

	for {
		time.Sleep(time.Second)
	}
	/*
		log.Println("reset servo")
		// Reset servo
		if err := trapController.DigitalPinWrite("6VEnable", 1); err != nil {
			return err
		}
		if err := trapController.ServoWrite("Activate", 10); err != nil {
			return err
		}
		time.Sleep(time.Second)

		log.Println("reset spring")
		// Reset trap
		if err := trapController.DigitalPinWrite("12VEnable", 1); err != nil {
			return err
		}
		if err := trapController.ActuatorWrite("Reset", 2); err != nil {
			return err
		}
		time.Sleep(time.Second * 40)
		if err := trapController.ActuatorWrite("Reset", 0); err != nil {
			return err
		}
		time.Sleep(time.Second)
		if err := trapController.ActuatorWrite("Reset", 1); err != nil {
			return err
		}
		time.Sleep(time.Second * 40)
		if err := trapController.DigitalPinWrite("12VEnable", 0); err != nil {
			return err
		}

		log.Println("waiting for PIR1")
		for {
			val, err := trapController.DigitalPinRead("PIR1", true)
			if err != nil {
				return err
			}
			if val.Value == 1 {
				log.Println("PIR1 triggered")
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		log.Println("triggering trap")
		if err := trapController.ServoWrite("Activate", 140); err != nil {
			return err
		}

		time.Sleep(time.Second)

		if err := trapController.ServoWrite("Activate", 10); err != nil {
			return err
		}

		time.Sleep(time.Second * 10)

		log.Println("waiting for PIR2")
		for {
			val, err := trapController.DigitalPinRead("PIR2", true)
			if err != nil {
				return err
			}
			if val.Value == 1 {
				break
			}
			time.Sleep(time.Millisecond * 200)
		}

		log.Println("reset spring")
		// Reset trap
		if err := trapController.DigitalPinWrite("12VEnable", 1); err != nil {
			return err
		}
		if err := trapController.ActuatorWrite("Reset", 2); err != nil {
			return err
		}
		time.Sleep(time.Second * 40)
		if err := trapController.ActuatorWrite("Reset", 0); err != nil {
			return err
		}
		time.Sleep(time.Second)
		if err := trapController.ActuatorWrite("Reset", 1); err != nil {
			return err
		}
		time.Sleep(time.Second * 40)
		if err := trapController.DigitalPinWrite("12VEnable", 0); err != nil {
			return err
		}

		log.Println("waiting for PIR1")
		for {
			val, err := trapController.DigitalPinRead("PIR1", true)
			if err != nil {
				return err
			}
			if val.Value == 1 {
				log.Println("PIR1 triggered")
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		if err := trapController.ServoWrite("Activate", 140); err != nil {
			return err
		}

		time.Sleep(time.Second)

		if err := trapController.ServoWrite("Activate", 10); err != nil {
			return err
		}

		return nil
	*/
}
