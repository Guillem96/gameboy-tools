package cartridge

import (
	"errors"
	"fmt"
	"os"
)

// Reference: https://gbdev.io/pandocs/The_Cartridge_Header.html

type Cartridge struct {
	Header   *CartridgeHeader
	ROMBanks [][]uint8
	RAMBanks [][]uint
}

func NewCartridge(h *CartridgeHeader, rbs [][]uint8) *Cartridge {
	return &Cartridge{
		Header:   h,
		ROMBanks: rbs,
		RAMBanks: nil,
	}
}

func (c *Cartridge) Validate() error {
	nb := 0
	var result uint16
	result = 0
	for _, bank := range c.ROMBanks {
		for _, b := range bank {
			if nb != 0x14E && nb != 0x14F {
				result += uint16(b)
			}
			nb += 1
		}
	}
	b0 := result & 0xFF
	b1 := result & 0xFF00
	fmt.Printf("-- Global Checksum ----------------------\n %02x %02x\n", b0, b1)
	return nil
}

func (c *Cartridge) Save(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		errMsg := fmt.Sprintf("error creating file %v: %+v\n", fname, err)
		return errors.New(errMsg)
	}
	defer f.Close()

	for i, b := range c.ROMBanks {
		_, err := f.Write(b)
		if err != nil {
			errMsg := fmt.Sprintf("error writing bank %d: %+v\n", i, err)
			return errors.New(errMsg)
		}
		err = f.Sync()
		if err != nil {
			errMsg := fmt.Sprintf("failed syncing the bank %d bytes: %+v\n", i, err)
			return errors.New(errMsg)
		}
	}

	return nil
}
