package bus

import (
	"github.com/StellarisJAY/gbgo/cartridge"
	"github.com/StellarisJAY/gbgo/interrupt"
	"github.com/StellarisJAY/gbgo/ppu"
)

// Bus 虚拟总线，cpu通过总线地址访问内存和硬件
type Bus struct {
	workRAMBanks [][]byte // 可切换的bank 0~7
	wRAMSelect   byte
	highRAM      []byte
	cartridge    *cartridge.BasicCartridge // 卡带数据
	ppu          *ppu.PPU                  // ppu 显卡

	iEnable *interrupt.Register // IE 寄存器
	iFlag   *interrupt.Register // IF 寄存器
}

func MakeBus(cart *cartridge.BasicCartridge) *Bus {
	b := &Bus{
		highRAM:   make([]byte, 256),
		cartridge: cart,
		iEnable:   interrupt.NewRegister(),
		iFlag:     interrupt.NewRegister(),
	}
	b.workRAMBanks = make([][]byte, 8)
	b.workRAMBanks[0] = make([]byte, 0x1000)
	// CGB模式下的1~7号bank
	if b.cartridge.IsCGBMode() {
		for i := 1; i <= 7; i++ {
			b.workRAMBanks[i] = make([]byte, 0x1000)
		}
	}
	return b
}

func (b *Bus) ConnectPPU(ppu *ppu.PPU) {
	b.ppu = ppu
}

func (b *Bus) ReadMem8(addr uint16) byte {
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // cartridge banks
		return b.cartridge.Read(addr)
	case addr >= 0x8000 && addr <= 0x9FFF: // 显存，可切换bank 0/1
		return b.ppu.ReadVRAM(addr)
	case addr >= 0xA000 && addr <= 0xBFFF: // 外部RAM
		return b.cartridge.Read(addr)
	case addr >= 0xC000 && addr <= 0xCFFF: // work ram bank0
		return b.workRAMBanks[0][addr-0xC000]
	case addr >= 0xD000 && addr <= 0xDFFF: // work ram bank 1~7
		return b.workRAMBanks[b.wRAMSelect][addr-0xD000]
	case addr >= 0xFF80 && addr <= 0xFFFE:
		return b.highRAM[addr-0xFF80]
	case addr == 0xFF44: // LY
		return b.ppu.ReadScanline()
	}
	return 0
}

func (b *Bus) WriteMem8(addr uint16, data byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x7FFF: // cartridge banks
		b.cartridge.Write(addr, data)
	case addr >= 0x8000 && addr <= 0x9FFF: // 显存，可切换bank 0/1
		b.ppu.WriteVRAM(addr, data)
	case addr >= 0xA000 && addr <= 0xBFFF: // 外部RAM
		b.cartridge.Write(addr, data)
	case addr >= 0xC000 && addr <= 0xCFFF: // work ram bank0
		b.workRAMBanks[0][addr-0xC000] = data
	case addr >= 0xD000 && addr <= 0xDFFF: // work ram bank 1~7
		b.workRAMBanks[b.wRAMSelect][addr-0xD000] = data
	case addr >= 0xFF80 && addr <= 0xFFFE:
		b.highRAM[addr-0xFF80] = data
	case addr == 0xFF46: // OAM DMA
		dmaAddr := uint16(data) << 8
		b.dmaWriteOAM(dmaAddr)
	case addr == 0xFF4F: // switch vRAM bank0/1
		b.ppu.SwitchVRAMBank(data & 1)
	case addr == 0xFF70: // switch work RAM bank 1~7
		b.switchWorkRAM(data & 7)
	}
}

func (b *Bus) switchWorkRAM(bankSel byte) {
	if b.cartridge.IsCGBMode() {
		if bankSel == 0 {
			b.wRAMSelect = 1
		} else {
			b.wRAMSelect = bankSel
		}
	}
}

func (b *Bus) ReadMem16(addr uint16) uint16 {
	low, high := b.ReadMem8(addr), b.ReadMem8(addr+1)
	return uint16(high)<<8 + uint16(low)
}

// OAM DMA 从addr地址拷贝OAM数据到PPU的OAM内存
func (b *Bus) dmaWriteOAM(addr uint16) {
	start, end := addr, addr+0x9F
	buffer := make([]byte, 160) // OAM内存大小是160字节
	for i := start; i <= end; i++ {
		buffer[i-start] = b.ReadMem8(i)
	}
	b.ppu.WriteOAM(buffer)
}

// RequestInterrupt 设备通过该方法向总线发起中断
func (b *Bus) RequestInterrupt(code interrupt.Code) {
	b.iFlag.Get(code)
}

func (b *Bus) disableAllInterrupts() {
	b.iEnable.Clear()
}
