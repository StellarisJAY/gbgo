package cpu

import (
	"fmt"
	"github.com/StellarisJAY/gbgo/bus"
	"time"
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

	pendingInterruptSwitch int // EI和DI都不会立即切换中断状态，都需要等EI和DI之后一条指令执行后才切换状态
	nextInterruptEnable    bool
	interruptEnabled       bool

	lastTickTime int64
	cycles       int64
}

const (
	// flag 寄存器的4~7位
	carryFlag     byte = 1 << (iota + 4)
	halfCarryFlag      // 半字节carry, 加法向第四位进位 或 减法不从第四位借位时设置
	subFlag
	zeroFlag

	cpuFrequency      int64 = 4194304 // CPU频率，4.194304MHz, 每毫秒4194周期
	cpuFrequencyMilli       = cpuFrequency / 1000
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

func MakeCPU(b *bus.Bus) *Processor {
	return &Processor{
		bus: b,
	}
}

// Tick 每渲染一帧画面CPU tick一次，通过计算两次渲染之间的时间间隔来控制CPU的指令周期
func (p *Processor) Tick(time int64, callback InstructionCallback) {
	oldCycles := p.cycles
	// 上一次tick与当前时间之间的cpu周期数
	curTickCycles := (time - p.lastTickTime) * cpuFrequencyMilli
	for p.cycles <= oldCycles+curTickCycles {
		oldPc := p.pc
		opCode := p.readOperand8(p.pc, immediate)
		p.pc++
		ins, exists := instructionSet[opCode]
		if !exists {
			panic(fmt.Errorf("unknown opcode at 0x%4X:  0x%2X", oldPc, opCode))
		}
		ins.execute(p, callback)
		if oldPc+1 == p.pc {
			p.pc = oldPc + ins.length
		}
		p.cycles += int64(ins.cycles)
		// EI和DI要等待下一条指令结束才切换interrupt状态
		if p.pendingInterruptSwitch == 0 {
			p.pendingInterruptSwitch = -1
			p.interruptEnabled = p.nextInterruptEnable
		} else if p.pendingInterruptSwitch > 0 {
			p.pendingInterruptSwitch -= 1
		}
	}
	p.lastTickTime = time
}

// Reset cpu status
func (p *Processor) Reset() {
	// entry point
	p.pc = 0x0100
	// stack bottom
	p.sp = 0xFFFE
	p.b, p.c = 0, 0
	p.d, p.e = 0, 0
	p.h, p.l = 0, 0
	p.a, p.f = 0, 0
	p.pendingInterruptSwitch = -1
	p.nextInterruptEnable = false
	p.interruptEnabled = false
	p.lastTickTime = time.Now().UnixMilli()
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

func (p *Processor) getFlag(flag byte) bool {
	return p.f&flag != 0
}

func (p *Processor) determineHalfCarry(original, delta byte) {
	// 加法，低位半字节向高位半字节进位
	// 减法，低位半字节没有从高位半字节借位
	halfCarry := original&0xf+delta > 0xf || original&0xf > delta
	p.setFlag(halfCarryFlag, halfCarry)
}

// 调用modifier修改一个内存地址
func (p *Processor) modifyMemory8(addr uint16, modifier func(byte) byte) {
	result := modifier(p.readMem8(addr))
	p.writeMem8(addr, result)
}
