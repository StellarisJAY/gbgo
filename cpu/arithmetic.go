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
	p.determineHalfCarry(original, result)
}
