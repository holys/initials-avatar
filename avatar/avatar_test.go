package avatar

import (
	"bytes"
	"image/jpeg"
	"image/png"
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
		name      string
		size      int
		encoding  string
		undersize bool
		oversize  bool
	}{
		{"Swordsmen", 22, "png", true, false},
		{"Condor Heroes", 30, "jpeg", false, false},
		//		{"Condor Heroes", 200, false, true},
	}

	for _, v := range stuffs {
		raw, err := av.DrawToBytes(v.name, v.size, v.encoding)
		if err != nil {
			t.Error(err)
		}
		switch v.encoding {
		case "png":
			if _, perr := png.Decode(bytes.NewReader(raw)); perr != nil {
				t.Error(perr)
			}
		case "jpeg":
			if _, perr := jpeg.Decode(bytes.NewReader(raw)); perr != nil {
				t.Error(perr)
			}
		}
	}
}

func TestGetInitials(t *testing.T) {
	names := []struct {
		full, intitials string
	}{
		{"David", "D"},
		{"Goliath", "G"},
		//		{"David Goliath", "DG"},
	}

	for _, v := range names {
		n := getInitials(v.full)
		if n != v.intitials {
			t.Errorf("expected %s got %s", v.intitials, n)
		}
	}
}
