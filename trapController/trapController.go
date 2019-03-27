package trapController

import (
	"github.com/godbus/dbus"
)

const (
	dbusPath   = "/org/cacophony/rs485controller"
	dbusDest   = "org.cacophony.rs485controller"
	methodBase = "org.cacophony.rs485controller"
)

type DigitalPin struct {
	name   string
	value  bool
	output bool
}

func DigitalPinWrite(pin string, val uint16) error {
	obj, err := getDbusObj()
	if err != nil {
		return err
	}
	return obj.Call(methodBase+".DigitalPinWrite", 0, pin, val).Store()
}

func DigitalPinRead(pin string, update bool) (res DigitalPin, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	err = obj.Call(methodBase+".DigitalPinRead", 0, pin, update).Store(&res.value)
	return
}

func DigitalPinReadAll(update bool) (names []string, outputs, values []bool, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	err = obj.Call(methodBase+".DigitalPinReadAll", 0, update).Store(&names, &outputs, &values)
	return
}

func ActuatorWrite(name string, value uint16) error {
	obj, err := getDbusObj()
	if err != nil {
		return err
	}
	return obj.Call(methodBase+".ActuatorWrite", 0, name, value).Store()
}

func ActuatorRead(name string, update bool) (value uint16, extended, retracted bool, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	err = obj.Call(methodBase+".ActuatorRead", 0, name, update).Store(&value, &extended, &retracted)
	return
}

func ActuatorReadAll(update bool) (names []string, values []uint16, extended, retracted []bool, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	err = obj.Call(methodBase+".ActuatorReadAll", 0, update).Store(&names, &values, &extended, &retracted)
	return
}

func ServoWrite(name string, val uint16) error {
	obj, err := getDbusObj()
	if err != nil {
		return err
	}
	return obj.Call(methodBase+".ServoWrite", 0, name, val).Store()
}

func ServoRead(name string, update bool) (value uint16, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	err = obj.Call(methodBase+".ServoRead", 0, name, update).Store(&value)
	return
}

func ServoReadAll(update bool) (names []string, values []uint16, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	err = obj.Call(methodBase+".ServoReadAll", 0, update).Store(&names, &values)
	return
}

func getDbusObj() (dbus.BusObject, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	obj := conn.Object(dbusDest, dbusPath)
	return obj, nil
}
