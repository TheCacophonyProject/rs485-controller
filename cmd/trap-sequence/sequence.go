package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/TheCacophonyProject/rs485-controller/trapController"
)

const (
	resetActuatorDuration = time.Second * 40
	servoAngle1           = 10
	servoAngle2           = 100
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
		if err := s.runSequence(); err != nil {
			s.error(err)
			log.Println(err)
		}
	}()
	return nil
}

func (s *sequence) Stop() error {
	if !s.running {
		return errors.New("sequence already stopped")
	}
	s.quit <- true
	trapController.ActuatorWrite("Reset", 0)
	return nil
}

func (s *sequence) updateState(state string) {
	s.state = state
	log.Println("new state: ", state)
}

func (s *sequence) error(err error) {
	s.updateState(fmt.Sprintf("erorr in sequence: %s", err.Error()))
}

func (s *sequence) wait(d time.Duration) bool {
	timeout := time.After(time.Second)
	select {
	case <-s.quit:
		s.updateState("quit")
		return true
	case <-timeout:
		return false
	}
}

func (s *sequence) resetSpring() (bool, error) {
	s.updateState("reset spring")
	if err := trapController.DigitalPinWrite("12VEnable", 1); err != nil {
		return false, err
	}
	if err := trapController.ActuatorWrite("Reset", 2); err != nil {
		return false, err
	}

	if exit := s.wait(resetActuatorDuration); exit {
		return true, nil
	}

	if err := trapController.ActuatorWrite("Reset", 0); err != nil {
		return false, err
	}

	if exit := s.wait(time.Second); exit {
		return true, nil
	}

	if err := trapController.ActuatorWrite("Reset", 1); err != nil {
		return false, err
	}

	if exit := s.wait(resetActuatorDuration); exit {
		return true, nil
	}
	return false, nil
}

func (s *sequence) runSequence() error {
	if s.running {
		return errors.New("already runnign sequence")
	}
	s.running = true
	defer func() {
		s.running = false
	}()

	s.updateState("reset servo")
	// Reset servo

	if err := trapController.DigitalPinWrite("6VEnable", 1); err != nil {
		return err
	}
	if err := trapController.DigitalPinWrite("EM", 0); err != nil {
		return nil
	}
	if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
		return nil
	}

	if exit := s.wait(time.Second); exit {
		return nil
	}

	if exit, err := s.resetSpring(); err != nil {
		return err
	} else if exit {
		return nil
	}
	if err := trapController.DigitalPinWrite("12VEnable", 0); err != nil {
		return err
	}

	//exit, err := s.waitForTrigger
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

		if exit := s.wait(200 * time.Millisecond); exit {
			return nil
		}
	}

	s.updateState("triggering trap")
	if err := trapController.ServoWrite("Activate", servoAngle2); err != nil {
		s.error(err)
		return nil
	}

	if exit := s.wait(time.Second); exit {
		return nil
	}

	if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
		s.error(err)
		return nil
	}

	if exit := s.wait(10 * time.Second); exit {
		return nil
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
		if exit := s.wait(200 * time.Millisecond); exit {
			return nil
		}
	}

	if err := trapController.DigitalPinWrite("6VEnable", 1); err != nil {
		s.error(err)
		return nil
	}
	if err := trapController.DigitalPinWrite("EM", 1); err != nil {
		s.error(err)
		return nil
	}

	if exit, err := s.resetSpring(); err != nil {
		return err
	} else if exit {
		return nil
	}
	if err := trapController.DigitalPinWrite("12VEnable", 0); err != nil {
		return err
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

		if exit := s.wait(200 * time.Millisecond); exit {
			return nil
		}
	}
	if err := trapController.ServoWrite("Activate", servoAngle2); err != nil {
		s.error(err)
		return nil
	}

	if exit := s.wait(time.Second); exit {
		return nil
	}

	if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
		s.error(err)
		return nil
	}
	return nil
}
