package cpu

// ld r8, n8
func loadImmediate8(p *Processor, op *instruction) {
	data := p.readOperand8(p.pc, op.mode)
	switch op.code {
	case 0x06:
		p.b = data
	case 0x0E:
		p.c = data
	case 0x16:
		p.d = data
	case 0x1E:
		p.e = data
	case 0x26:
		p.h = data
	case 0x2E:
		p.l = data
	default:
		panic("impossible opcode")
	}
}

// ld r16, n16
func loadImmediate16(p *Processor, op *instruction) {
	data := p.readOperand16(p.pc, immediate)
	switch op.code {
	case 0x01: // ld BC, n16
		p.writeBC(data)
	case 0x11: // ld DE, n16
		p.writeDE(data)
	case 0x21: // ld HL, n16
		p.writeHL(data)
	case 0x31: // ld SP, n16
		p.sp = data
	}
}

// ld A, r8
func loadAFromReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x7F: // ld A, A
	case 0x78:
		p.a = p.b
	case 0x79:
		p.a = p.c
	case 0x7A:
		p.a = p.d
	case 0x7B:
		p.a = p.e
	case 0x7C:
		p.a = p.h
	case 0x7D:
		p.a = p.l
	}
}

// ld B, r8
func loadBFromReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x40: // ld B, B
	case 0x41:
		p.b = p.c
	case 0x42:
		p.b = p.d
	case 0x43:
		p.b = p.e
	case 0x44:
		p.b = p.h
	case 0x45:
		p.b = p.l
	}
}

// ld C, r8
func loadCFromReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x48:
		p.c = p.b
	case 0x49: // ld C,C
	case 0x4A:
		p.c = p.d
	case 0x4B:
		p.c = p.e
	case 0x4C:
		p.c = p.h
	case 0x4D:
		p.c = p.l
	}
}

// ld D, r8
func loadDFromReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x50:
		p.d = p.b
	case 0x51:
		p.d = p.c
	case 0x52: // ld D, D
	case 0x53:
		p.d = p.e
	case 0x54:
		p.d = p.h
	case 0x55:
		p.d = p.l
	}
}

// ld E, r8
func loadEFromReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x58:
		p.e = p.b
	case 0x59:
		p.e = p.c
	case 0x5A:
		p.e = p.d
	case 0x5B: // ld E, E
	case 0x5C:
		p.e = p.h
	case 0x5D:
		p.e = p.l
	}
}

// ld H, r8
func loadHFromReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x60:
		p.h = p.b
	case 0x61:
		p.h = p.c
	case 0x62:
		p.h = p.d
	case 0x63:
		p.h = p.e
	case 0x64: // ld H, H
	case 0x65:
		p.h = p.l
	}
}

// ld L, r8
func loadLFromReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x68:
		p.l = p.b
	case 0x69:
		p.l = p.c
	case 0x6A:
		p.l = p.d
	case 0x6B:
		p.l = p.e
	case 0x6C:
		p.l = p.h
	case 0x6D: // ld L, L
	}
}

// ld r8, (HL)
func loadRegFromHL(p *Processor, op *instruction) {
	data := p.readMem8(p.reg16(p.h, p.l))
	switch op.code {
	case 0x7E:
		p.a = data
	case 0x46:
		p.b = data
	case 0x4E:
		p.c = data
	case 0x56:
		p.d = data
	case 0x5E:
		p.e = data
	case 0x66:
		p.h = data
	case 0x6E:
		p.l = data
	}
}

// ld (HL), r8
func storeRegInHL(p *Processor, op *instruction) {
	addr := p.reg16(p.h, p.l)
	switch op.code {
	case 0x70:
		p.writeMem8(addr, p.b)
	case 0x71:
		p.writeMem8(addr, p.c)
	case 0x72:
		p.writeMem8(addr, p.d)
	case 0x73:
		p.writeMem8(addr, p.e)
	case 0x74:
		p.writeMem8(addr, p.h)
	case 0x75:
		p.writeMem8(addr, p.l)
	}
}

// ld A, (N)
func loadA(p *Processor, op *instruction) {
	switch op.code {
	case 0x0A: // ld A, (BC)
		p.a = p.readMem8(p.reg16(p.b, p.c))
	case 0x1A: // ld A, (DE)
		p.a = p.readMem8(p.reg16(p.d, p.e))
	case 0x7E: // ld A, (HL)
		p.a = p.readMem8(p.reg16(p.h, p.l))
	case 0xFA: // ld A, (nn)
		p.a = p.readOperand8(p.pc, absolute)
	}
}

// ld (N), A
func storeA(p *Processor, op *instruction) {
	switch op.code {
	case 0x02: // ld (BC), A
		p.writeMem8(p.reg16(p.b, p.c), p.a)
	case 0x12: // ld (DE), A
		p.writeMem8(p.reg16(p.d, p.e), p.a)
	case 0x77: // ld (HL), A
		p.writeMem8(p.reg16(p.h, p.l), p.a)
	case 0xEA: // ld (nn), A
		addr := p.readOperand16(p.pc, immediate)
		p.writeMem8(addr, p.a)
	}
}

// ld r8, A
func storeAInReg(p *Processor, op *instruction) {
	switch op.code {
	case 0x47:
		p.b = p.a
	case 0x4F:
		p.c = p.a
	case 0x57:
		p.d = p.a
	case 0x5F:
		p.e = p.a
	case 0x67:
		p.h = p.a
	case 0x6F:
		p.l = p.a
	}
}

// ld (HL), n8
func storeImmediateInHL(p *Processor, op *instruction) {
	addr := p.reg16(p.h, p.l)
	data := p.readOperand8(p.pc, op.mode)
	p.writeMem8(addr, data)
}

// ld A, (0xFF00 + N); N = C or n8
func readIOPort(p *Processor, op *instruction) {
	var addr uint16 = 0xFF00
	switch op.code {
	case 0xF0: // ld A, (0xFF00 + n8)
		addr += uint16(p.readOperand8(p.pc, op.mode))
	case 0xF2: // ld A, (0XFF00 + C)
		addr += uint16(p.c)
	}
	p.a = p.readMem8(addr)
}

// ld (0XFF00 + N), A; N = C or n8
func writeIOPort(p *Processor, op *instruction) {
	var addr uint16 = 0xFF00
	switch op.code {
	case 0xE0: // ld (0xFF00 + n8), A
		addr += uint16(p.readOperand8(p.pc, op.mode))
	case 0xE2: // ld (0xFF00 + C), A
		addr += uint16(p.c)
	}
	p.writeMem8(addr, p.a)
}

// A = (HL), HL += 1
// (HL) = A, HL += 1
func ldi(p *Processor, op *instruction) {
	addr := p.reg16(p.h, p.l)
	switch op.code {
	case 0x22: // ldi (HL), A
		p.writeMem8(addr, p.a)
	case 0x2A: // ldi A, (HL)
		p.a = p.readMem8(addr)
	default:
		panic("impossible ldi opcode")
	}
	addr += 1
	p.writeHL(addr)
}

// A = (HL), HL -= 1
// (HL) = A, HL -= 1
func ldd(p *Processor, op *instruction) {
	addr := p.reg16(p.h, p.l)
	switch op.code {
	case 0x32: // ldd (HL), A
		p.writeMem8(addr, p.a)
	case 0x3A: // ldd A, (HL)
		p.a = p.readMem8(addr)
	}
	addr -= 1
	p.writeHL(addr)
}
