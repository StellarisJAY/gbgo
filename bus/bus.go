package bus

import "github.com/StellarisJAY/gbgo/cartridge"

// Bus 虚拟总线，cpu通过总线地址访问内存和硬件
type Bus struct {
	workRAM0  []byte // 固定的bank 0
	workRAM1  []byte // cgbMode，可切换的workRam的bank 1~7
	highRAM   []byte
	cartridge *cartridge.BasicCartridge // 卡带数据
}

func MakeBus(cart *cartridge.BasicCartridge) *Bus {
	return &Bus{
		workRAM0:  make([]byte, 0x1000),
		workRAM1:  make([]byte, 0x1000),
		highRAM:   make([]byte, 0xFFFF-0xFF80),
		cartridge: cart,
	}
}

func (b *Bus) ReadMem8(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // cartridge banks
		return b.cartridge.Read(addr)
	case addr >= 0x8000 && addr <= 0x9FFF: // 显存，可切换bank 0/1

	case addr >= 0xA000 && addr <= 0xBFFF: // 外部RAM
		return b.cartridge.Read(addr)
	case addr >= 0xC000 && addr <= 0xCFFF: // work ram bank0
		return b.workRAM0[addr-0xC000]
	case addr >= 0xD000 && addr <= 0xDFFF: // work ram bank 1~7
		return b.workRAM1[addr-0xD000]
	case addr >= 0xFF80 && addr <= 0xFFFE:
		return b.highRAM[addr-0xFF80]
	}
	return 0
}

func (b *Bus) WriteMem8(addr uint16, data byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // cartridge banks
		b.cartridge.Write(addr, data)
	case addr >= 0x8000 && addr <= 0x9FFF: // 显存，可切换bank 0/1

	case addr >= 0xA000 && addr <= 0xBFFF: // 外部RAM
		b.cartridge.Write(addr, data)
	case addr >= 0xC000 && addr <= 0xCFFF: // work ram bank0
		b.workRAM0[addr-0xC000] = data
	case addr >= 0xD000 && addr <= 0xDFFF: // work ram bank 1~7
		b.workRAM1[addr-0xD000] = data
	case addr >= 0xFF80 && addr <= 0xFFFE:
		b.highRAM[addr-0xFF80] = data
	}
}

func (b *Bus) ReadMem16(addr uint16) uint16 {
	low, high := b.ReadMem8(addr), b.ReadMem8(addr+1)
	return uint16(high)<<8 + uint16(low)
}
