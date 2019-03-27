package main

import (
	"flag"
	"log"

	"github.com/TheCacophonyProject/rs485-controller/trapController"
)

var (
	version = "<not set>"
)

func main() {
	log.SetFlags(0) // Removes timestamp output
	log.Printf("running version: %s", version)

	var digitalPinReadAll, actuatorReadAll, servoReadAll, skipUpdate bool
	var digitalPinRead, digitalPinWrite, servoRead, servoWrite, actuatorRead, actuatorWrite string
	var val64 uint64

	flag.BoolVar(&skipUpdate, "skip-update", false, "Won't update values from trap")
	flag.Uint64Var(&val64, "value", 0, "Value when writing")

	flag.StringVar(&digitalPinWrite, "digital-pin-write", "", "Write to a digital pin")
	flag.StringVar(&digitalPinRead, "digital-pin-read", "", "Read a digital pin")
	flag.BoolVar(&digitalPinReadAll, "digital-pin-read-all", false, "Read all digital pins")

	flag.StringVar(&actuatorWrite, "actuator-write", "", "Write to a actuator")
	flag.StringVar(&actuatorRead, "actuator-read", "", "Read a actuator")
	flag.BoolVar(&actuatorReadAll, "actuator-read-all", false, "Read all actuators")

	flag.StringVar(&servoWrite, "servo-write", "", "Write to a servo")
	flag.StringVar(&servoRead, "servo-read", "", "Read a servo")
	flag.BoolVar(&servoReadAll, "servo-read-all", false, "Read all servos")

	flag.Parse()

	val := uint16(val64)

	if digitalPinWrite != "" {
		err := trapController.DigitalPinWrite(digitalPinWrite, val)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("updated")
		return
	}
	if digitalPinRead != "" {
		result, err := trapController.DigitalPinRead(digitalPinRead, !skipUpdate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(result)
		return
	}
	if digitalPinReadAll {
		names, outputs, values, err := trapController.DigitalPinReadAll(!skipUpdate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(names)
		log.Println(outputs)
		log.Println(values)
		return
	}

	if actuatorWrite != "" {
		err := trapController.ActuatorWrite(actuatorWrite, val)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("updated")
		return
	}
	if actuatorRead != "" {
		value, extended, retracted, err := trapController.ActuatorRead(actuatorRead, !skipUpdate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(value)
		log.Println(extended)
		log.Println(retracted)
	}
	if actuatorReadAll {
		names, vals, extended, retracted, err := trapController.ActuatorReadAll(!skipUpdate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(names)
		log.Println(vals)
		log.Println(extended)
		log.Println(retracted)
	}

	if servoWrite != "" {
		err := trapController.ServoWrite(servoWrite, val)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("updated")
		return
	}
	if servoRead != "" {
		value, err := trapController.ServoRead(servoRead, !skipUpdate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(value)
		return
	}
	if servoReadAll {
		names, values, err := trapController.ServoReadAll(!skipUpdate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(names)
		log.Println(values)
		return
	}

}
