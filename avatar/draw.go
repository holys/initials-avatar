package avatar

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	//"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	errFontRequired = errors.New("font file is required")
	errInvalidTTF   = errors.New("invalid ttf")
)

// drawer draws an image.Image
type drawer struct {
	fontSize    float64
	dpi         float64
	fontHinting font.Hinting
	face        font.Face
	font        *truetype.Font
}

func newDrawer(fontFile string) (*drawer, error) {
	if fontFile == "" {
		return nil, errFontRequired
	}
	g := new(drawer)
	g.fontSize = 75.0
	g.dpi = 72.0
	g.fontHinting = font.HintingNone

	ttf, err := getTTF(fontFile)
	if err != nil {
		return nil, errInvalidTTF
	}
	g.face = truetype.NewFace(ttf, &truetype.Options{
		Size:    g.fontSize,
		DPI:     g.dpi,
		Hinting: g.fontHinting,
	})

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, errInvalidTTF
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, errInvalidTTF
	}
	g.font = font
	return g, nil
}

// our avatar image is square
func (g *drawer) Draw(s string, size int, bg *color.RGBA) image.Image {
	// draw the background
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)

	// draw the text
	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: g.face,
	}

	// font index
	fi := g.font.Index([]rune(s)[0])

	// glyph example: http://www.freetype.org/freetype2/docs/tutorial/metrics.png
	var gbuf truetype.GlyphBuf
	var err error
	var _fsize fixed.Int26_6 = fixed.Int26_6(g.fontSize * g.dpi * (64.0 / 72.0))
	err = gbuf.Load(g.font, _fsize, fi, font.HintingFull)
	if err != nil {
		// fixme
		drawer.DrawString("")
		return dst
	}

	// center
	dY := int((size - int(gbuf.Bounds.Max.Y-gbuf.Bounds.Min.Y)>>6) / 2)
	dX := int((size - int(gbuf.Bounds.Max.X-gbuf.Bounds.Min.X)>>6) / 2)
	y := int(gbuf.Bounds.Max.Y>>6) + dY
	x := 0 - int(gbuf.Bounds.Min.X>>6) + dX

	//y := 10 + int(math.Ceil(g.fontSize*g.dpi/72)) //FIXME: what does it mean?
	drawer.Dot = fixed.Point26_6{
		X: fixed.I(x), //(fixed.I(size) - drawer.MeasureString(s)) / 2,
		Y: fixed.I(y),
	}
	drawer.DrawString(s)

	return dst
}

// read the font file as *truetype.Font
func getTTF(fontFile string) (*truetype.Font, error) {
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
