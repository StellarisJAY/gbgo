package ppu

type PPU struct {
}

func MakePPU() *PPU {
	return &PPU{}
}

func (p *PPU) Render() {
	renderStartingPage()
}
