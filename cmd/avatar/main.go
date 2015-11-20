package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/holys/initials-avatar"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
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

func server(ctx *cli.Context) {
	fontFile := ctx.String("fontFile")
	// port := ctx.Int("port")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	fFile, err := filepath.Abs(fontFile)
	if err != nil {
		log.Fatal("invalid font file path")
	}
	h := newAvatarHandler(fFile)
	e.Get("/:name", h.Get)

	fmt.Printf("starting at :%s ...\n", port)
	e.Run(fmt.Sprintf(":%s", port))
}

func serverCommand() cli.Command {
	return cli.Command{
		Name:      "server",
		ShortName: "s",
		Usage:     "runs the webserver",
		Action:    server,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "fontFile",
				Usage: "tty font file path",
				Value: "./resource/fonts/Hiragino_Sans_GB_W3.ttf",
			},
			cli.IntFlag{
				Name:  "port",
				Usage: "http port to run",
				Value: 3000,
			},
		},
	}
}
func main() {
	a := cli.NewApp()
	a.Name = "Initials-avatar"
	a.Version = "0.0.1"
	a.Usage = "Generate an avatar image from a user's initials"
	a.Authors = []cli.Author{
		{"holys", "chendahui007@gmail.com"},
	}
	a.Commands = []cli.Command{
		serverCommand(),
	}
	a.RunAndExitOnError()
}
