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
	servoDoorAngleOpen    = 1
	servoDoorAngleClose   = 179
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
		for {
			if err := s.runSequence(); err != nil {
				s.error(err)
				log.Println(err)
				time.Sleep(10 * time.Second)
			} else {
				break
			}
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

func (s *sequence) runSequence() error {
	if s.running {
		return errors.New("already runnign sequence")
	}
	s.running = true
	defer func() {
		s.running = false
	}()

	//======================START OF SEQUENCE==============================
	s.updateState("reset trap")
	if err := trapController.DigitalPinWrite("12VEnable", 1); err != nil { // Servos need to be powered all the time for the
		return err
	}
	if err := trapController.DigitalPinWrite("6VEnable", 1); err != nil { // Servos need to be powered all the time for the
		return err
	}
	if exit := s.wait(2 * time.Second); exit {
		return nil
	}
	if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
		return nil
	}
	if exit := s.wait(2 * time.Second); exit {
		return nil
	}
	if err := trapController.ServoWrite("Door1", servoDoorAngleOpen); err != nil {
		return nil
	}
	if exit := s.wait(2 * time.Second); exit {
		return nil
	}
	if err := trapController.ServoWrite("Door2", servoDoorAngleClose); err != nil {
		return nil
	}
	if exit := s.wait(2 * time.Second); exit {
		return nil
	}
	// reset spring
	if exit, err := s.resetSpring(); err != nil {
		return err
	} else if exit {
		return nil
	}

	//======================WAIT FOR PIR1==============================
	s.updateState("waiting for PIR1 (first pest)")
	if exit, err := s.waitForDigitalPin("PIR1", 1); err != nil {
		s.error(err)
		return nil
	} else if exit {
		return nil
	}
	s.updateState("PIR1 triggered")

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

	//======================WAIT FOR PIR2==============================
	s.updateState("waiting for PIR2 (waiting for pest to move into holding chamber)")
	if exit, err := s.waitForDigitalPin("PIR2", 1); err != nil {
		s.error(err)
		return nil
	} else if exit {
		return nil
	}

	if err := trapController.ServoWrite("Door1", servoDoorAngleClose); err != nil {
		s.error(err)
		return nil
	}

	if exit := s.wait(2 * time.Second); exit {
		return nil
	}

	if exit, err := s.resetSpring(); err != nil {
		return err
	} else if exit {
		return nil
	}

	i := 1
	for {
		//======================WAIT FOR PIR1==============================
		s.updateState(fmt.Sprintf("waiting for PIR1 (capture pest number %d)", i+1))
		if exit, err := s.waitForDigitalPin("PIR1", 1); err != nil {
			s.error(err)
			return nil
		} else if exit {
			return nil
		}

		// Close trap
		if err := trapController.ServoWrite("Activate", servoAngle2); err != nil {
			s.error(err)
			return nil
		}
		if exit := s.wait(2 * time.Second); exit {
			return nil
		}
		if err := trapController.ServoWrite("Activate", servoAngle1); err != nil {
			s.error(err)
			return nil
		}
		if exit := s.wait(2 * time.Second); exit {
			return nil
		}
		// Open Door2
		if err := trapController.ServoWrite("Door2", servoDoorAngleOpen); err != nil {
			s.error(err)
			return nil
		}
		if exit := s.wait(10 * time.Second); exit {
			return nil
		}

		//======================WAIT FOR PIR3==============================
		s.updateState(fmt.Sprintf("waiting for PIR3 (pest %d to move into last chamber)", i))
		if exit, err := s.waitForDigitalPin("PIR3", 1); err != nil {
			s.error(err)
			return nil
		} else if exit {
			return nil
		}

		// Close Door2
		if err := trapController.ServoWrite("Door2", servoDoorAngleClose); err != nil {
			s.error(err)
			return nil
		}

		if exit := s.wait(2 * time.Second); exit {
			return nil
		}

		// Open Door1
		if err := trapController.ServoWrite("Door1", servoDoorAngleOpen); err != nil {
			s.error(err)
			return nil
		}

		if exit := s.wait(10 * time.Second); exit {
			return nil
		}

		//======================WAIT FOR PIR 2==============================
		s.updateState(fmt.Sprintf("waiting for PIR2 (pest %d to move into holding chamber)", i+1))
		if exit, err := s.waitForDigitalPin("PIR2", 1); err != nil {
			s.error(err)
			return nil
		} else if exit {
			return nil
		}

		// Open trap
		s.resetSpring()

		if exit := s.wait(2 * time.Second); exit {
			return nil
		}

		// Clse Door 1
		if err := trapController.ServoWrite("Door1", servoDoorAngleClose); err != nil {
			s.error(err)
			return nil
		}
		i++

		if exit := s.wait(10 * time.Second); exit {
			return nil
		}
	}
}
