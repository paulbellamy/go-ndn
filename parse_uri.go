package ndn

import (
	"errors"
	"strconv"
	"strings"
)

func ParseURI(rawURI string) (name Name, err error) {
	scheme, rest, err := getScheme(rawURI)
	if scheme != "ndn" {
		return nil, errors.New("protocol scheme must be ndn")
	}

	var segments []string
	if !strings.HasPrefix(rest, "/") {
		return nil, errors.New("invalid uri")
	} else if strings.HasPrefix(rest, "//") {
		segments = strings.Split(rest[2:], "/")[1:]
	} else {
		segments = strings.Split(rest[1:], "/")
	}

	for _, segment := range segments {
		unescaped, err := unescapeHex(segment)
		if err != nil {
			return nil, err
		}
		unescaped, err = unescapeDots(unescaped)
		if err != nil {
			return nil, err
		}

		name = name.AppendComponent(Component{unescaped})
	}

	return name, nil
}

func getScheme(rawURI string) (scheme, rest string, err error) {
	for i := 0; i < len(rawURI); i++ {
		c := rawURI[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
			// do nothing
		case '0' <= c && c <= '9' || c == '+' || c == '-' || c == '.':
			if i == 0 {
				return "", rawURI, nil
			}
		case c == ':':
			if i == 0 {
				return "", "", errors.New("missing protocol scheme")
			}
			return rawURI[0:i], rawURI[i+1:], nil
		default:
			// we have encountered an invalid character,
			// so there is no valid scheme
			return "", rawURI, nil
		}
	}
	return "", rawURI, nil
}

// unescape unescapes a string
func unescapeHex(s string) (string, error) {
	// Count %, check that they're well-formed.
	n := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[0:3]
				}
				return "", EscapeError(s)
			}
			i += 3
		default:
			i++
		}
	}

	if n == 0 {
		return s, nil
	}

	t := make([]byte, len(s)-2*n)
	j := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			t[j] = unhex(s[i+1])<<4 | unhex(s[i+2])
			j++
			i += 3
		case '+':
			t[j] = '+'
			j++
			i++
		default:
			t[j] = s[i]
			j++
			i++
		}
	}
	return string(t), nil
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func unescapeDots(s string) (string, error) {
	// ensure all are dots
	leadingDotCount := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '.' {
			return s, nil
		}
		leadingDotCount++
	}

	if leadingDotCount < 3 {
		return "", EscapeError(s)
	}

	return s[3:], nil
}

type EscapeError string

func (e EscapeError) Error() string {
	return "invalid URI escape " + strconv.Quote(string(e))
}
