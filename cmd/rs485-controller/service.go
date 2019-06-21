package main

import (
	"errors"

	trapController "github.com/TheCacophonyProject/rs485-controller/trapController"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	dbusName = "org.cacophony.rs485controller"
	dbusPath = "/org/cacophony/rs485controller"
)

type service struct {
	device trapController.Trap
}

func startDbusService(device trapController.Trap) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	reply, err := conn.RequestName(dbusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("name already taken")
	}

	s := &service{
		device: device,
	}
	conn.Export(s, dbusPath, dbusName)
	conn.Export(genIntrospectable(s), dbusPath, "org.freedesktop.DBus.Introspectable")
	return nil
}

func genIntrospectable(v interface{}) introspect.Introspectable {
	node := &introspect.Node{
		Interfaces: []introspect.Interface{{
			Name:    dbusName,
			Methods: introspect.Methods(v),
		}},
	}
	return introspect.NewIntrospectable(node)
}

func (s service) DigitalPinWrite(name string, value uint16) *dbus.Error {
	d, err := s.device.GetDigitalPin(name)
	if err != nil {
		return dbusError("DigitalPinWrite", err.Error())
	}
	err = s.device.Write(d.Address, value)
	if err != nil {
		return dbusError("DigitalPinWrite", err.Error())
	}
	return nil
}

func (s service) DigitalPinRead(pinName string, update bool) (uint16, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return 0, dbusError("DigitalPinRead", err.Error())
		}
	}

	res, err := s.device.ReadDigitalPin(pinName)
	if err != nil {
		return 0, dbusError("DigitalPinRead", err.Error())
	}
	return res, nil
}

func (s service) DigitalPinReadAll(update bool) ([]string, []bool, []uint16, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return nil, nil, nil, dbusError("DigitalPinReadAll", err.Error())
		}
	}

	l := len(s.device.DigitalPins)
	names := make([]string, l)
	output := make([]bool, l)
	values := make([]uint16, l)

	for i, d := range s.device.DigitalPins {
		val, err := s.device.ReadDigitalPin(d.Name)
		if err != nil {
			return nil, nil, nil, dbusError("DigitalPinReadAll", err.Error())
		}
		names[i] = d.Name
		output[i] = d.Output
		values[i] = val
	}
	return names, output, values, nil
}

func (s service) ActuatorWrite(actuatorName string, value uint16) *dbus.Error {
	if err := s.device.WriteActuator(actuatorName, value); err != nil {
		return dbusError("ActuatorWrite", err.Error())
	}
	return nil
}

func (s service) ActuatorRead(actuatorName string, update bool) (uint16, bool, bool, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return 0, false, false, dbusError("ActuatorRead", err.Error())
		}
	}
	val, extended, retracted, err := s.device.ReadActuator(actuatorName)
	if err != nil {
		return 0, false, false, dbusError("ActuatorRead", err.Error())
	}
	return val, extended, retracted, nil
}

func (s service) ActuatorReadAll(update bool) ([]string, []uint16, []bool, []bool, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return nil, nil, nil, nil, dbusError("ActuatorReadAll", err.Error())
		}
	}
	l := len(s.device.Actuators)
	name := make([]string, l)
	state := make([]uint16, l)
	exnteded := make([]bool, l)
	retracted := make([]bool, l)

	for i, a := range s.device.Actuators {
		val, e, r, err := s.device.ReadActuator(a.Name)
		if err != nil {
			return nil, nil, nil, nil, dbusError("ActuatorReadAll", err.Error())
		}
		name[i] = a.Name
		state[i] = val
		exnteded[i] = e
		retracted[i] = r
	}
	return name, state, exnteded, retracted, nil
}

func (s service) ServoWrite(servoName string, value uint16) *dbus.Error {
	if err := s.device.WriteServo(servoName, value); err != nil {
		return dbusError("ServoWrite", err.Error())
	}
	return nil
}

func (s service) ServoRead(servoName string, update bool) (uint16, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return 0, dbusError("ServoRead", err.Error())
		}
	}
	val, err := s.device.ReadServo(servoName)
	if err != nil {
		return 0, dbusError("ServoRead", err.Error())
	}
	return val, nil
}

func (s service) ServoReadAll(update bool) ([]string, []uint16, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return nil, nil, dbusError("ServosReadAll", err.Error())
		}
	}

	l := len(s.device.Servos)
	names := make([]string, l)
	values := make([]uint16, l)

	for i, servo := range s.device.Servos {
		val, err := s.device.ReadServo(servo.Name)
		if err != nil {
			return nil, nil, dbusError("ServosReadAll", err.Error())
		}
		names[i] = servo.Name
		values[i] = val
	}
	return names, values, nil
}

func dbusError(name string, body string) *dbus.Error {
	return &dbus.Error{
		Name: dbusName + "." + name,
		Body: []interface{}{body},
	}
}
