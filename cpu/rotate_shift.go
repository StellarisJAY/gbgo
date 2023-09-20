package cpu

func rotateA(p *Processor, op *Instruction) {
	switch op.code {
	case 0x07:
		p.a = p.rotateLeft(p.a)
	case 0x17:
		p.a = p.rotateLeftCarry(p.a)
	case 0x0f:
		p.a = p.rotateRight(p.a)
	case 0x1f:
		p.a = p.rotateRightCarry(p.a)
	}
}

// rlc,rl, rrc, rr等opcode都是0xCB, 它们通过后面的立即数来编号
func rotatesAndShifts(p *Processor, _ *Instruction) {
	code := p.readOperand8(p.pc, immediate)
	p.pc += 1
	switch {
	case code <= 0x07:
		rlc(p, code)
	case code >= 0x08 && code <= 0x0F:
		rrc(p, code)
	case code >= 0x10 && code <= 0x17:
		rl(p, code)
	case code >= 0x18 && code <= 0x1F:
		rr(p, code)
	case code >= 0x20 && code <= 0x27:
		sla(p, code)
	case code >= 0x28 && code <= 0x2F:
		sra(p, code)
	case code >= 0x30 && code <= 0x37:
		swap(p, code)
	case code >= 0x38 && code <= 0x3F:
		srl(p, code)
	case code >= 0x40 && code <= 0x47:
		bit(p, code)
	case code >= 0x80 && code <= 0x87:
		resetBit(p, code)
	case code >= 0xC0 && code <= 0xC7:
		setBit(p, code)
	}
}

func rlc(p *Processor, code byte) {
	switch code {
	case 0x00:
		p.b = p.rotateLeft(p.b)
	case 0x01:
		p.c = p.rotateLeft(p.c)
	case 0x02:
		p.d = p.rotateLeft(p.d)
	case 0x03:
		p.e = p.rotateLeft(p.e)
	case 0x04:
		p.h = p.rotateLeft(p.h)
	case 0x05:
		p.l = p.rotateLeft(p.l)
	case 0x06:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, p.rotateLeft)
	case 0x07:
		p.a = p.rotateLeft(p.a)
	}
}

func rl(p *Processor, code byte) {
	switch code {
	case 0x10:
		p.b = p.rotateLeftCarry(p.b)
	case 0x11:
		p.c = p.rotateLeftCarry(p.c)
	case 0x12:
		p.d = p.rotateLeftCarry(p.d)
	case 0x13:
		p.e = p.rotateLeftCarry(p.e)
	case 0x14:
		p.h = p.rotateLeftCarry(p.h)
	case 0x15:
		p.l = p.rotateLeftCarry(p.l)
	case 0x16:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, p.rotateLeftCarry)
	case 0x17:
		p.a = p.rotateLeftCarry(p.a)
	}
}

func rrc(p *Processor, code byte) {
	switch code {
	case 0x08:
		p.b = p.rotateRight(p.b)
	case 0x09:
		p.c = p.rotateRight(p.c)
	case 0x0A:
		p.d = p.rotateRight(p.d)
	case 0x0B:
		p.e = p.rotateRight(p.e)
	case 0x0C:
		p.h = p.rotateRight(p.h)
	case 0x0D:
		p.l = p.rotateRight(p.l)
	case 0x0E:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, p.rotateRight)
	case 0x0F:
		p.a = p.rotateRight(p.a)
	}
}

func rr(p *Processor, code byte) {
	switch code {
	case 0x18:
		p.b = p.rotateRightCarry(p.b)
	case 0x19:
		p.c = p.rotateRightCarry(p.c)
	case 0x1A:
		p.d = p.rotateRightCarry(p.d)
	case 0x1B:
		p.e = p.rotateRightCarry(p.e)
	case 0x1C:
		p.h = p.rotateRightCarry(p.h)
	case 0x1D:
		p.l = p.rotateRightCarry(p.l)
	case 0x1E:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, p.rotateRightCarry)
	case 0x1F:
		p.a = p.rotateRightCarry(p.a)
	}
}

// SLA N; N = r8
// N = N << 1
func sla(p *Processor, code byte) {
	switch code {
	case 0x20:
		p.b = p.shiftLeft(p.b)
	case 0x21:
		p.c = p.shiftLeft(p.c)
	case 0x22:
		p.d = p.shiftLeft(p.d)
	case 0x23:
		p.e = p.shiftLeft(p.e)
	case 0x24:
		p.h = p.shiftLeft(p.h)
	case 0x25:
		p.l = p.shiftLeft(p.l)
	case 0x26:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, p.shiftLeft)
	case 0x27:
		p.a = p.shiftLeft(p.a)
	}
}

// SRA N; N = r8
// N = N >> 1; 不改变最高位
func sra(p *Processor, code byte) {
	switch code {
	case 0x28:
		p.b = p.shiftRight(p.b, false)
	case 0x29:
		p.c = p.shiftRight(p.c, false)
	case 0x2A:
		p.d = p.shiftRight(p.d, false)
	case 0x2B:
		p.e = p.shiftRight(p.e, false)
	case 0x2C:
		p.h = p.shiftRight(p.h, false)
	case 0x2D:
		p.l = p.shiftRight(p.l, false)
	case 0x2E:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, func(val byte) byte {
			return p.shiftRight(val, false)
		})
	case 0x2F:
		p.a = p.shiftRight(p.a, false)
	}
}

// SRL N; N=r8
// N = N >> 1; 最高位变为0
func srl(p *Processor, code byte) {
	switch code {
	case 0x38:
		p.b = p.shiftRight(p.b, true)
	case 0x39:
		p.c = p.shiftRight(p.c, true)
	case 0x3A:
		p.d = p.shiftRight(p.d, true)
	case 0x3B:
		p.e = p.shiftRight(p.e, true)
	case 0x3C:
		p.h = p.shiftRight(p.h, true)
	case 0x3D:
		p.l = p.shiftRight(p.l, true)
	case 0x3E:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, func(val byte) byte {
			return p.shiftRight(val, true)
		})
	case 0x3F:
		p.a = p.shiftRight(p.a, true)
	}
}

// SWAP N
func swap(p *Processor, code byte) {
	switch code {
	case 0x30:
		p.b = p.swapHighLow(p.b)
	case 0x31:
		p.c = p.swapHighLow(p.c)
	case 0x32:
		p.d = p.swapHighLow(p.d)
	case 0x33:
		p.e = p.swapHighLow(p.e)
	case 0x34:
		p.h = p.swapHighLow(p.h)
	case 0x35:
		p.l = p.swapHighLow(p.l)
	case 0x36:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, p.swapHighLow)
	case 0x37:
		p.a = p.swapHighLow(p.a)
	}
}

// BIT n, N; n = 0~7
func bit(p *Processor, code byte) {
	n := p.readOperand8(p.pc, immediate)
	p.pc += 1
	switch code {
	case 0x40:
		p.testBit(p.b, n)
	case 0x41:
		p.testBit(p.c, n)
	case 0x42:
		p.testBit(p.d, n)
	case 0x43:
		p.testBit(p.e, n)
	case 0x44:
		p.testBit(p.h, n)
	case 0x45:
		p.testBit(p.l, n)
	case 0x46:
		val := p.readMem8(p.reg16(p.h, p.l))
		p.testBit(val, n)
	case 0x47:
		p.testBit(p.a, n)
	}
}

// SET n, N; n = 0~7
func setBit(p *Processor, code byte) {
	n := p.readOperand8(p.pc, immediate)
	p.pc += 1
	switch code {
	case 0xC0:
		p.b = p.setBit(p.b, n)
	case 0xC1:
		p.c = p.setBit(p.c, n)
	case 0xC2:
		p.d = p.setBit(p.d, n)
	case 0xC3:
		p.e = p.setBit(p.e, n)
	case 0xC4:
		p.h = p.setBit(p.h, n)
	case 0xC5:
		p.l = p.setBit(p.l, n)
	case 0xC6:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, func(val byte) byte {
			return p.setBit(val, n)
		})
	case 0xC7:
		p.a = p.setBit(p.a, n)
	}
}

func resetBit(p *Processor, code byte) {
	n := p.readOperand8(p.pc, immediate)
	p.pc += 1
	switch code {
	case 0x80:
		p.b = p.resetBit(p.b, n)
	case 0x81:
		p.c = p.resetBit(p.c, n)
	case 0x82:
		p.d = p.resetBit(p.d, n)
	case 0x83:
		p.e = p.resetBit(p.e, n)
	case 0x84:
		p.h = p.resetBit(p.h, n)
	case 0x85:
		p.l = p.resetBit(p.l, n)
	case 0x86:
		addr := p.reg16(p.h, p.l)
		p.modifyMemory8(addr, func(val byte) byte {
			return p.resetBit(val, n)
		})
	case 0x87:
		p.a = p.resetBit(p.a, n)
	}
}

func (p *Processor) rotateLeft(val byte) byte {
	// carry设置为旧的bit 7
	p.setFlag(carryFlag, val&0x80 != 0)
	result := val << 1
	p.setFlag(zeroFlag, result == 0)
	p.setFlag(subFlag, false)
	p.setFlag(halfCarryFlag, false)
	return result
}

func (p *Processor) rotateLeftCarry(val byte) byte {
	// 当前的carryFlag写入结果的bit0
	var carry byte = 0
	if p.getFlag(carryFlag) {
		carry = 1
	}
	// carry设置为旧的bit 7
	p.setFlag(carryFlag, val&0x80 != 0)
	result := val<<1 | carry
	p.setFlag(zeroFlag, result == 0)
	p.setFlag(subFlag, false)
	p.setFlag(halfCarryFlag, false)
	return result
}

func (p *Processor) rotateRight(val byte) byte {
	// carryFlag设置为旧的bit0
	p.setFlag(carryFlag, val&1 != 0)
	result := val >> 1
	p.setFlag(zeroFlag, result == 0)
	p.setFlag(subFlag, false)
	p.setFlag(halfCarryFlag, false)
	return result
}

func (p *Processor) rotateRightCarry(val byte) byte {
	// 当前的carryFlag写入结果的bit7
	var carry byte = 0
	if p.getFlag(carryFlag) {
		carry = 0b1000_0000
	}
	// carryFlag设置为旧的bit0
	p.setFlag(carryFlag, val&1 != 0)
	result := val>>1 | carry
	p.setFlag(zeroFlag, result == 0)
	p.setFlag(subFlag, false)
	p.setFlag(halfCarryFlag, false)
	return result
}

func (p *Processor) shiftLeft(val byte) byte {
	// carry设置为旧的bit7
	p.setFlag(carryFlag, val&0x80 != 0)
	result := val << 1
	p.setFlag(zeroFlag, result == 0)
	p.setFlag(subFlag, false)
	p.setFlag(halfCarryFlag, false)
	return result
}

// shiftRight 右移，keepMSB是否保留最高位
func (p *Processor) shiftRight(val byte, changeMSB bool) byte {
	var msb byte
	// 保留最高位
	if !changeMSB {
		msb = val & 0x80
	}
	// carry设置为旧的bit0
	p.setFlag(carryFlag, val&1 != 0)
	result := val>>1 | msb
	p.setFlag(zeroFlag, result == 0)
	p.setFlag(subFlag, false)
	p.setFlag(halfCarryFlag, false)
	return result
}

// swap 交换高低半字节
func (p *Processor) swapHighLow(val byte) byte {
	high, low := val>>4, val&0xf
	result := low<<4 | high
	p.setFlag(zeroFlag, result == 0)
	p.setFlag(carryFlag, false)
	p.setFlag(subFlag, false)
	p.setFlag(halfCarryFlag, false)
	return result
}

// testBit 获取val的bit n
func (p *Processor) testBit(val, n byte) {
	result := val&(1<<n) == 0
	p.setFlag(zeroFlag, result)
	p.setFlag(halfCarryFlag, true)
	p.setFlag(subFlag, false)
}

// setBit 设置val的bit n
func (p *Processor) setBit(val, n byte) byte {
	return val | (1 << n)
}

// resetBit 清除val的bit n
func (p *Processor) resetBit(val, n byte) byte {
	return val & ^(1 << n)
}
