package cpu

func (p *Processor) conditionalCall(target uint16, condition bool) {
	if condition {
		// push return address
		ra := p.pc + 2
		p.stackPush16(ra)
		p.pc = target
	}
}

func (p *Processor) conditionalReturn(condition bool) {
	if condition {
		// pop return address
		ra := p.stackPop16()
		p.pc = ra
	}
}

func call(p *Processor, op *Instruction) {
	target := p.readOperand16(p.pc, op.mode)
	p.conditionalCall(target, true)
}

func callC(p *Processor, op *Instruction) {
	target := p.readOperand16(p.pc, op.mode)
	condition := false
	switch op.code {
	case 0xC4:
		condition = !p.getFlag(zeroFlag)
	case 0xCC:
		condition = p.getFlag(zeroFlag)
	case 0xD4:
		condition = !p.getFlag(carryFlag)
	case 0xDC:
		condition = p.getFlag(carryFlag)
	}
	p.conditionalCall(target, condition)
}

func ret(p *Processor, _ *Instruction) {
	p.conditionalReturn(true)
}

func retc(p *Processor, op *Instruction) {
	condition := false
	switch op.code {
	case 0xC0:
		condition = !p.getFlag(zeroFlag)
	case 0xC8:
		condition = p.getFlag(zeroFlag)
	case 0xD0:
		condition = !p.getFlag(carryFlag)
	case 0xD8:
		condition = p.getFlag(carryFlag)
	}
	p.conditionalReturn(condition)
}
