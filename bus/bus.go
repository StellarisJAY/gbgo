package bus

type Bus struct {
	workRAM0 []byte
	workRAM1 []byte
}

func (b *Bus) ReadMem8(addr uint16) byte {
	// todo read ROM
	switch {
	case addr >= 0xC000 && addr <= 0xCFFF:
		return b.workRAM0[addr-0xC000]
	case addr >= 0xD000 && addr <= 0xDFFF:
		return b.workRAM1[addr-0xD000]
	}
	return 0
}

func (b *Bus) WriteMem8(addr uint16, data byte) {
	// todo bus write
}

func (b *Bus) ReadMem16(addr uint16) uint16 {
	return 0
}
