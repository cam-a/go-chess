// Copyright 2016 The Ebiten Authors

package chess

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	ScreenWidth  = 900
	ScreenHeight = 900
	boardSize    = 18
	sampleRate   = 44100
)

var (
	shadowColor = color.NRGBA{0, 0, 0, 0x80}
)

// Game represents a game state.
type Game struct {
	input      *Input
	board      *Board
	boardImage *ebiten.Image
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{
		input: NewInput(),
	}
	var err error
	sounds, err := newAudio()
	if err != nil {
		return nil, err
	}
	g.board, err = NewBoard(sounds, boardSize)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.input.Update()
	if err := g.board.Update(g.input); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		w, h := g.board.Size()
		g.boardImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)

	text.Draw(screen, "Points: "+strconv.Itoa(g.board.GetPoints()), getArcadeFonts(3), 10, 30, color.White)

	if g.board.playMode == false {
		screen.Fill(backgroundColor)
		text.Draw(screen, "you died", getArcadeFonts(3), 10, 30, shadowColor)
		text.Draw(screen, "time to restart from terminal! :)", getArcadeFonts(3), 10, 70, shadowColor)
		text.Draw(screen, "your score: "+strconv.Itoa(g.board.GetPoints()), getArcadeFonts(3), 10, 120, shadowColor)
	}
}

func getArcadeFonts(scale int) font.Face {
	if arcadeFonts == nil {
		tt, err := truetype.Parse(fonts.ArcadeN_ttf)
		if err != nil {
			log.Fatal(err)
		}

		arcadeFonts = map[int]font.Face{}
		for i := 1; i <= 4; i++ {
			const dpi = 72
			arcadeFonts[i] = truetype.NewFace(tt, &truetype.Options{
				Size:    float64(arcadeFontBaseSize * i),
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
		}
	}
	return arcadeFonts[scale]
}
