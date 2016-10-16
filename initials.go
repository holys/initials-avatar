package avatar

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
)

const email = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

var regxEmail *regexp.Regexp

func init() {
	regxEmail = regexp.MustCompile(email)
}

// This controls how parsing of initials is handled.
type opts struct {
	// set all initials to capital letters
	allCaps bool

	// allow parsing of initials from email messages
	allowEmail bool

	// set the maximum number of initials allowed
	limit int
}

// Tries to find initials in a given src. The src is a name, the logic that is
// used to decide which characters are used as initials is adopted from the
// initials project https://github.com/gr2m/initials.
//
// You can pass an opts object to contorl the parsing like setting maximum
// number of initials and allowing parsing of initials from emails etc.
func parseInitials(src io.Reader, o opts) (string, error) {
	scanner := bufio.NewScanner(src)
	scanner.Split(bufio.ScanWords)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	count := 0
	for i, w := range words {
		if count >= o.limit {
			break
		}
		if regxEmail.MatchString(w) {
			if i == 0 && o.allowEmail {
				s := strings.Split(w, "@")
				sr := strings.NewReader(s[0])
				return parseInitials(sr, o)
			}
			continue
		}
		r := strings.NewReader(w)
		x, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		switch {
		case unicode.IsLetter(x):
			if o.allCaps {
				if unicode.IsLower(x) {
					x = unicode.ToUpper(x)
				}
			}
			_, _ = buf.WriteRune(x)
			count++
		case x == '(':
			if i > 0 {
				rb := &bytes.Buffer{}
			DONE:
				for {
					next, _, err := r.ReadRune()
					switch next {
					case ')':
						_, _, err = r.ReadRune()
						if err != nil {
							if err != io.EOF {
								rb.Reset()
							}
						}
						break DONE
					default:
						if err != nil {
							rb.Reset()
							break DONE
						}
						_, _ = rb.WriteRune(next)
					}

				}
				return rb.String(), nil
			}
		}

	}
	return buf.String(), nil
}
