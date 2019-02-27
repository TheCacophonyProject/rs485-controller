package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/goburrow/modbus"
	yaml "gopkg.in/yaml.v2"
)

const deviceTypeID = 1

type Trap struct {
	Name          string        `yaml:"name"`
	Version       uint16        `yaml:"version"`
	RS485id       byte          `yaml:"rs485-id"`
	DigitalPins   []*DigitalPin `yaml:"digital-pins"`
	Servos        []*Servo      `yaml:"servos"`
	Actuators     []*Actuator   `yaml:"actuators"`
	handler       *modbus.RTUClientHandler
	updateDetails updateDetails
}

type updateDetails struct {
	maxAddress uint16
	minAddress uint16
	Values     []*uint16
	Address    []uint16
}

type device interface {
	GetName() string
	GetAddress() uint16
}

type DigitalPin struct {
	Name    string `yaml:"name"`
	Address uint16 `yaml:"address"`
	Output  bool   `yaml:"output"`
	Value   uint16
}

type Servo struct {
	Name    string `yaml:"name"`
	Address uint16 `yaml:"address"`
	Value   uint16
}

type Actuator struct {
	Name             string `yaml:"name"`
	Address          uint16 `yaml:"address"`
	ExtendedAddress  uint16 `yaml:"extended-address"`
	RetractedAddress uint16 `yaml:"retracted-address"`
	Value            uint16
	Extended         uint16
	Retracted        uint16
}

func NewTrap(filename string, serialPort string, baudRate int, timeout int) (*Trap, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	trap := &Trap{}
	err = yaml.Unmarshal(buf, trap)
	if err != nil {
		return nil, err
	}
	trap.handler = modbus.NewRTUClientHandler(serialPort)
	trap.handler.BaudRate = baudRate
	trap.handler.SlaveId = trap.RS485id
	trap.handler.DataBits = 8
	trap.handler.Parity = "N"
	trap.handler.StopBits = 1
	trap.handler.Timeout = time.Duration(timeout) * time.Millisecond
	trap.makeupdateDetails()

	log.Println(trap.updateDetails)

	return trap, nil
}

// This will extract data to make updating easier and faster
func (t *Trap) makeupdateDetails() {
	l := len(t.DigitalPins) + len(t.Servos) + len(t.Actuators)*3

	t.updateDetails = updateDetails{
		maxAddress: 0,
		minAddress: math.MaxUint16,
		Values:     make([]*uint16, l),
		Address:    make([]uint16, l),
	}
	i := 0

	for _, digitalPin := range t.DigitalPins {
		t.updateDetails.Values[i] = &digitalPin.Value
		t.updateDetails.Address[i] = digitalPin.Address
		i++
	}

	for _, servo := range t.Servos {
		t.updateDetails.Values[i] = &servo.Value
		t.updateDetails.Address[i] = servo.Address
		i++
	}

	for _, actuator := range t.Actuators {
		t.updateDetails.Values[i] = &actuator.Value
		t.updateDetails.Address[i] = actuator.Address
		i++
		t.updateDetails.Values[i] = &actuator.Extended
		t.updateDetails.Address[i] = actuator.ExtendedAddress
		i++
		t.updateDetails.Values[i] = &actuator.Retracted
		t.updateDetails.Address[i] = actuator.RetractedAddress
		i++
	}
	log.Println("asdasd")
	log.Println(&t.Actuators[0].Retracted)

	for _, val := range t.updateDetails.Address {
		t.updateDetails.maxAddress = uint16Max(t.updateDetails.maxAddress, val)
		t.updateDetails.minAddress = uint16Min(t.updateDetails.minAddress, val)
	}
}

// Test will check that it can connect and that the deviceTypeID and version matches
func (t *Trap) Test() error {
	res, err := t.read(0, 2)
	if err != nil {
		return err
	}
	if res[0] != deviceTypeID {
		return fmt.Errorf("invald deviceTypeID, got %d, expecting, %d", res[0], deviceTypeID)
	}
	if res[1] != t.Version {
		return fmt.Errorf("invald Version, got %d, expecting, %d", res[1], t.Version)
	}

	return nil
}

// Update will read the values from the device
func (t *Trap) Update() error {
	res, err := t.read(t.updateDetails.minAddress, t.updateDetails.maxAddress-t.updateDetails.minAddress+1)
	if err != nil {
		return err
	}
	log.Println(res)
	for i, address := range t.updateDetails.Address {
		*t.updateDetails.Values[i] = res[address-t.updateDetails.minAddress]
	}
	return nil
}

func (t Trap) Write(address uint16, value uint16) error {
	return nil
}

func (t *Trap) read(start uint16, len uint16) ([]uint16, error) {
	err := t.handler.Connect()
	if err != nil {
		return nil, err
	}
	defer t.handler.Close()
	client := modbus.NewClient(t.handler)
	holdingResults, err := client.ReadHoldingRegisters(start, len)
	if err != nil {
		return nil, err
	}
	return Unit16fromBytes(holdingResults)
}

func (t Trap) String() string {
	res := fmt.Sprintf("Trap '%s'\n", t.Name)
	res += fmt.Sprintf("Version: %d\n", t.Version)
	res += fmt.Sprintf("RS485 ID: %d\n", t.RS485id)
	res += fmt.Sprintf("Digital Pins:\n")
	for _, digitalPin := range t.DigitalPins {
		res += "\t" + digitalPin.String()
	}
	res += fmt.Sprintf("Servos:\n")
	for _, servo := range t.Servos {
		res += "\t" + servo.String()
	}
	res += fmt.Sprintf("Actuators:\n")
	for _, actuator := range t.Actuators {
		res += "\t" + actuator.String()
	}
	return res
}

func (d DigitalPin) String() string {
	return fmt.Sprintf("DigitalPin: '%s', Address: '%d', Output: '%t'\n", d.Name, d.Address, d.Output)
}

func (s Servo) String() string {
	return fmt.Sprintf("Servo: `%s`, Address: `%d`\n", s.Name, s.Address)
}

func (a Actuator) String() string {
	return fmt.Sprintf("Actuator: `%s`, Address: `%d`, Extended Address: `%d`, Retracted Address: `%d`\n",
		a.Name, a.Address, a.ExtendedAddress, a.RetractedAddress)
}

func uint16Max(a uint16, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

func uint16Min(a uint16, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
