package avatar

import (
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
		undersize bool
		oversize  bool
	}{
		{"Swordsmen", 22, true, false},
		{"Condor Heroes", 30, false, false},
		{"Condor Heroes", 30, false, false},
		//		{"Condor Heroes", 200, false, true},
	}

	for _, v := range stuffs {
		_, err := av.DrawToBytes(v.name, v.size)
		if err != nil {
			t.Error(err)
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
