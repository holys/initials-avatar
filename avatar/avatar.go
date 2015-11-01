package avatar

import (
	"errors"
	"image/color"
	"strings"
	"unicode"

	"stathat.com/c/consistent"
)

var c = consistent.New()

var (
	AvatarBgColors = map[string]*color.RGBA{
		"45BDF3": &color.RGBA{69, 189, 243, 255},
		"E08F70": &color.RGBA{224, 143, 112, 255},
		"4DB6AC": &color.RGBA{77, 182, 172, 255},
		"9575CD": &color.RGBA{149, 117, 205, 255},
		"B0855E": &color.RGBA{176, 133, 94, 255},
		"F06292": &color.RGBA{240, 98, 146, 255},
		"A3D36C": &color.RGBA{163, 211, 108, 255},
		"7986CB": &color.RGBA{121, 134, 203, 255},
		"F1B91D": &color.RGBA{241, 185, 29, 255},
	}

	DefaultColorKey = "45BDF3"

	ErrUnsupportChar = errors.New("unsupport character")
)

func init() {
	for key := range AvatarBgColors {
		c.Add(key)
	}
}

type InitialsAvatar struct {
	Name     string
	Initials string
	Color    *color.RGBA
}

func NewInitialsAvatar(name string) (*InitialsAvatar, error) {
	name = strings.TrimSpace(name)
	firstRune := []rune(name)[0]
	if !isHan(firstRune) && !unicode.IsLetter(firstRune) {
		return nil, ErrUnsupportChar
	}

	avatar := new(InitialsAvatar)
	avatar.Name = name
	avatar.Initials = strings.ToUpper(getInitials(name))
	avatar.Color = getColorByName(name)

	return avatar, nil
}

func isHan(r rune) bool {
	if unicode.Is(unicode.Scripts["Han"], r) {
		return true
	}
	return false
}

func getColorByName(name string) *color.RGBA {
	key, err := c.Get(name)
	if err != nil {
		key = DefaultColorKey
	}
	return AvatarBgColors[key]
}

func getInitials(name string) string {
	if len(name) <= 0 {
		return ""
	}
	nameRunes := []rune(name)
	return string(nameRunes[0])
}
