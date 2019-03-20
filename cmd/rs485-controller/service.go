package main

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	dbusName = "org.cacophony.rs485controller"
	dbusPath = "/org/cacophony/rs485controller"
)

type service struct {
	device Trap
}

func startDbusService(device Trap) error {
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

func (s service) WriteDigitalPin(name string, value uint16) *dbus.Error {
	d, err := s.device.GetDigitalPin(name)
	if err != nil {
		return dbusError("WriteDigitalPin", err.Error())
	}
	err = s.device.Write(d.Address, value)
	if err != nil {
		return dbusError("WriteDigitalPin", err.Error())
	}
	return nil
}

func (s service) ReadDigitalPin(pinName string, update bool) (bool, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return false, dbusError(".Read", err.Error())
		}
	}

	res, err := s.device.ReadDigitalPin(pinName)
	if err != nil {
		return false, dbusError(".Read", err.Error())
	}
	return res, nil
}

func (s service) ReadAllDigitalPins(update bool) ([]string, []bool, []bool, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return nil, nil, nil, dbusError(".Read", err.Error())
		}
	}

	l := len(s.device.DigitalPins)
	names := make([]string, l)
	output := make([]bool, l)
	values := make([]bool, l)

	for i, d := range s.device.DigitalPins {
		val, err := s.device.ReadDigitalPin(d.Name)
		if err != nil {
			return nil, nil, nil, dbusError(".ReadAllDigitalPins", err.Error())
		}
		names[i] = d.Name
		output[i] = d.Output
		values[i] = val
	}
	return names, output, values, nil
}

func (s service) WriteActuator(actuatorName string, value uint16) *dbus.Error {
	if err := s.device.WriteActuator(actuatorName, value); err != nil {
		return dbusError(".WriteActuator", err.Error())
	}
	return nil
}

func (s service) ReadActuator(actuatorName string, update bool) (uint16, bool, bool, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return 0, false, false, dbusError(".Read", err.Error())
		}
	}
	val, extended, retracted, err := s.device.ReadActuator(actuatorName)
	if err != nil {
		return 0, false, false, dbusError("ReadActuator", err.Error())
	}
	return val, extended, retracted, nil
}

func (s service) ReadAllActuators(update bool) ([]string, []uint16, []bool, []bool, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return nil, nil, nil, nil, dbusError(".Read", err.Error())
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
			return nil, nil, nil, nil, dbusError("ReadAllActuators", err.Error())
		}
		name[i] = a.Name
		state[i] = val
		exnteded[i] = e
		retracted[i] = r
	}
	return name, state, exnteded, retracted, nil
}

func (s service) WriteServo(servoName string, value uint16) *dbus.Error {
	if err := s.device.WriteServo(servoName, value); err != nil {
		return dbusError("WriteServo", err.Error())
	}
	return nil
}

func (s service) ReadServo(servoName string, update bool) (uint16, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return 0, dbusError(".Read", err.Error())
		}
	}
	val, err := s.device.ReadServo(servoName)
	if err != nil {
		return 0, dbusError("ReadServo", err.Error())
	}
	return val, nil
}

func (s service) ReadAllServos(update bool) ([]string, []uint16, *dbus.Error) {
	if update {
		if err := s.device.Update(); err != nil {
			return nil, nil, dbusError(".Read", err.Error())
		}
	}

	l := len(s.device.Servos)
	names := make([]string, l)
	values := make([]uint16, l)

	for i, servo := range s.device.Servos {
		val, err := s.device.ReadServo(servo.Name)
		if err != nil {
			return nil, nil, dbusError("ReadAllServos", err.Error())
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
