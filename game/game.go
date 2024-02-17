package game

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Game struct {
	grid         grid
	charX, charY int
}

type grid struct {
	Width  int
	Height int
	Cells  [][]rune
}

func GetGrid() grid {
	g := grid{Width: 10, Height: 10}
	g.Cells = make([][]rune, g.Height)
	for i := range g.Cells {
		g.Cells[i] = make([]rune, g.Width)
	}
	return g
}

func (g *grid) String() string {
	var s string
	for _, row := range g.Cells {
		for _, cell := range row {
			if cell == 0 {
				cell = '.'
			}
			s += string(cell)
		}
		s += "\n"
	}
	return s
}

func (g *grid) Set(x, y int, r rune) {
	g.Cells[y][x] = r
}

func Start() {
	g := Game{grid: GetGrid(), charX: 0, charY: 0}
	fps := 30
	frameTime := time.Duration(1000/fps) * time.Millisecond

	for {
		lastTime := time.Now()
		g.draw()
		g.update()
		time.Sleep(frameTime - time.Since(lastTime))
	}
}

func (g *Game) clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (g *Game) draw() {
	g.clearScreen()
	fmt.Print(g.grid.String())
}

func (g *Game) update() {
	g.grid.Set(g.charX, g.charY, '.')
	if g.charX == g.grid.Width-1 {
		g.charX = 0
		g.charY++
		if g.charY == g.grid.Height-1 {
			g.charY = 0
		}
	}
	g.charX++
	g.grid.Set(g.charX, g.charY, 'X')
}
