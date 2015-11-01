package avatar

import (
	"image"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Generator  generates image.Image
type Generator struct {
	sideLength  int
	fontFile    string
	dpi         float64
	fontSize    float64
	fontSpacing float64
	fontHinting font.Hinting
	face        font.Face
}

func NewGenerator(fontFile string) *Generator {
	if fontFile == "" {
		panic("fontFile must be specific font path")
	}
	g := new(Generator)
	//TODO: configurable
	g.sideLength = 120
	g.fontFile = fontFile
	g.dpi = 72.0
	g.fontSize = 75.0
	g.fontSpacing = 1.5 //不需要字符间空隙
	g.fontHinting = font.HintingNone

	ttFont, err := getTrueTypeFont(g.fontFile)
	if err != nil {
		panic(err)
	}

	g.face = truetype.NewFace(ttFont, &truetype.Options{
		Size:    g.fontSize,
		DPI:     g.dpi,
		Hinting: g.fontHinting,
	})

	return g
}

func (g *Generator) Generate(avatar *InitialsAvatar) image.Image {
	width := g.sideLength
	height := g.sideLength
	fontSize := g.fontSize
	dpi := g.dpi
	face := g.face

	bgColor := avatar.Color
	charToDraw := avatar.Initials

	dstImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dstImage, dstImage.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)
	drawer := &font.Drawer{
		Dst:  dstImage,
		Src:  image.White,
		Face: face,
	}

	y := 10 + int(math.Ceil(fontSize*dpi/72))
	drawer.Dot = fixed.Point26_6{
		X: (fixed.I(width) - drawer.MeasureString(charToDraw)) / 2,
		Y: fixed.I(y),
	}

	drawer.DrawString(charToDraw)

	return dstImage
}

func getTrueTypeFont(fontFile string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}

	ttf, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return ttf, nil
}
