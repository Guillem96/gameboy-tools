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
