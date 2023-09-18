package cartridge

// NoMBC 没有MBC的ROM ONLY卡带
type NoMBC struct {
	rom      []byte
	ram      []byte
	usingRam bool
}

func (n *NoMBC) Read(addr uint16) byte {
	switch {
	case addr <= 0x7FFF:
		return n.rom[addr]
	case addr >= 0xA000 && addr <= 0xBFFF:
		if n.usingRam {
			return n.ram[addr-0xA000]
		}
	}
	return 0
}

func (n *NoMBC) Write(addr uint16, data byte) {
	switch {
	case addr <= 0x7FFF:
	case addr >= 0xA000 && addr <= 0xBFFF:
		if n.usingRam {
			n.ram[addr-0xA000] = data
		}
	}
}
