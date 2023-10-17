package drummer

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const BACKGROUND_IMG = "assets/background.png"

type Display struct {
	backgroundImg *ebiten.Image
	regularFont   font.Face
	bigFont       font.Face
}

func NewDisplay() *Display {
	display := &Display{}

	// load background.png and draw it to the screen
	img, _, err := ebitenutil.NewImageFromFile(BACKGROUND_IMG)
	if err != nil {
		log.Fatal(err)
	}
	display.backgroundImg = img

	// set up fonts
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)

	const dpi = 72
	display.regularFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	display.bigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	return display
}

func (d *Display) Draw(screen *ebiten.Image, dm *DrumMachine) {
	screen.DrawImage(d.backgroundImg, nil)

	// draw the current name of the pattern
	text.Draw(screen, dm._current_pattern.Name, d.regularFont, 10, 180, color.White)

	// draw the current index of the pattern
	text.Draw(screen, fmt.Sprint(dm.Current_pattern_idx), d.bigFont, 38, 60, color.Black)

	// draw the current tempo
	text.Draw(screen, fmt.Sprint(dm.Tempo), d.bigFont, 177, 40, color.White)

	// draw the beat-indicating dot
	if dm.Beat%4 == 0 && dm.State == "playing" {
		//	pygame.draw.circle(self.screen, "yellow", (self.SCREEN_WIDTH/2, 40), 4)
		vector.DrawFilledCircle(screen, 120, 40, 4, color.White, true)
	}
	// draw the beat tracker
	for i := 0; i < dm.Pattern_length; i++ {
		const (
			bar_width      = 10.0
			bar_height     = 18.0
			bar_spacing    = 4.0
			bar_y_position = 240.0 - bar_height - 10.0
		)

		vector.StrokeRect(screen,
			10.0+float32(i)*(bar_width+bar_spacing), bar_y_position,
			bar_width, bar_height,
			2, color.Black, true)
		if dm.Beat == i {
			vector.DrawFilledRect(screen,
				10.0+float32(i)*(bar_width+bar_spacing),
				bar_y_position, bar_width, bar_height,
				color.Black, true)
		}
	}
}
