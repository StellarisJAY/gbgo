package ppu

type LCDStatus struct {
	data byte
}

const (
	modeFlag byte = 1
)

const (
	hBlankMode byte = iota
	vBlankMode
	searchingOAMMode
	transferDataToLCDCMode
)

const (
	lycEqualFlag byte = 1 << (iota + 2)
	hBlankSource
	vBlankSource
	oamSource
	lycEqualSource
)

func (s *LCDStatus) get(attribute byte) bool {
	return s.data&attribute != 0
}

func (s *LCDStatus) set(attribute byte, stat bool) {
	if stat {
		s.data |= attribute
	} else {
		s.data &= ^attribute
	}
}
