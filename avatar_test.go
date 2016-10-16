package avatar

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitialsAvatar_DrawToBytes(t *testing.T) {
	fontFile := os.Getenv("AVATAR_FONT")
	if fontFile == "" {
		t.Skip("Font file is needed")
	}

	av := New(fontFile)

	stuffs := []struct {
		name     string
		size     int
		encoding string
	}{
		{"Swordsmen", 22, "png"},
		{"Condor Heroes", 30, "jpeg"},
		{"孔子", 22, "png"},
		{"Swordsmen", 0, "png"},
		{"*", 22, "png"},
	}

	for _, v := range stuffs {
		raw, err := av.DrawToBytes(v.name, v.size, v.encoding)
		if err != nil {
			if err == ErrUnsupportChar {
				t.Skip("ErrUnsupportChar")
			}
			t.Error(err)
		}
		switch v.encoding {
		case "png":
			if _, perr := png.Decode(bytes.NewReader(raw)); perr != nil {
				t.Error(perr)
			}
		case "jpeg":
			if _, perr := jpeg.Decode(bytes.NewReader(raw)); perr != nil {
				t.Error(perr, v)
			}
		}
	}
}

func TestGetInitials(t *testing.T) {
	names := []struct {
		full, intitials string
	}{
		{"John", "J"},
		{"Doe", "D"},
		{"", ""},
		{"John Doe", "JD"},
		{"john doe", "jd"},
		{"joe@example.com", "j"},
		{"John Doe (dj)", "dj"},
	}

	for _, v := range names {
		n := getInitials(v.full)
		if n != v.intitials {
			t.Errorf("expected %s got %s", v.intitials, n)
		}
	}
}

func TestParseFont(t *testing.T) {
	fileNotExists := "xxxxxxx.ttf"
	_, err := parseFont(fileNotExists)
	if err == nil {
		t.Error("should return error")
	}

	_, err = newDrawer(fileNotExists)
	if err == nil {
		t.Error("should return error")
	}

	fileExistsButNotTTF, _ := ioutil.TempFile(os.TempDir(), "prefix")
	defer os.Remove(fileExistsButNotTTF.Name())

	_, err = parseFont(fileExistsButNotTTF.Name())
	if err == nil {
		t.Error("should return error")
	}
	_, err = newDrawer(fileExistsButNotTTF.Name())
	if err == nil {
		t.Error("should return error")
	}
	_, err = newDrawer("")
	if err == nil {
		t.Error("should return error")
	}

}
