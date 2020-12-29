package starter

import (
	"errors"
	"strings"
)

func getHostPort(addr string) (h string, p string) {
	s := strings.Split(addr, ":")
	if len(s) == 2 {
		return s[0], s[1]
	} else if len(s) == 3 {
		return strings.TrimLeft(s[1], "////"), s[2]
	}
	panic(errors.New("wrong addr format"))
}
