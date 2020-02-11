package main

import (
	"time"

	"github.com/TheCacophonyProject/rs485-controller/trapController"
)

func (s *sequence) waitForDigitalPin(pin string, value uint16) (bool, error) {
	for {
		val, err := trapController.DigitalPinRead(pin, true)
		if err != nil {
			return false, err
		}
		if val.Value == value {
			return false, nil
		}
		if exit := s.wait(200 * time.Millisecond); exit {
			return true, nil
		}
	}
}

func (s *sequence) wait(d time.Duration) bool {
	// Check recording window
	// Have s.quit be called at end of window
	timeout := time.After(d)
	select {
	case <-s.quit:
		s.updateState("quit")
		return true
	case <-timeout:
		return false
	}
}

func (s *sequence) resetSpring() (bool, error) {
	// reset trap spring
	s.updateState("reset spring")
	if exit := s.wait(time.Second); exit {
		return true, nil
	}
	if err := trapController.DigitalPinWrite("12VEnable", 1); err != nil {
		return false, err
	}
	if exit := s.wait(time.Second); exit {
		return true, nil
	}
	if err := trapController.ActuatorWrite("Reset", 0); err != nil {
		return false, err
	}
	if exit := s.wait(time.Second); exit {
		return true, nil
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

	if err := trapController.ActuatorWrite("Reset", 0); err != nil {
		return false, err
	}

	if exit := s.wait(time.Second); exit {
		return true, nil
	}

	return false, nil
}
