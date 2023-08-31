package cpu

// instruction 指令信息
type instruction struct {
	code    byte               // opcode
	name    string             // 指令别名
	length  uint16             // 指令长度，字节
	cycles  uint64             // cpu cycles
	mode    memoryMode         // 内存访问模式
	handler instructionHandler // 指令处理函数
}

// ProcessorContext cpu当前的上下文，在callback中使用
type ProcessorContext struct {
}

type instructionHandler func(p *Processor, op *instruction)
type InstructionCallback func(ctx ProcessorContext, op *instruction)

// execute 执行指令并调用回调函数
func (ins *instruction) execute(p *Processor, callback InstructionCallback) {
	ins.handler(p, ins)
	if callback != nil {
		callback(p.getContext(), ins)
	}
}

// getContext 获取当前cpu上下文
func (p *Processor) getContext() ProcessorContext {
	return ProcessorContext{}
}

// 指令集，按照指令类型排序
var instructionSet = map[byte]*instruction{
	// ld r8, n8
	0x06: {0x06, "LD", 2, 8, immediate, loadImmediate8},
	0x0E: {0x0E, "LD", 2, 8, immediate, loadImmediate8},
	0x16: {0x16, "LD", 2, 8, immediate, loadImmediate8},
	0x1E: {0x1E, "LD", 2, 8, immediate, loadImmediate8},
	0x26: {0x26, "LD", 2, 8, immediate, loadImmediate8},
	0x2E: {0x2E, "LD", 2, 8, immediate, loadImmediate8},
	// ld r16, n16
	0x01: {0x01, "LD", 3, 12, immediate, loadImmediate16},
	0x11: {0x11, "LD", 3, 12, immediate, loadImmediate16},
	0x21: {0x21, "LD", 3, 12, immediate, loadImmediate16},
	0x31: {0x31, "LD", 3, 12, immediate, loadImmediate16},
	// ld B, r8
	0x40: {0x40, "LD", 1, 4, none, loadBFromReg},
	0x41: {0x41, "LD", 1, 4, none, loadBFromReg},
	0x42: {0x42, "LD", 1, 4, none, loadBFromReg},
	0x43: {0x43, "LD", 1, 4, none, loadBFromReg},
	0x44: {0x44, "LD", 1, 4, none, loadBFromReg},
	0x45: {0x45, "LD", 1, 4, none, loadBFromReg},
	// ld C, r8
	0x48: {0x48, "LD", 1, 4, none, loadCFromReg},
	0x49: {0x49, "LD", 1, 4, none, loadCFromReg},
	0x4A: {0x4A, "LD", 1, 4, none, loadCFromReg},
	0x4B: {0x4B, "LD", 1, 4, none, loadCFromReg},
	0x4C: {0x4C, "LD", 1, 4, none, loadCFromReg},
	0x4D: {0x4D, "LD", 1, 4, none, loadCFromReg},
	// ld D, r8
	0x50: {0x50, "LD", 1, 4, none, loadDFromReg},
	0x51: {0x51, "LD", 1, 4, none, loadDFromReg},
	0x52: {0x52, "LD", 1, 4, none, loadDFromReg},
	0x53: {0x53, "LD", 1, 4, none, loadDFromReg},
	0x54: {0x54, "LD", 1, 4, none, loadDFromReg},
	0x55: {0x55, "LD", 1, 4, none, loadDFromReg},
	// ld E, r8
	0x58: {0x58, "LD", 1, 4, none, loadEFromReg},
	0x59: {0x59, "LD", 1, 4, none, loadEFromReg},
	0x5A: {0x5A, "LD", 1, 4, none, loadEFromReg},
	0x5B: {0x5B, "LD", 1, 4, none, loadEFromReg},
	0x5C: {0x5C, "LD", 1, 4, none, loadEFromReg},
	0x5D: {0x5D, "LD", 1, 4, none, loadEFromReg},
	// ld H, r8
	0x60: {0x60, "LD", 1, 4, none, loadHFromReg},
	0x61: {0x61, "LD", 1, 4, none, loadHFromReg},
	0x62: {0x62, "LD", 1, 4, none, loadHFromReg},
	0x63: {0x63, "LD", 1, 4, none, loadHFromReg},
	0x64: {0x64, "LD", 1, 4, none, loadHFromReg},
	0x65: {0x65, "LD", 1, 4, none, loadHFromReg},
	// ld L, r8
	0x68: {0x68, "LD", 1, 4, none, loadLFromReg},
	0x69: {0x69, "LD", 1, 4, none, loadLFromReg},
	0x6A: {0x6A, "LD", 1, 4, none, loadLFromReg},
	0x6B: {0x6B, "LD", 1, 4, none, loadLFromReg},
	0x6C: {0x6C, "LD", 1, 4, none, loadLFromReg},
	0x6D: {0x6D, "LD", 1, 4, none, loadLFromReg},
	// ld r8, (HL)
	0x46: {0x46, "LD", 1, 8, none, loadRegFromHL},
	0x4E: {0x4E, "LD", 1, 8, none, loadRegFromHL},
	0x56: {0x56, "LD", 1, 8, none, loadRegFromHL},
	0x5E: {0x5E, "LD", 1, 8, none, loadRegFromHL},
	0x66: {0x66, "LD", 1, 8, none, loadRegFromHL},
	0x6E: {0x6E, "LD", 1, 8, none, loadRegFromHL},
	// ld (HL), r8
	0x70: {0x70, "LD", 1, 8, none, storeRegInHL},
	0x71: {0x71, "LD", 1, 8, none, storeRegInHL},
	0x72: {0x72, "LD", 1, 8, none, storeRegInHL},
	0x73: {0x73, "LD", 1, 8, none, storeRegInHL},
	0x74: {0x74, "LD", 1, 8, none, storeRegInHL},
	0x75: {0x75, "LD", 1, 8, none, storeRegInHL},
	// ld A, (NN)
	0x0A: {0x0A, "LD", 1, 8, none, loadA},
	0x1A: {0x0A, "LD", 1, 8, none, loadA},
	0x7E: {0x7E, "LD", 1, 8, none, loadA},
	0xFA: {0xFA, "LD", 3, 16, none, loadA},
	// ld (N), A
	0x02: {0x02, "LD", 1, 8, none, storeA},
	0x12: {0x12, "LD", 1, 8, none, storeA},
	0x77: {0x77, "LD", 1, 8, none, storeA},
	0xEA: {0xEA, "LD", 3, 16, none, storeA},
	// ld r8, A
	0x47: {0x47, "LD", 1, 4, none, storeAInReg},
	0x4F: {0x4F, "LD", 1, 4, none, storeAInReg},
	0x57: {0x57, "LD", 1, 4, none, storeAInReg},
	0x5F: {0x5F, "LD", 1, 4, none, storeAInReg},
	0x67: {0x67, "LD", 1, 4, none, storeAInReg},
	0x6F: {0x6F, "LD", 1, 4, none, storeAInReg},
	// ld A, r8
	0x78: {0x78, "LD", 1, 4, none, loadAFromReg},
	0x79: {0x79, "LD", 1, 4, none, loadAFromReg},
	0x7A: {0x7A, "LD", 1, 4, none, loadAFromReg},
	0x7B: {0x7B, "LD", 1, 4, none, loadAFromReg},
	0x7C: {0x7C, "LD", 1, 4, none, loadAFromReg},
	0x7D: {0x7D, "LD", 1, 4, none, loadAFromReg},
	0x7F: {0x7F, "LD", 1, 4, none, loadAFromReg},
	// ld (HL), n8
	0x36: {0x36, "LD", 2, 12, immediate, storeImmediateInHL},
	// read IO Port
	0xF0: {0xF0, "LD", 2, 12, immediate, readIOPort},
	0xF2: {0xF2, "LD", 1, 8, none, readIOPort},
	// write IO Port
	0xE0: {0xE0, "LD", 2, 12, immediate, writeIOPort},
	0xE2: {0xE2, "LD", 1, 8, none, writeIOPort},
	// ldi
	0x22: {0x22, "LDI", 1, 8, none, ldi},
	0x2A: {0x2A, "LDI", 1, 8, none, ldi},
	// ldd
	0x32: {0x32, "LDD", 1, 8, none, ldd},
	0x3A: {0x3A, "LDD", 1, 8, none, ldd},
	// ld SP, HL
	0xF9: {0xF9, "LD", 1, 8, none, loadSP},
	// stack push
	0xC5: {0xC5, "PUSH", 1, 12, none, pushReg},
	0xD5: {0xD5, "PUSH", 1, 12, none, pushReg},
	0xE5: {0xE5, "PUSH", 1, 12, none, pushReg},
	0xF5: {0xF5, "PUSH", 1, 12, none, pushReg},
	// stack pop
	0xC1: {0xC1, "PUSH", 1, 12, none, popReg},
	0xD1: {0xD1, "PUSH", 1, 12, none, popReg},
	0xE1: {0xE1, "PUSH", 1, 12, none, popReg},
	0xF1: {0xF1, "PUSH", 1, 12, none, popReg},
}
