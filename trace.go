package main

import (
	"fmt"
	"github.com/StellarisJAY/gbgo/cpu"
)

func logInstruction(ctx cpu.ProcessorContext, ins *cpu.Instruction) {
	fmt.Printf("%04X  %02X\t%6s\t%04X %04X %04X %04X %d\n", ctx.PC, ins.Code(), ins.Name(), ctx.AF, ctx.BC, ctx.DE, ctx.HL, ctx.Cycles)
}
