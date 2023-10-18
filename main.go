package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/tgerla/camina/drummer"
)

type Game struct {
	drumMachine *drummer.DrumMachine
	display     *drummer.Display
}

func (g *Game) Update() error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		{
			if g.drumMachine.State == "stopped" {
				g.drumMachine.Play()
			} else {
				g.drumMachine.Stop()
			}
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		{
			if g.drumMachine.CurrentPatternIndex > 1 {
				g.drumMachine.SwitchPattern(g.drumMachine.CurrentPatternIndex - 1)
			} else {
				g.drumMachine.SwitchPattern(len(g.drumMachine.Patterns))
			}
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		{
			if g.drumMachine.CurrentPatternIndex < len(g.drumMachine.Patterns) {
				g.drumMachine.SwitchPattern(g.drumMachine.CurrentPatternIndex + 1)
			} else {
				g.drumMachine.SwitchPattern(1)
			}
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		{
			if g.drumMachine.Tempo < 300 {
				g.drumMachine.SetTempo(g.drumMachine.Tempo + 1)
			}
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		{
			if g.drumMachine.Tempo > 30 {
				g.drumMachine.SetTempo(g.drumMachine.Tempo - 1)
			}
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
	case inpututil.IsKeyJustPressed(ebiten.KeyQ):
		{
			g.drumMachine.Stop()
			os.Exit(0)
		}
	}

	g.drumMachine.Tick()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.display.Draw(screen, g.drumMachine)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 240, 240
}

func main() {
	dm := drummer.NewDrumMachine()
	dm.Play()

	display := drummer.NewDisplay()

	game := &Game{
		drumMachine: dm,
		display:     display,
	}

	ebiten.SetWindowSize(240, 240)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
