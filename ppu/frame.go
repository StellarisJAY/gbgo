package ppu

const (
	Width  = 160
	Height = 144
)

type color struct {
	r byte
	g byte
	b byte
}

var frame []byte = make([]byte, Width*Height*3)

func setPixel(x, y uint32, c color) {
	first := y*Width*3 + x*3
	if first+2 < uint32(len(frame)) {
		frame[first] = c.r
		frame[first+1] = c.g
		frame[first+2] = c.b
	}
}

func FrameData() []byte {
	return frame
}

func renderStartingPage() {
	var x, y uint32 = Width/2 - 8, Height/2 - 8
	for i := uint32(0); i < 8; i++ {
		for j := uint32(0); j < 8; j++ {
			setPixel(x+i, y+j, color{255, 255, 255})
		}
	}
}
