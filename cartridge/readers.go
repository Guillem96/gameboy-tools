package cartridge

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// FileROMReader is the object responsible of reading and parsing a Game Boy
// local file ROM.
type FileROMReader struct {
	l         *log.Logger
	fname     string
	inmemfile []byte
	header    *CartridgeHeader
}

// FileROMReader creates a new rom file reader and returns a pointer to it
func NewFileROMReader(fname string) *FileROMReader {
	return &FileROMReader{
		header:    nil,
		fname:     fname,
		inmemfile: nil,
		l:         log.New(os.Stdout, "[GB ROM File Reader]", log.LstdFlags),
	}
}

// ReadHeader reads the whole cartridge header
func (frr *FileROMReader) ReadHeader() (*CartridgeHeader, error) {
	err := frr.loadROMInMemory()
	if err != nil {
		rerr := fmt.Errorf("reading cartridge header: %v", err)
		return nil, rerr
	}

	if frr.header != nil {
		return frr.header, nil
	}

	frr.l.Println("Reading ROM header data.")
	bytes := frr.inmemfile[:0x150]
	frr.header = ROMHeaderFromBytes(bytes)
	return frr.header, nil
}

// ReadCartridge dumps the whole cartridge data. Reads all ROM & RAM banks
func (frr *FileROMReader) ReadCartridge() (*Cartridge, error) {
	h, err := frr.ReadHeader()
	if err != nil {
		return nil, err
	}

	// Dump Rom banks
	nb := h.GetNumROMBanks()
	banks := make([][]uint8, nb)
	frr.l.Printf("The cartridge has %d banks.\n", nb)

	romSize := 0x4000
	for b := 0; b < nb; b++ {
		start := b * romSize
		end := start + romSize
		banks[b] = frr.inmemfile[start:end]
	}

	// TODO: RAM is in a separate file in this case

	return NewCartridge(h, banks), nil
}

func (frr *FileROMReader) loadROMInMemory() error {
	if frr.inmemfile == nil {
		rb, err := loadFile(frr.fname)
		frr.inmemfile = rb
		if err != nil {
			return err
		}
	}
	return nil
}

func loadFile(fname string) ([]byte, error) {
	file, err := os.Open(fname)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	size := stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
