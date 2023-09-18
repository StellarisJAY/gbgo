package cartridge

import "fmt"

type header struct {
	title            string // game title in upper case
	manufacturerCode string // manufacturer code in upper case
	cgbOnly          bool
	mapper           byte   // mapper number
	romSize          uint32 // rom size
	ramSize          uint32 // ram size

	entryPoint []byte
}

type BasicCartridge struct {
	h   header
	raw []byte // raw data
}

func MakeBasicCartridge(raw []byte) BasicCartridge {
	return BasicCartridge{
		h:   makeHeader(raw),
		raw: raw,
	}
}

func makeHeader(raw []byte) header {
	title := string(raw[0x134:0x13F])
	code := string(raw[0x13F:0x143])
	var cgbOnly bool
	switch raw[0x143] {
	case 0x80:
		cgbOnly = false
	case 0xC0:
		cgbOnly = true
	}
	// romSize = 32KiB * (1 << value)
	romSize := uint32(1<<raw[0x148]) * 32 * 1024
	var ramSize uint32
	switch raw[0x149] {
	case 0:
		ramSize = 0
	case 1:
		ramSize = 0
	case 2:
		ramSize = 8 * 1024
	case 3:
		ramSize = 32 * 1024
	case 4:
		ramSize = 128 * 1024
	case 5:
		ramSize = 64 * 1024
	}
	return header{
		title:            title,
		manufacturerCode: code,
		cgbOnly:          cgbOnly,
		mapper:           raw[0x147],
		romSize:          romSize,
		ramSize:          ramSize,
		entryPoint:       raw[0x100:0x104],
	}
}

func (bc *BasicCartridge) Info() {
	fmt.Println("title: ", bc.h.title)
	fmt.Println("manufacturer: ", bc.h.manufacturerCode)
	fmt.Printf("mapper num: %d\n", bc.h.mapper)
	fmt.Printf("rom size: %d KiB\n", bc.h.romSize>>10)
	fmt.Printf("ram size: %d KiB\n", bc.h.ramSize>>10)
}
