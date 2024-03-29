package cartridge

import (
	"errors"
	"fmt"
)

const (
	// Cartridge type
	RomOnly                    uint8 = 0x00
	MBC1                       uint8 = 0x01
	MBC1RAM                    uint8 = 0x02
	MBC1RAMBattery             uint8 = 0x03
	MBC2                       uint8 = 0x05
	MBC2Battery                uint8 = 0x06
	ROMRAM                     uint8 = 0x08
	ROMRAMBattery              uint8 = 0x09
	MMM01                      uint8 = 0x0B
	MMM01RAM                   uint8 = 0x0C
	MMM01RAMBattery            uint8 = 0x0D
	MBC3TimerBattery           uint8 = 0x0F
	MBC3TimerRAMBattery        uint8 = 0x10
	MBC3                       uint8 = 0x11
	MBC3RAM                    uint8 = 0x12
	MBC3RAMBattery             uint8 = 0x13
	MBC5                       uint8 = 0x19
	MBC5RAM                    uint8 = 0x1A
	MBC5RAMBattery             uint8 = 0x1B
	MBC5Rumble                 uint8 = 0x1C
	MBC5RumbleRAM              uint8 = 0x1D
	MBC5RumbleRAMBattery       uint8 = 0x1E
	MBC6                       uint8 = 0x20
	MBC7SensorRumbleRAMBattery uint8 = 0x22
	PocketCamera               uint8 = 0xFC
	BandaiTAMA5                uint8 = 0xFD
	HuC3                       uint8 = 0xFE
	HuC1RAMBattery             uint8 = 0xFF
)

// ROM Sizes
const (
	ROM32KB  uint8 = 0x00
	ROM64KB  uint8 = 0x01
	ROM128KB uint8 = 0x02
	ROM256KB uint8 = 0x03
	ROM512KB uint8 = 0x04
	ROM1MB   uint8 = 0x05
	ROM2MB   uint8 = 0x06
	ROM4MB   uint8 = 0x07
	ROM8MB   uint8 = 0x08
)

// RAM Sizes
const (
	None     uint8 = 0x00
	Unused   uint8 = 0x01
	RAM8KB   uint8 = 0x02
	RAM32KB  uint8 = 0x03
	RAM128KB uint8 = 0x04
	RAM64KB  uint8 = 0x05
)

// CartridgeHeader contains all the information stored in the GB cartridge header
type CartridgeHeader struct {
	rawBytes         []uint8
	NintendoLogo     []uint8 // 0104-0133 - Nintendo Logo
	Title            []uint8 // 0134-0143 - Title
	ManufacturerCode []uint8 // 013F-0142 - Manufacturer Code
	CGBFlag          uint8   // 0143 - CGB Flag
	LicenseeCode     []uint8 // 0144-0145 - New Licensee Code
	SGBFlag          uint8   // 0146 - SGB Flag
	CartridgeType    uint8   // 0147 - Cartridge Type
	ROMSize          uint8   // 0148 - ROM Size
	RAMSize          uint8   // 0149 - RAM Size
	DestinationCode  uint8   // 014A - Destination Code
	OldLicenseeCode  uint8   // 014B - Old Licensee Code
	MaskROMVersion   uint8   // 014C - Mask ROM Version number
	HeaderChecksum   uint8   // 014D - Header Checksum
	GlobalChecksum   []uint8
}

// ROMHeaderFromBytes loads the given bytes into the CartridgeHeader structure and returns a
// reference to the recently created structure
func ROMHeaderFromBytes(bytes []uint8) *CartridgeHeader {
	return &CartridgeHeader{
		rawBytes:         bytes,
		NintendoLogo:     bytes[0x104:0x134],
		Title:            bytes[0x134:0x144],
		ManufacturerCode: bytes[0x13F:0x143],
		CGBFlag:          bytes[0x143],
		LicenseeCode:     bytes[0x144:0x145],
		SGBFlag:          bytes[0x146],
		CartridgeType:    bytes[0x147],
		ROMSize:          bytes[0x148],
		RAMSize:          bytes[0x149],
		DestinationCode:  bytes[0x14A],
		OldLicenseeCode:  bytes[0x14B],
		MaskROMVersion:   bytes[0x14C],
		HeaderChecksum:   bytes[0x14D],
		GlobalChecksum:   bytes[0x14E:0x150],
	}
}

func (ch *CartridgeHeader) PrintInfo() {
	fmt.Println("*** Cartridge Header ***")
	fmt.Println("Title:", string(ch.Title))
	fmt.Printf("Header Checksum: 0x%02x\n", ch.HeaderChecksum)
	fmt.Printf("Global Checksum: 0x%02x 0x%02x\n", ch.GlobalChecksum[0], ch.GlobalChecksum[1])
	fmt.Println("Cartridge Type:", ch.CartridgeTypeText())
	fmt.Print("Supports SGB: ")
	if ch.SupportsSGB() {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
	fmt.Println("# ROM Banks:", ch.GetNumROMBanks())
	fmt.Println("# RAM Banks:", ch.GetNumRAMBanks())
}

// IsGBCOnly returns true if the cartridge can only run in a GameBoy color
func (ch *CartridgeHeader) IsGBCOnly() bool {
	return ch.CGBFlag == 0xC0
}

// SupportsSGB returns true if the cartridge supports Super GameBoy
func (ch *CartridgeHeader) SupportsSGB() bool {
	return ch.SGBFlag == 0x03
}

// HasMBC returns true if the cartridge has a MBC
func (ch *CartridgeHeader) HasMBC() bool {
	return ch.CartridgeType != RomOnly && ch.CartridgeType != ROMRAM && ch.CartridgeType != ROMRAMBattery
}

// IsMBC1 returns true if the cartridge is MBC1 type
func (ch *CartridgeHeader) IsMBC1() bool {
	return ch.CartridgeType == MBC1 || ch.CartridgeType == MBC1RAM || ch.CartridgeType == MBC1RAMBattery
}

// IsMBC2 returns true if the cartridge is MBC2 type
func (ch *CartridgeHeader) IsMBC2() bool {
	return ch.CartridgeType == MBC2 || ch.CartridgeType == MBC2Battery
}

// IsMBC3 returns true if the cartridge is MBC3 type
func (ch *CartridgeHeader) IsMBC3() bool {
	return ch.CartridgeType == MBC3 || ch.CartridgeType == MBC3RAM || ch.CartridgeType == MBC3RAMBattery ||
		ch.CartridgeType == MBC3TimerBattery || ch.CartridgeType == MBC3TimerRAMBattery
}

// IsMBC5 returns true if the cartridge is MBC5 type
func (ch *CartridgeHeader) IsMBC5() bool {
	return ch.CartridgeType == MBC5 || ch.CartridgeType == MBC5RAM || ch.CartridgeType == MBC5RAMBattery ||
		ch.CartridgeType == MBC5Rumble || ch.CartridgeType == MBC5RumbleRAM ||
		ch.CartridgeType == MBC5RumbleRAMBattery
}

func (ch *CartridgeHeader) HasBattery() bool {
	return ch.CartridgeType == MBC1RAMBattery || ch.CartridgeType == MBC2Battery || ch.CartridgeType == ROMRAMBattery ||
		ch.CartridgeType == MMM01RAMBattery || ch.CartridgeType == MBC3TimerBattery || ch.CartridgeType == MBC3TimerRAMBattery ||
		ch.CartridgeType == MBC3RAMBattery || ch.CartridgeType == MBC5RAMBattery || ch.CartridgeType == MBC5RumbleRAMBattery ||
		ch.CartridgeType == MBC7SensorRumbleRAMBattery || ch.CartridgeType == HuC1RAMBattery
}

func (ch *CartridgeHeader) HasRAM() bool {
	return ch.CartridgeType == MBC1RAMBattery || ch.CartridgeType == MBC1RAM || ch.CartridgeType == ROMRAM ||
		ch.CartridgeType == ROMRAMBattery || ch.CartridgeType == MMM01RAM || ch.CartridgeType == MMM01RAMBattery ||
		ch.CartridgeType == MBC3RAMBattery || ch.CartridgeType == MBC3RAM || ch.CartridgeType == MBC5RumbleRAMBattery ||
		ch.CartridgeType == MBC7SensorRumbleRAMBattery || ch.CartridgeType == HuC1RAMBattery ||
		ch.CartridgeType == MBC5RAMBattery || ch.CartridgeType == MBC5RumbleRAM
}

func (ch *CartridgeHeader) CartridgeTypeText() string {
	var msg string
	if !ch.HasMBC() {
		msg = "ROM (no MBC)"
	} else if ch.IsMBC1() {
		msg = "MBC1"
	} else if ch.IsMBC2() {
		msg = "MBC2"
	} else if ch.IsMBC3() {
		msg = "MBC3"
	} else if ch.IsMBC5() {
		msg = "MBC5"
	}

	if ch.HasRAM() {
		msg = msg + " + RAM"
	}

	if ch.HasBattery() {
		msg = msg + " + Battery"
	}

	if ch.IsGBCOnly() {
		msg = msg + " (GBC Only)"
	}

	return msg
}

// GetNumROMBanks returns the number of ROM banks in the cartridge
func (ch *CartridgeHeader) GetNumROMBanks() int {
	return map[uint8]int{
		ROM32KB:  2,
		ROM64KB:  4,
		ROM128KB: 8,
		ROM256KB: 16,
		ROM512KB: 32,
		ROM1MB:   64,
		ROM2MB:   128,
		ROM4MB:   256,
		ROM8MB:   512,
	}[ch.ROMSize]
}

// GetNumRAMBanks returns the number of RAM banks in the cartridge
func (ch *CartridgeHeader) GetNumRAMBanks() int {
	return map[uint8]int{
		None:     0,
		Unused:   0,
		RAM8KB:   1,
		RAM32KB:  4,
		RAM64KB:  8,
		RAM128KB: 16,
	}[ch.RAMSize]
}

// Validate runs the checksum procedure and compares te result agains the byte located at
// 0x14D. If the result matches with the predefined checksum means that the dump has been successful
func (ch *CartridgeHeader) Validate() error {
	var x uint
	x = 0x00

	for i := 0x0134; i < 0x014D; i++ {
		x = x - uint(ch.rawBytes[i]) - 1
		x = x & 0xFF
	}

	valid := uint8(x) == ch.HeaderChecksum
	if !valid {
		return errors.New("invalid header checksum")
	} else {
		return nil
	}
}
