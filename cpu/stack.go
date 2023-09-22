package cpu

// ld SP, HL
func loadSP(p *Processor, _ *Instruction) {
	p.sp = p.reg16(p.h, p.l)
}

// ld (nn), SP
func saveSP(p *Processor, op *Instruction) {
	addr := p.readOperand16(p.pc, op.mode)
	p.writeMem16(addr, p.sp)
}

func (p *Processor) stackPush16(data uint16) {
	p.sp = p.sp - 2
	p.writeMem16(p.sp, data)
}

func (p *Processor) stackPop16() uint16 {
	data := p.readMem16(p.sp)
	p.sp = p.sp + 2
	return data
}

// push
func pushReg(p *Processor, op *Instruction) {
	switch op.code {
	case 0xC5:
		p.stackPush16(p.reg16(p.b, p.c))
	case 0xD5:
		p.stackPush16(p.reg16(p.d, p.e))
	case 0xE5:
		p.stackPush16(p.reg16(p.h, p.l))
	case 0xF5:
		p.stackPush16(p.reg16(p.a, p.f))
	}
}

// pop
func popReg(p *Processor, op *Instruction) {
	switch op.code {
	case 0xC1:
		p.writeBC(p.stackPop16())
	case 0xD1:
		p.writeDE(p.stackPop16())
	case 0xE1:
		p.writeHL(p.stackPop16())
	case 0xF1:
		p.writeAF(p.stackPop16())
	}
}

// ADD SP, n8
func addSP(p *Processor, op *Instruction) {
	delta := p.readOperand8(p.pc, op.mode)
	if delta&0x80 == 0 {
		// positive
		p.sp += uint16(delta)
	} else {
		// negative
		p.sp -= uint16(-int16(int8(delta)))
	}
	p.setFlag(zeroFlag, false)
	p.setFlag(subFlag, false)
}
