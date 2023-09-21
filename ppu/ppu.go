package ppu

import (
	"github.com/StellarisJAY/gbgo/interrupt"
)

type PPU struct {
	lcdc       LCDControl
	scanline   byte
	oam        []byte
	vRAMBanks  [][]byte // 两个8KiB的VRAM bank
	vRAMSelect byte     // CGB mode可切换bank

	interruptRequester interrupt.Requester
}

func MakePPU(requester interrupt.Requester) *PPU {
	return &PPU{
		vRAMBanks: [][]byte{
			make([]byte, 0x2000),
			make([]byte, 0x2000),
		},
		vRAMSelect: 0,

		interruptRequester: requester,
	}
}

func (p *PPU) Render() {
	renderStartingPage()
}

func (p *PPU) ReadScanline() byte {
	return p.scanline
}

func (p *PPU) WriteOAM(data []byte) {
	p.oam = data
}

func (p *PPU) ReadVRAM(addr uint16) byte {
	return p.vRAMBanks[p.vRAMSelect][addr-0x8000]
}

func (p *PPU) WriteVRAM(addr uint16, data byte) {
	p.vRAMBanks[p.vRAMSelect][addr-0x8000] = data
}

func (p *PPU) SwitchVRAMBank(bank byte) {
	p.vRAMSelect = bank
}
