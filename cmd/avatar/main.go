package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"path/filepath"

	"github.com/holys/initials-avatar/avatar"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var (
	fontFile = flag.String("fontFile", "./resource/fonts/Hiragino Sans GB W3.ttf", "tty font file path")
	port     = flag.Int("port", 3000, "http port to run")
)

type avatarHandler struct {
	fontFile string
}

func newAvatarHandler(fontFile string) *avatarHandler {
	h := new(avatarHandler)
	h.fontFile = fontFile
	return h
}

func (h *avatarHandler) Get(ctx *echo.Context) error {
	name := ctx.Param("name")

	a, err := avatar.NewInitialsAvatar(name)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	g := avatar.NewGenerator(h.fontFile)
	m := g.Generate(a)

	ctx.Response().Header().Set("Content-Type", "image/png")
	ctx.Response().Header().Set("Cache-Control", "max-age=600")
	ctx.Response().WriteHeader(http.StatusOK)

	return png.Encode(ctx.Response().Writer(), m)
}

func main() {
	flag.Parse()
	if len(*fontFile) == 0 {
		log.Fatal("invalid font file path")
	}
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	fontFile, err := filepath.Abs(*fontFile)
	if err != nil {
		log.Fatal("invalid font file path")
	}
	h := newAvatarHandler(fontFile)
	e.Get("/avatar/:name", h.Get)

	fmt.Printf("starting at :%d ...\n", *port)
	e.Run(fmt.Sprintf(":%d", *port))

}
