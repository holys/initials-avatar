package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/holys/initials-avatar/avatar"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var (
	fontFile = flag.String("fontFile", "./resource/fonts/Hiragino Sans GB W3.ttf", "tty font file path")
	port     = flag.Int("port", 3000, "http port to run")
)

type avatarHandler struct {
	avatar *avatar.InitialsAvatar
}

func newAvatarHandler(fontFile string) *avatarHandler {
	h := new(avatarHandler)
	h.avatar = avatar.New(fontFile)
	return h
}

func (h *avatarHandler) Get(ctx *echo.Context) error {
	name := ctx.Param("name")
	size := ctx.Query("size")
	if size == "" {
		size = "120"
	}
	//FIXME: auto size
	sz, err := strconv.Atoi(size)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data, err := h.avatar.DrawToBytes(name, sz)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx.Response().Header().Set("Content-Type", "image/png")
	ctx.Response().Header().Set("Cache-Control", "max-age=600")
	ctx.Response().WriteHeader(http.StatusOK)
	ctx.Response().Write(data)

	return nil
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
	e.Get("/:name", h.Get)

	fmt.Printf("starting at :%d ...\n", *port)
	e.Run(fmt.Sprintf(":%d", *port))

}
