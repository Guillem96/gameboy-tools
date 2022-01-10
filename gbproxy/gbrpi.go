package gbproxy

import (
	"fmt"
	"os"
	"time"

	"github.com/Guillem96/gameboy-tools/conmap"
	"github.com/stianeikeland/go-rpio/v4"
)

const waitTime = 50 * time.Microsecond

// GameBoyRPiPin implements the GameBoyPin interface. This implementation maps connections between
// a RaspberryPi and the GameBoy via GPIO
type GameBoyRPiPin rpio.Pin

// Read returns the RaspberryPi GPIO pin status
func (p GameBoyRPiPin) Read() bool {
	return rpio.ReadPin(rpio.Pin(p)) == rpio.High
}

// High sets the the RaspberryPi GPIO pin status to high
func (p GameBoyRPiPin) High() {
	rpio.WritePin(rpio.Pin(p), rpio.High)
}

// Low sets the the RaspberryPi GPIO pin status to low
func (p GameBoyRPiPin) Low() {
	rpio.WritePin(rpio.Pin(p), rpio.Low)
}

// SetState sets the the RaspberryPi GPIO pin to the given status
func (p GameBoyRPiPin) SetState(state bool) {
	if state {
		rpio.Pin(p).High()
	} else {
		rpio.Pin(p).Low()
	}
}

// Input sets the the RaspberryPi GPIO pin mode to input (populated by the GameBoy)
func (p GameBoyRPiPin) Input() {
	rpio.PinMode(rpio.Pin(p), rpio.Input)
}

// Output sets the the RaspberryPi GPIO pin mode to input (populated by the RPi)
func (p GameBoyRPiPin) Output() {
	rpio.PinMode(rpio.Pin(p), rpio.Output)
}

// RPiGameBoyProxy implements the GameBoyProxy to provide a working data transfer between
// a RaspberryPi and the GameBoy
type RPiGameBoyProxy struct {
	As []GameBoyRPiPin
	Db []GameBoyRPiPin
	Rd GameBoyRPiPin
	Wr GameBoyRPiPin
}

func initGPIO() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewRPiGameBoyProxy creates a new RPiGameBoyProxy. if isMaster is set to
// true then it means that the raspberry is the one in charge of managing the
// cartridge, meaning that it has to send the read/write operations along side
// selecting the address. Contrarily, if isMaster is set to false, the raspberry
// acts as a slave just forwarding to the GameBoy the requested byte
func NewRPiGameBoyProxy(cm *conmap.GameBoyRaspberryMapping, isMaster bool) *RPiGameBoyProxy {
	initGPIO()

	as := []GameBoyRPiPin{GameBoyRPiPin(cm.A0), GameBoyRPiPin(cm.A1), GameBoyRPiPin(cm.A2),
		GameBoyRPiPin(cm.A3), GameBoyRPiPin(cm.A4), GameBoyRPiPin(cm.A5), GameBoyRPiPin(cm.A6),
		GameBoyRPiPin(cm.A7), GameBoyRPiPin(cm.A8), GameBoyRPiPin(cm.A9), GameBoyRPiPin(cm.A10),
		GameBoyRPiPin(cm.A11), GameBoyRPiPin(cm.A12), GameBoyRPiPin(cm.A13), GameBoyRPiPin(cm.A14),
		GameBoyRPiPin(cm.A15),
	}

	db := []GameBoyRPiPin{GameBoyRPiPin(cm.D0), GameBoyRPiPin(cm.D1), GameBoyRPiPin(cm.D2),
		GameBoyRPiPin(cm.D3), GameBoyRPiPin(cm.D4), GameBoyRPiPin(cm.D5), GameBoyRPiPin(cm.D6),
		GameBoyRPiPin(cm.D7),
	}

	// AX pins are the address selector.
	for _, a := range as {
		// If Raspberry manages the cartridge then the address selector must
		// be in output mode
		if isMaster {
			a.Output()
			a.Low()
		} else {
			a.Input()
		}
	}

	for _, d := range db {
		d.Output()
		d.Low()
		time.Sleep(waitTime)
		d.Input()
	}

	rd := GameBoyRPiPin(cm.RD)
	if isMaster {
		rd.Output()
		rd.High()
	} else {
		rd.Input()
	}

	wr := GameBoyRPiPin(cm.WR)
	if isMaster {
		wr.Output()
		wr.High()
	} else {
		wr.Input()
	}

	return &RPiGameBoyProxy{
		As: as,
		Db: db,
		Rd: rd,
		Wr: wr,
	}
}

// End clears the Raspberry pi GPIO
func (rpigb *RPiGameBoyProxy) End() {
	// Unmap gpio memory when done
	rpio.Close()
}

// Read reads the byte located in the address specified with the SelectAddress method.
func (rpigb *RPiGameBoyProxy) Read() uint8 {
	var result uint8

	rpigb.Rd.Low()
	time.Sleep(waitTime)

	result = 0x00
	for i := 0; i < 8; i++ {
		if rpigb.Db[i].Read() {
			result += (1 << i)
		}
	}

	rpigb.Rd.High()
	time.Sleep(waitTime)

	return result
}

// Write writes the provided value to the selected address with the SelectAddress function
func (rpigb *RPiGameBoyProxy) Write(value uint8) {
	// When writing we set DX pins to output mode
	writeToRPiPins(uint(value), rpigb.Db)

	rpigb.Wr.Low()
	time.Sleep(waitTime)

	rpigb.Wr.High()
	time.Sleep(waitTime)

	rpigb.SetReadMode()
}

// SelectAddress sets the GPIO pins status so the referenced address in the cartridge is the given one
func (rpigb *RPiGameBoyProxy) SelectAddress(addr uint) {
	writeToRPiPins(addr, rpigb.As)
}

func (rpigb *RPiGameBoyProxy) SetReadMode() {
	for _, d := range rpigb.Db {
		d.Low()
		d.Input()
	}
	time.Sleep(waitTime)
}

func (rpigb *RPiGameBoyProxy) SetWriteMode() {
	for _, d := range rpigb.Db {
		d.Output()
		d.Low()
	}
	time.Sleep(waitTime)
}

func writeToRPiPins(value uint, pins []GameBoyRPiPin) {
	gbPins := make([]GameBoyPin, 0)
	for _, p := range pins {
		gbPins = append(gbPins, GameBoyPin(p))
	}
	writeToPins(value, gbPins)

	// Wait for GameBoy to do the write
	time.Sleep(waitTime)
}
