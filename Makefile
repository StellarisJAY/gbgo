export TARGET="./gbgo"
export ROM = "roms/games/Tetris.gb"
export FPS = 30
build:
	@go build -o $(TARGET)
run:build
	@$(TARGET) -file $(ROM) -fps $(FPS) -scale 3
trace:build
	@$(TARGET) -file $(ROM) -fps 60 -trace
