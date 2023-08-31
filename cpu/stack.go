package cpu

// ld SP, HL
func loadSP(p *Processor, _ *instruction) {
	p.sp = p.reg16(p.h, p.l)
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
func pushReg(p *Processor, op *instruction) {
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
func popReg(p *Processor, op *instruction) {
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
