package cartridge

import "fmt"

type MBC1 struct {
	raw               []byte
	romBank0          []byte
	switchableRomBank []byte
	romBankSelect     byte

	ramEnabled    bool
	ramBanks      [][]byte
	ramBankSelect byte
}

func makeMBC1(raw []byte, ramSize uint32) *MBC1 {
	mbc1 := &MBC1{
		raw:               raw,
		romBank0:          raw[0:0x4000],
		switchableRomBank: raw[0x4000:],
		romBankSelect:     1,
		ramEnabled:        false,
		ramBankSelect:     0,
	}
	if ramSize != 0 {
		mbc1.ramEnabled = true
		mbc1.ramBanks = make([][]byte, 4)
		for i := 0; i < 4; i++ {
			mbc1.ramBanks[i] = make([]byte, 0x2000)
		}
	}
	return mbc1
}

func (m *MBC1) Read(addr uint16) byte {
	switch {
	case addr <= 0x3FFF:
		return m.romBank0[addr]
	case addr >= 0x4000 && addr <= 0x7FFF: // ROM Bank 01~7F
		return m.switchableRomBank[addr-0x4000]
	case addr >= 0xA000 && addr <= 0xBFFF: // RAM Bank 0~3
		if m.ramEnabled {
			return m.ramBanks[m.ramBankSelect][addr-0xA000]
		}
	}
	panic(fmt.Errorf("invalid cartridge address: 0x%X", addr))
}

func (m *MBC1) Write(addr uint16, data byte) {
	switch {
	case addr <= 0x1FFF: // RAM Enable register
		m.ramEnabled = data&0xF == 0xA
	case addr >= 0x2000 && addr <= 0x3FFF: // ROM Bank Number
		m.switchRomBank(data)
	case addr >= 0x4000 && addr <= 0x5FFF: // RAM Bank Number

	case addr >= 0x6000 && addr <= 0x7FFF: // Banking Mode select
	}
}

func (m *MBC1) switchRomBank(bankSel byte) {
	m.romBankSelect = bankSel & 0x1F
	start := uint32(m.romBankSelect) * 0x4000
	m.switchableRomBank = m.raw[start : start+0x4000]
}
