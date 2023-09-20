package cpu

func (p *Processor) restart(vector uint16) {
	p.stackPush16(p.pc)
	// jump to 0x8000 + n
	p.pc = 0x8000 + vector
}

// rst: 重启程序，跳转到restart地址：0x00,0x08...0x30,0x38
func rst(p *Processor, op *Instruction) {
	switch op.code {
	case 0xC7:
		p.restart(0x00)
	case 0xCF:
		p.restart(0x08)
	case 0xD7:
		p.restart(0x10)
	case 0xDF:
		p.restart(0x18)
	case 0xE7:
		p.restart(0x20)
	case 0xEF:
		p.restart(0x28)
	case 0xF7:
		p.restart(0x30)
	case 0xFF:
		p.restart(0x38)
	}
}

func nop(_ *Processor, _ *Instruction) {}

func enableInterrupt(p *Processor, _ *Instruction) {
	p.pendingInterruptSwitch = 1
	p.nextInterruptEnable = true
}

func disableInterrupt(p *Processor, _ *Instruction) {
	p.pendingInterruptSwitch = 1
	p.nextInterruptEnable = false
}

// carryFlag取反
func ccf(p *Processor, _ *Instruction) {
	carry := p.getFlag(carryFlag)
	if carry {
		p.setFlag(carryFlag, false)
	} else {
		p.setFlag(carryFlag, true)
	}
}

// 设置carryFlag
func scf(p *Processor, _ *Instruction) {
	p.setFlag(carryFlag, true)
}
