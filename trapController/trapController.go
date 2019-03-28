package trapController

import (
	"github.com/godbus/dbus"
)

const (
	dbusPath   = "/org/cacophony/rs485controller"
	dbusDest   = "org.cacophony.rs485controller"
	methodBase = "org.cacophony.rs485controller"
)

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
	err = obj.Call(methodBase+".DigitalPinRead", 0, pin, update).Store(&res.Value)
	return
}

func DigitalPinReadAll(update bool) (res []DigitalPin, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	var names []string
	var outputs, values []bool

	err = obj.Call(methodBase+".DigitalPinReadAll", 0, update).Store(&names, &outputs, &values)
	res = make([]DigitalPin, len(names))
	for i := range res {
		res[i].Name = names[i]
		res[i].Output = outputs[i]
		res[i].Value = boolToUint16(values[i])
	}
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

func ActuatorReadAll(update bool) (res []Actuator, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	var names []string
	var values []uint16
	var extended, retracted []bool
	err = obj.Call(methodBase+".ActuatorReadAll", 0, update).Store(&names, &values, &extended, &retracted)
	res = make([]Actuator, len(names))
	for i := range res {
		res[i].Name = names[i]
		res[i].Value = values[i]
		res[i].Extended = boolToUint16(extended[i])
		res[i].Retracted = boolToUint16(retracted[i])
	}
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

func ServoReadAll(update bool) (res []Servo, err error) {
	obj, err := getDbusObj()
	if err != nil {
		return
	}
	var names []string
	var values []uint16
	err = obj.Call(methodBase+".ServoReadAll", 0, update).Store(&names, &values)
	res = make([]Servo, len(names))
	for i := range res {
		res[i].Name = names[i]
		res[i].Value = values[i]
	}
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

func boolToUint16(val bool) uint16 {
	if val {
		return 1
	}
	return 0
}
