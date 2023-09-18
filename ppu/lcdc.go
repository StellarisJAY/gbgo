package ppu

type LDCControl struct {
	data byte
}

const (
	bgWindowEnable   byte = 1 << iota // 是否显示背景和窗口
	objEnable                         // 是否显示对象
	objSize                           // 对象大小，8x8,8x16
	bgTileMap                         // 背景tile map地址，0: 0x9800~0x9BFF,1:0x9C00~0x9FFF
	bgWindowTileData                  //  背景和窗口tile数据地址, 8800~97FF,8000~8FFF
	windowEnable                      // 是否显示窗口
	windowTileMap                     // 窗口tile map, 0: 0x9800~0x9BFF,1:0x9C00~0x9FFF
	ppuEnable                         // 是否开启ppu和LCD
)

func (c *LDCControl) get(attribute byte) bool {
	return c.data&attribute != 0
}

func (c *LDCControl) set(attribute byte) {
	c.data = c.data | attribute
}

func (c *LDCControl) reset(attribute byte) {
	c.data = c.data & ^attribute
}
