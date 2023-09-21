package interrupt

type Register struct {
	data byte
}

type Code byte
type Requester func(Code)

const (
	VBlankInterrupt Code = 1 << iota
	LCDStatInterrupt
	TimerInterrupt
	SerialInterrupt
	JoyPadInterrupt
)

func NewRegister() *Register {
	return &Register{0}
}

func (r *Register) Get(code Code) bool {
	return r.data&byte(code) != 0
}

func (r *Register) Set(code Code, stat bool) {
	if stat {
		r.data |= byte(code)
	} else {
		r.data &= ^byte(code)
	}
}

func (r *Register) Clear() {
	r.data = 0
}

func (r *Register) Poll() byte {
	result := r.data
	r.data = 0
	return result
}
