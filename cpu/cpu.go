package cpu

import (
	"fmt"
	"github.com/StellarisJAY/gbgo/bus"
)

type Processor struct {
	a  byte // accumulator
	b  byte
	c  byte
	d  byte
	e  byte
	f  byte // flags
	h  byte
	l  byte
	sp uint16 // stack pointer
	pc uint16 // program counter

	bus *bus.Bus
}

const (
	// flag 寄存器的4~7位
	carryFlag     byte = 1 << (iota + 4)
	halfCarryFlag      // 半字节carry
	subFlag
	zeroFlag
)

// memoryMode 指令寻址模式
type memoryMode byte

const (
	// 立即数，从pc地址读取操作数
	immediate memoryMode = iota
	// 直接寻址，从pc位置读取地址，再从地址读取操作数
	absolute
	// 无寻址
	none
)

func (p *Processor) run() {
	for {
		if p.pc == 0 {
			break
		}
		oldPc := p.pc
		opCode := p.readOperand8(p.pc, immediate)
		p.pc++
		ins, exists := instructionSet[opCode]
		if !exists {
			panic(fmt.Errorf("unknown opcode at 0x%x:  %x", oldPc, opCode))
		}
		ins.execute(p, nil)

		if oldPc+1 == p.pc {
			p.pc = oldPc + ins.length
		}
	}
}

func (p *Processor) readOperand8(pc uint16, mode memoryMode) byte {
	switch mode {
	case immediate:
		return p.readMem8(pc + 1)
	case absolute:
		addr := p.readMem16(pc + 1)
		return p.readMem8(addr)
	case none:
		return 0
	default:
		return 0
	}
}

func (p *Processor) readOperand16(pc uint16, mode memoryMode) uint16 {
	switch mode {
	case immediate:
		return p.readMem16(pc + 1)
	case absolute:
		addr := p.readMem16(pc + 1)
		return p.readMem16(addr)
	case none:
		return 0
	default:
		return 0
	}
}

func (p *Processor) readMem8(addr uint16) byte {
	return p.bus.ReadMem8(addr)
}

func (p *Processor) readMem16(addr uint16) uint16 {
	return p.bus.ReadMem16(addr)
}

func (p *Processor) writeMem8(addr uint16, data byte) {
	p.bus.WriteMem8(addr, data)
}

func (p *Processor) writeMem16(addr, data uint16) {
	low, high := byte(data&0xFF), byte(data>>8)
	p.writeMem8(addr, low)
	p.writeMem8(addr+1, high)
}

// 获取16位组合寄存器值
func (p *Processor) reg16(regHigh, regLow byte) uint16 {
	return uint16(regHigh)<<8 | uint16(regLow)
}

func (p *Processor) writeBC(data uint16) {
	high := byte(data >> 8)
	low := byte(data & 0xFF)
	p.b, p.c = high, low
}

func (p *Processor) writeDE(data uint16) {
	high := byte(data >> 8)
	low := byte(data & 0xFF)
	p.d, p.e = high, low
}

func (p *Processor) writeHL(data uint16) {
	high := byte(data >> 8)
	low := byte(data & 0xFF)
	p.h, p.l = high, low
}

func (p *Processor) writeAF(data uint16) {
	high := byte(data >> 8)
	low := byte(data & 0xFF)
	p.a, p.f = high, low
}

func (p *Processor) setFlag(flag byte, status bool) {
	if status {
		p.f |= flag
	} else {
		p.f &= ^flag
	}
}

func (p *Processor) determineHalfCarry(original, result byte) {
	// 加法，低位半字节向高位半字节进位
	// 减法，低位半字节从高位半字节借位
	// 原值和结果值的高位半字节不相等
	halfCarry := original&0xf0 != result&0xf0
	p.setFlag(halfCarryFlag, halfCarry)
}
