package cpu

// inc N; N=r8,(HL)
// N = N + 1, z0h-
func inc(p *Processor, op *Instruction) {
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
func dec(p *Processor, op *Instruction) {
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
func addA(p *Processor, op *Instruction) {
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
func addAWithCarry(p *Processor, op *Instruction) {
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

// SUB A, N; N = r8,n8,(HL)
// A = A - N; z1hc
func subA(p *Processor, op *Instruction) {
	original := p.a
	var delta byte
	switch op.code {
	case 0x90:
		delta = p.b
	case 0x91:
		delta = p.c
	case 0x92:
		delta = p.d
	case 0x93:
		delta = p.e
	case 0x94:
		delta = p.h
	case 0x95:
		delta = p.l
	case 0x96: // SUB A, (HL)
		delta = p.readMem8(p.reg16(p.h, p.l))
	case 0x97: // SUB A, A
		delta = p.a
	case 0xD6: // SUB A, n8
		delta = p.readOperand8(p.pc, op.mode)
	}
	p.a = p.a - delta
	p.setFlag(carryFlag, original < delta)
	p.setFlag(zeroFlag, p.a == 0)
	p.setFlag(subFlag, true)
	p.determineHalfCarry(original, delta)
}

// SBC A, N; N = r8,n8,(HL)
// A = A - N - carry; z1hc
func subAWithCarry(p *Processor, op *Instruction) {
	original := p.a
	var delta, carry byte
	switch op.code {
	case 0x98:
		delta = p.b
	case 0x99:
		delta = p.c
	case 0x9A:
		delta = p.d
	case 0x9B:
		delta = p.e
	case 0x9C:
		delta = p.h
	case 0x9D:
		delta = p.l
	case 0x9E: // SBC A, (HL)
		delta = p.readMem8(p.reg16(p.h, p.l))
	case 0x9F: // SBC A, A
		delta = p.a
	case 0xDE: // SBC A, n8
		delta = p.readOperand8(p.pc, op.mode)
	}
	if p.getFlag(carryFlag) {
		carry = 1
	}
	p.a = p.a - delta - carry
	p.setFlag(carryFlag, original < delta+carry)
	p.setFlag(zeroFlag, p.a == 0)
	p.setFlag(subFlag, true)
	p.determineHalfCarry(original, delta)
}

// CP A,N; N = r8,n8,(HL)
// z1hc
// c: A < N
// z: A == N
func compareA(p *Processor, op *Instruction) {
	var data byte
	switch op.code {
	case 0xB8:
		data = p.b
	case 0xB9:
		data = p.c
	case 0xBA:
		data = p.d
	case 0xBB:
		data = p.e
	case 0xBC:
		data = p.h
	case 0xBD:
		data = p.l
	case 0xBE:
		data = p.readMem8(p.reg16(p.h, p.l))
	case 0xBF:
		data = p.a
	case 0xFE:
		data = p.readOperand8(p.pc, op.mode)
	}
	p.compare(p.a, data)
}

func (p *Processor) compare(val1, val2 byte) {
	p.setFlag(zeroFlag, val1 == val2)
	p.setFlag(subFlag, true)
	p.setFlag(carryFlag, val1 < val2)
	p.determineHalfCarry(val1, val2)
}

// ADD HL, N; N = BC,DE,HL,SP
// -0hc
// h: carry from bit 11
// c: carry from bit 15
func addHL(p *Processor, op *Instruction) {
	original := p.reg16(p.h, p.l)
	var delta uint16
	switch op.code {
	case 0x09:
		delta = p.reg16(p.b, p.c)
	case 0x19:
		delta = p.reg16(p.d, p.e)
	case 0x29:
		delta = p.reg16(p.h, p.l)
	case 0x39:
		delta = p.sp
	}
	result32 := uint32(original) + uint32(delta)
	p.setFlag(carryFlag, result32 > 0xffff)
	p.setFlag(subFlag, false)
	// 第11位是否有进位
	p.setFlag(halfCarryFlag, original&0xfff+delta > 0xfff)
	p.writeHL(original + delta)
}

// INC r16
func inc16(p *Processor, op *Instruction) {
	switch op.code {
	case 0x03:
		p.writeBC(p.reg16(p.b, p.c) + 1)
	case 0x13:
		p.writeDE(p.reg16(p.d, p.e) + 1)
	case 0x23:
		p.writeHL(p.reg16(p.h, p.l) + 1)
	case 0x33:
		p.sp += 1
	}
}

// DEC r16
func dec16(p *Processor, op *Instruction) {
	switch op.code {
	case 0x0B:
		p.writeBC(p.reg16(p.b, p.c) - 1)
	case 0x1B:
		p.writeDE(p.reg16(p.d, p.e) - 1)
	case 0x2B:
		p.writeHL(p.reg16(p.h, p.l) - 1)
	case 0x3B:
		p.sp -= 1
	}
}

func cpl(p *Processor, _ *Instruction) {
	p.a = ^p.a
	p.setFlag(halfCarryFlag, true)
	p.setFlag(subFlag, true)
}

func daa(p *Processor, _ *Instruction) {
	if !p.getFlag(subFlag) {
		if p.getFlag(carryFlag) || p.a > 0x99 {
			p.a += 0x60
			p.setFlag(carryFlag, true)
		}
		if p.getFlag(halfCarryFlag) || p.a&0x0f > 0x09 {
			p.a += 0x06
			p.setFlag(halfCarryFlag, false)
		}
	} else if p.getFlag(carryFlag) && p.getFlag(halfCarryFlag) {
		p.a = 0x9A
		p.setFlag(halfCarryFlag, false)
	} else if p.getFlag(carryFlag) {
		p.a = 0xA0
	} else if p.getFlag(halfCarryFlag) {
		p.a = 0xFA
		p.setFlag(halfCarryFlag, false)
	}
}
