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

// NewCartridge creates a pointer to a Cartridge struct
func NewCartridge(h *CartridgeHeader, rbs [][]uint8) *Cartridge {
	return &Cartridge{
		Header:   h,
		ROMBanks: rbs,
		RAMBanks: nil,
	}
}

// Validate validates the global checksum of the cartridge. Global checksum is a 16 bin number located within
// the cartridge header at the address range 014E-014F. Actually, Game Boy does not validate this
// checksum, therefore there might be games outthere that contain an invalid checksum.
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

	b0 := uint8(result & 0xFF)
	b1 := uint8((result & 0xFF00) >> 8)

	if b0 != c.Header.GlobalChecksum[1] || b1 != c.Header.GlobalChecksum[0] {
		errMsg := fmt.Sprintf("invalid global checksum. Expected 0x%02x 0x%02x found 0x%02x 0x%02x\n",
			c.Header.GlobalChecksum[0], c.Header.GlobalChecksum[1], b1, b0)
		return errors.New(errMsg)
	}

	return nil
}

// Save serializes the cartridge in a binary file storing all the ROM banks sequentially
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
