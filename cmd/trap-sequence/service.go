package main

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	dbusName = "org.cacophony.trapsequence"
	dbusPath = "/org/cacophony/trapsequence"
)

type service struct {
	sequence *sequence
}

func startDbusService() error {
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
		sequence: getSequence(),
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

func dbusError(name string, body string) *dbus.Error {
	return &dbus.Error{
		Name: dbusName + "." + name,
		Body: []interface{}{body},
	}
}

// StartSequence will start the trap sequence of the given name
func (s service) StartSequence() *dbus.Error {
	err := s.sequence.Start()
	if err != nil {
		return dbusError("StartSequence", err.Error())
	}
	return nil
}

// GetState will return a string describing the current state of the sequence
func (s service) GetState() (string, *dbus.Error) {
	return s.sequence.state, nil
}

// StopSequence will stop the trap sequence
func (s service) StopSequence() *dbus.Error {
	err := s.sequence.Stop()
	if err != nil {
		return dbusError("StopSequence", err.Error())
	}
	return nil
}
