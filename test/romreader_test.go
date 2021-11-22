package test

import (
	"testing"

	"github.com/Guillem96/gameboy-tools/cartridge"
)

const romFile = "../roms/pkmn_red.gb"

var expectedNintendoLogo = [...]uint8{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00,
	0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC,
	0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC,
	0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

func TestReadHeaderFromFile(t *testing.T) {
	frr := cartridge.NewFileROMReader(romFile)
	header, err := frr.ReadHeader()
	if err != nil {
		t.Error(err)
	}

	err = header.ValidateHeader()
	if err != nil {
		t.Error(err)
	}
}

func TestReadWholeCartridge(t *testing.T) {
	frr := cartridge.NewFileROMReader(romFile)
	cart, err := frr.ReadCartridge()
	if err != nil {
		t.Error(err)
	}
	err = cart.Save("../roms/test.gb")
	if err != nil {
		t.Error(err)
	}
}

func TestReadHeaderAndCheckNintendoLogo(t *testing.T) {
	frr := cartridge.NewFileROMReader(romFile)
	header, err := frr.ReadHeader()
	if err != nil {
		t.Error(err)
	}

	if len(header.NintendoLogo) != len(expectedNintendoLogo) {
		t.Errorf("Parsed logo and expected does not match. Length %d != %d", len(header.NintendoLogo), len(expectedNintendoLogo))
	}
	for i, eb := range expectedNintendoLogo {
		if eb != header.NintendoLogo[i] {
			t.Errorf("Parsed logo byte %d and expected does not match.", i)
		}
	}
}
