package main

import (
	"flag"
	"fmt"
	"github.com/StellarisJAY/gbgo/cartridge"
	"io"
	"os"
)

type config struct {
	file string
}

type Emulator struct {
	game cartridge.BasicCartridge
}

func parseConfigs() *config {
	conf := &config{}
	flag.StringVar(&conf.file, "file", "", "game rom file")
	flag.Parse()
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

func MakeEmulator() *Emulator {
	conf := parseConfigs()
	raw := readGbFile(conf.file)
	c := cartridge.MakeBasicCartridge(raw)
	return &Emulator{game: c}
}
