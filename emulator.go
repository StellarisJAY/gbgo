package main

import (
	"flag"
	"fmt"
	"github.com/StellarisJAY/gbgo/bus"
	"github.com/StellarisJAY/gbgo/cartridge"
	"github.com/StellarisJAY/gbgo/cpu"
	"github.com/StellarisJAY/gbgo/ppu"
	"github.com/veandco/go-sdl2/sdl"
	"io"
	"os"
	"time"
	"unsafe"
)

type config struct {
	file  string
	scale int
	fps   int64
}

type Emulator struct {
	conf     *config
	game     cartridge.BasicCartridge
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture

	cpu           *cpu.Processor
	ppu           *ppu.PPU
	bus           *bus.Bus
	lastFrameTime int64
}

func parseConfigs() *config {
	conf := &config{}
	flag.StringVar(&conf.file, "file", "", "game rom file")
	flag.IntVar(&conf.scale, "scale", 1, "window scale")
	flag.Int64Var(&conf.fps, "fps", 30, "frame rate")
	flag.Parse()
	if conf.fps < 20 {
		conf.fps = 20
	} else if conf.fps > 60 {
		conf.fps = 60
	}
	return conf
}

func readGbFile(fileName string) []byte {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("open gb file error %w", err))
	}
	defer file.Close()
	raw, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("read gb file error %w", err))
	}
	return raw
}

func initSDL(conf *config) (*sdl.Window, *sdl.Renderer, *sdl.Texture) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Errorf("init sdl error %w", err))
	}
	window, err := sdl.CreateWindow("gbgo", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(ppu.Width*conf.scale), int32(ppu.Height*conf.scale), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(fmt.Errorf("sdl create window error %w", err))
	}
	renderer, err := sdl.CreateRenderer(window, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Errorf("sdl create renderer error %w", err))
	}
	_ = renderer.SetScale(float32(conf.scale), float32(conf.scale))

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, ppu.Width, ppu.Height)
	if err != nil {
		panic(fmt.Errorf("sdl create texture error %w", err))
	}
	return window, renderer, texture
}

func MakeEmulator() *Emulator {
	conf := parseConfigs()
	raw := readGbFile(conf.file)
	c := cartridge.MakeBasicCartridge(raw)
	window, renderer, texture := initSDL(conf)
	b := bus.MakeBus(&c)
	processor := cpu.MakeCPU(b)
	return &Emulator{
		conf:     conf,
		game:     c,
		window:   window,
		renderer: renderer,
		texture:  texture,
		bus:      b,
		cpu:      processor,
		ppu:      ppu.MakePPU(),
	}
}

func (e *Emulator) start() {
	// 重置cpu各个寄存器
	e.cpu.Reset()
	interval := 1000 / e.conf.fps
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	e.lastFrameTime = time.Now().UnixMilli()
	for {
		select {
		case _ = <-ticker.C:
			e.Update()
		}
	}
}

// Update 每一帧更新一次，每次Update进行IO、执行CPU指令、渲染画面
func (e *Emulator) Update() {
	frameTime := time.Now().UnixMilli()
	// 输入事件处理
	e.handleEvents()
	// cpu tick
	e.cpu.Tick(frameTime, logInstruction)
	// 渲染画面
	e.renderFrame()
	e.lastFrameTime = frameTime
}

func (e *Emulator) renderFrame() {
	// ppu渲染tiles和sprites，生成frame数据
	e.ppu.Render()
	// frame数据渲染到屏幕
	frame := ppu.FrameData()
	_ = e.texture.Update(nil, unsafe.Pointer(&frame[0]), ppu.Width*3)
	_ = e.renderer.Copy(e.texture, nil, nil)
	e.renderer.Present()
}

func (e *Emulator) handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			ev := event.(*sdl.KeyboardEvent)
			switch ev.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				e.onShutdown()
				os.Exit(0)
				return
			case sdl.SCANCODE_W:
			case sdl.SCANCODE_S:
			case sdl.SCANCODE_A:
			case sdl.SCANCODE_D:
			}
		case *sdl.QuitEvent:
			e.onShutdown()
			os.Exit(0)
			return
		}
	}
}

func (e *Emulator) onShutdown() {
	_ = e.texture.Destroy()
	_ = e.renderer.Destroy()
	_ = e.window.Destroy()
}
