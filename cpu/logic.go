package cpu

// AND A, N; N = r8, n8, (HL)
func andWithA(p *Processor, op *Instruction) {
	var param byte
	switch op.code {
	case 0xA0: // AND A, B
		param = p.b
	case 0xA1:
		param = p.c
	case 0xA2:
		param = p.d
	case 0xA3:
		param = p.e
	case 0xA4:
		param = p.h
	case 0xA5:
		param = p.l
	case 0xA6:
		param = p.readMem8(p.reg16(p.h, p.l))
	case 0xA7:
		param = p.a
	case 0xE6:
		param = p.readOperand8(p.pc, immediate)
	}
	p.a &= param
	// flag z010
	p.f = 0b0010_0000
	p.setFlag(zeroFlag, p.a == 0)
}

// OR A, N; N = r8, n8, (HL)
func orWithA(p *Processor, op *Instruction) {
	var param byte
	switch op.code {
	case 0xB0: // AND A, B
		param = p.b
	case 0xB1:
		param = p.c
	case 0xB2:
		param = p.d
	case 0xB3:
		param = p.e
	case 0xB4:
		param = p.h
	case 0xB5:
		param = p.l
	case 0xB6:
		param = p.readMem8(p.reg16(p.h, p.l))
	case 0xB7:
		param = p.a
	case 0xF6:
		param = p.readOperand8(p.pc, immediate)
	}
	p.a |= param
	// flag z000
	p.f = 0
	p.setFlag(zeroFlag, p.a == 0)
}

// XOR A, N; N = r8, n8, (HL)
func xorWithA(p *Processor, op *Instruction) {
	var param byte
	switch op.code {
	case 0xA8:
		param = p.b
	case 0xA9:
		param = p.c
	case 0xAA:
		param = p.d
	case 0xAB:
		param = p.d
	case 0xAC:
		param = p.h
	case 0xAD:
		param = p.l
	case 0xAE:
		param = p.readMem8(p.reg16(p.h, p.l))
	case 0xAF:
		param = p.a
	case 0xEE:
		param = p.readOperand8(p.pc, immediate)
	}
	p.a ^= param
	// flag z000
	p.f = 0
	p.setFlag(zeroFlag, p.a == 0)
}
