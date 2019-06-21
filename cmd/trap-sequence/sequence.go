package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/TheCacophonyProject/rs485-controller/trapController"
)

const (
	resetActuatorTime = time.Second * 40
	servoAngle1       = 10
	servoAngle2       = 100
)

type sequence struct {
	state   string
	quit    chan bool
	running bool
}

func getSequence() *sequence {
	return &sequence{
		state: "Nothing",
		quit:  make(chan bool),
	}
}

func (s *sequence) Start() error {
	if s.running {
		return errors.New("already runnign sequence")
	}
	go func() {
		s.runSequence()
	}()
	return nil
}

func (s *sequence) updateState(state string) {
	s.state = state
	log.Println("new state: ", state)
}

func (s *sequence) error(err error) {
	s.updateState(fmt.Sprintf("erorr in sequence: ", err.Error()))
}

func (s *sequence) Stop() {
	s.quit <- true
	trapController.ActuatorWrite("Reset", 2)
}

func (s *sequence) runSequence() error {
	if s.running {
		return errors.New("already runnign sequence")
	}
	s.running = true
	defer func() {
		s.running = false
	}()
	s.updateState("Starting")
	s.updateState("reset servo")
	// Reset servo
	if err := trapController.DigitalPinWrite("6VEnable", 1); err != nil {
		s.error(err)
		return nil
	}
	if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
		s.error(err)
		return nil
	}

	timeout := time.After(time.Second)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}

	s.updateState("reset spring")
	// Reset trap
	if err := trapController.DigitalPinWrite("12VEnable", 1); err != nil {
		s.error(err)
		return nil
	}
	if err := trapController.ActuatorWrite("Reset", 2); err != nil {
		s.error(err)
		return nil
	}
	timeout = time.After(resetActuatorTime)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}
	if err := trapController.ActuatorWrite("Reset", 0); err != nil {
		s.error(err)
		return nil
	}
	timeout = time.After(time.Second)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}
	if err := trapController.ActuatorWrite("Reset", 1); err != nil {
		s.error(err)
		return nil
	}
	timeout = time.After(resetActuatorTime)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}
	if err := trapController.DigitalPinWrite("12VEnable", 0); err != nil {
		s.error(err)
		return nil
	}

	s.updateState("waiting for PIR1")
	for {
		val, err := trapController.DigitalPinRead("PIR1", true)
		if err != nil {
			s.error(err)
			return nil
		}
		if val.Value == 1 {
			s.updateState("PIR1 triggered")
			break
		}
		timeout = time.After(time.Millisecond * 200)
		select {
		case <-s.quit:
			s.updateState("quit")
			return nil
		case <-timeout:
		}
	}
	s.updateState("triggering trap")
	if err := trapController.ServoWrite("Activate", servoAngle2); err != nil {
		s.error(err)
		return nil
	}

	timeout = time.After(time.Second)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}

	if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
		s.error(err)
		return nil
	}

	timeout = time.After(time.Second * 10)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}

	s.updateState("waiting for PIR2")
	for {
		val, err := trapController.DigitalPinRead("PIR2", true)
		if err != nil {
			s.error(err)
			return nil
		}
		if val.Value == 1 {
			break
		}
		timeout = time.After(time.Millisecond * 200)
		select {
		case <-s.quit:
			s.updateState("quit")
			return nil
		case <-timeout:
		}
	}

	s.updateState("reset spring")
	// Reset trap
	if err := trapController.DigitalPinWrite("12VEnable", 1); err != nil {
		s.error(err)
		return nil
	}
	if err := trapController.ActuatorWrite("Reset", 2); err != nil {
		s.error(err)
		return nil
	}
	timeout = time.After(resetActuatorTime)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}
	if err := trapController.ActuatorWrite("Reset", 0); err != nil {
		s.error(err)
		return nil
	}
	timeout = time.After(time.Second)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}
	if err := trapController.ActuatorWrite("Reset", 1); err != nil {
		s.error(err)
		return nil
	}
	timeout = time.After(resetActuatorTime)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}
	if err := trapController.DigitalPinWrite("12VEnable", 0); err != nil {
		s.error(err)
		return nil
	}

	s.updateState("waiting for PIR1")
	for {
		val, err := trapController.DigitalPinRead("PIR1", true)
		if err != nil {
			s.error(err)
			return nil
		}
		if val.Value == 1 {
			log.Println("PIR1 triggered")
			break
		}
		timeout = time.After(time.Millisecond * 200)
		select {
		case <-s.quit:
			s.updateState("quit")
			return nil
		case <-timeout:
		}
	}
	if err := trapController.ServoWrite("Activate", servoAngle2); err != nil {
		s.error(err)
		return nil
	}

	timeout = time.After(time.Second)
	select {
	case <-s.quit:
		s.updateState("quit")
		return nil
	case <-timeout:
	}

	if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
		s.error(err)
		return nil
	}
	return nil
}
