package cartridge

import "fmt"

type header struct {
	title            string // 游戏名称，大写ASCII
	manufacturerCode string // 生产商编号，大写ASCII
	cgbOnly          bool
	mbc              byte   // MBC
	romSize          uint32 // rom 字节大小
	ramSize          uint32 // ram 字节大小

	entryPoint []byte
}

type BasicCartridge struct {
	h   header
	raw []byte
	mbc MBC // MBC接口
}

// MBC 接口，所有MBC必须实现读写地址
type MBC interface {
	Read(addr uint16) byte
	Write(addr uint16, data byte)
}

func MakeBasicCartridge(raw []byte) BasicCartridge {
	header := makeHeader(raw)
	var mbc MBC
	switch header.mbc {
	case 0:
		mbc = makeNoMBC(raw, header.ramSize)
	default:
		panic(fmt.Errorf("unsupported mbc %d", mbc))
	}
	return BasicCartridge{
		h:   header,
		raw: raw,
		mbc: mbc,
	}
}

func (bc *BasicCartridge) Read(addr uint16) byte {
	return bc.mbc.Read(addr)
}

func (bc *BasicCartridge) Write(addr uint16, data byte) {
	bc.mbc.Write(addr, data)
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
		mbc:              raw[0x147],
		romSize:          romSize,
		ramSize:          ramSize,
		entryPoint:       raw[0x100:0x104],
	}
}

func (bc *BasicCartridge) Info() {
	fmt.Println("title: ", bc.h.title)
	fmt.Println("manufacturer: ", bc.h.manufacturerCode)
	fmt.Printf("mbc num: %d\n", bc.h.mbc)
	fmt.Println("cgb mode: ", bc.h.cgbOnly)
	fmt.Printf("rom size: %d KiB\n", bc.h.romSize>>10)
	fmt.Printf("ram size: %d KiB\n", bc.h.ramSize>>10)
}
