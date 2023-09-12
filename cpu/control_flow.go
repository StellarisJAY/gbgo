package cpu

func (p *Processor) conditionalJump(target uint16, condition bool) {
	if condition {
		p.pc = target
	}
}

func (p *Processor) conditionalJumpRelative(offset uint16, condition bool) {
	target := int32(p.pc) + int32(int16(offset))
	if condition {
		p.pc = uint16(target)
	}
}

func jp(p *Processor, op *instruction) {
	target := p.readOperand16(p.pc, op.mode)
	p.pc = target
}

func jpHL(p *Processor, _ *instruction) {
	p.pc = p.reg16(p.h, p.l)
}

func jpc(p *Processor, op *instruction) {
	target := p.readOperand16(p.pc, op.mode)
	condition := false
	switch op.code {
	case 0xC2:
		condition = !p.getFlag(zeroFlag)
	case 0xCA:
		condition = p.getFlag(zeroFlag)
	case 0xD2:
		condition = !p.getFlag(carryFlag)
	case 0xDA:
		condition = p.getFlag(carryFlag)
	}
	p.conditionalJump(target, condition)
}

func jr(p *Processor, op *instruction) {
	offset := p.readOperand16(p.pc, op.mode)
	p.conditionalJumpRelative(offset, true)
}

func jrc(p *Processor, op *instruction) {
	offset := p.readOperand16(p.pc, op.mode)
	condition := false
	switch op.code {
	case 0x20:
		condition = !p.getFlag(zeroFlag)
	case 0x28:
		condition = p.getFlag(zeroFlag)
	case 0x30:
		condition = !p.getFlag(carryFlag)
	case 0x38:
		condition = p.getFlag(carryFlag)
	}
	p.conditionalJumpRelative(offset, condition)
}
