package cpu

// inc N; N=r8,(HL)
// N = N + 1, z0h-
func inc(p *Processor, op *instruction) {
	var result, original byte
	switch op.code {
	case 0x04:
		original = p.b
		p.b += 1
		result = p.b
	case 0x0C:
		original = p.c
		p.c += 1
		result = p.c
	case 0x14:
		original = p.d
		p.d += 1
		result = p.d
	case 0x1C:
		original = p.e
		p.e += 1
		result = p.e
	case 0x24:
		original = p.h
		p.h += 1
		result = p.h
	case 0x2C:
		original = p.l
		p.l += 1
		result = p.l
	case 0x34:
		result = p.readMem8(p.reg16(p.h, p.l))
		original = result
		p.writeMem8(p.reg16(p.h, p.l), result)
	case 0x3C:
		original = p.a
		p.a += 1
		result = p.a
	}
	p.setFlag(subFlag, false)
	p.setFlag(zeroFlag, result == 0)
	p.determineHalfCarry(original, result)
}

// dec N; N = r8, (HL)
// N = N-1, z1h-
func dec(p *Processor, op *instruction) {
	var result, original byte
	switch op.code {
	case 0x05:
		original = p.b
		p.b -= 1
		result = p.b
	case 0x0D:
		original = p.c
		p.c -= 1
		result = p.c
	case 0x15:
		original = p.d
		p.d -= 1
		result = p.d
	case 0x1D:
		original = p.e
		p.e -= 1
		result = p.e
	case 0x25:
		original = p.h
		p.h -= 1
		result = p.h
	case 0x2D:
		original = p.l
		p.l -= 1
		result = p.l
	case 0x35:
		result = p.readMem8(p.reg16(p.h, p.l))
		original = result
		p.writeMem8(p.reg16(p.h, p.l), result)
	case 0x3D:
		original = p.a
		p.a -= 1
		result = p.a
	}
	p.setFlag(subFlag, true)
	p.setFlag(zeroFlag, result == 0)
	p.determineHalfCarry(original, 1)
}

// ADD A, N; N = r8, n8, (HL)
// A = A + N; z0hc
func addA(p *Processor, op *instruction) {
	original := p.a
	var delta byte
	switch op.code {
	case 0x80:
		delta = p.b
	case 0x81:
		delta = p.c
	case 0x82:
		delta = p.d
	case 0x83:
		delta = p.e
	case 0x84:
		delta = p.h
	case 0x85:
		delta = p.l
	case 0x86: // ADD A, (HL)
		delta = p.readMem8(p.reg16(p.h, p.l))
	case 0x87: // ADD A, A
		delta = p.a
	case 0xC6: // ADD A, n8
		delta = p.readOperand8(p.pc, op.mode)
	}
	result16 := uint16(p.a) + uint16(delta)
	p.a = byte(result16 & 0xff)
	p.setFlag(carryFlag, result16 > 0xff)
	p.setFlag(zeroFlag, p.a == 0)
	p.setFlag(subFlag, false)
	p.determineHalfCarry(original, delta)
}

// ADC A, N; N = r8, (HL), n8
// A = A + N + carry; z0hc
func addAWithCarry(p *Processor, op *instruction) {
	original := p.a
	var delta, carry byte
	switch op.code {
	case 0x88:
		delta = p.b
	case 0x89:
		delta = p.c
	case 0x8A:
		delta = p.d
	case 0x8B:
		delta = p.e
	case 0x8C:
		delta = p.h
	case 0x8D:
		delta = p.l
	case 0x8E: // ADC A, (HL)
		delta = p.readMem8(p.reg16(p.h, p.l))
	case 0x8F: // ADC A, A
		delta = p.a
	case 0xCE: // ADC A, n8
		delta = p.readOperand8(p.pc, op.mode)
	}
	if p.getFlag(carryFlag) {
		carry = 1
	}
	result16 := uint16(p.a) + uint16(delta) + uint16(carry)
	p.a = byte(result16 & 0xff)
	p.setFlag(carry, result16 > 0xff)
	p.setFlag(zeroFlag, p.a == 0)
	p.setFlag(subFlag, false)
	p.determineHalfCarry(original, delta+carry)
}
