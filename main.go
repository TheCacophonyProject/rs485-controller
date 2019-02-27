package main

import (
	"encoding/binary"
	"errors"
	"log"
	"time"
)

func main() {

	trap, err := NewTrap("./trap-config.yaml", "/dev/ttyAMA0", 9600, 3000)
	if err != nil {
		log.Println("trap fail")
		log.Fatal(err)
	}
	log.Print(trap)

	err = trap.Test()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("device test passed")

	for {
		err = trap.Update()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(trap.DigitalPins[0].Name)
		log.Println(trap.DigitalPins[0].Value)
		log.Println(trap.DigitalPins[1].Name)
		log.Println(trap.DigitalPins[1].Value)
		log.Println(trap.DigitalPins[2].Name)
		log.Println(trap.DigitalPins[2].Value)
		log.Println(trap.DigitalPins[3].Name)
		log.Println(trap.DigitalPins[3].Value)

		log.Println(trap.Servos[0].Name)
		log.Println(trap.Servos[0].Value)
		log.Println(trap.Servos[1].Name)
		log.Println(trap.Servos[1].Value)

		log.Println(trap.Actuators[0].Name)
		log.Println(trap.Actuators[0].Value)
		log.Println(trap.Actuators[0].Extended)
		log.Println(trap.Actuators[0].Retracted)
		time.Sleep(time.Second)
	}
}

func Unit16fromBytes(bytes []byte) ([]uint16, error) {
	l := len(bytes)
	if l%2 == 1 {
		return nil, errors.New("length must be divisible by 2")
	}
	res := make([]uint16, l/2)
	for i := 0; i < l/2; i++ {
		res[i] = binary.BigEndian.Uint16(bytes[i*2 : i*2+2])
	}
	return res, nil
}
