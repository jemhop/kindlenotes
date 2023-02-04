package main

import (
	"errors"
	"strings"
)

// thanks to Hasan Yousef (https://forum.golangbridge.org/t/removing-first-and-last-empty-lines-from-a-string/24285)
func trimEmptyLines(b []byte) string {
	strs := strings.Split(string(b), "\n")
	str := ""
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		str += s + "\n"
	}
	str = strings.TrimSuffix(str, "\n")

	return str
}

// Thanks to peterSO (https://codereview.stackexchange.com/questions/122831/parse-numerals-from-a-string-in-golang)
var ErrRuneNotInt = errors.New("type: rune was not int")

func CharToNum(r rune) (int, error) {
	if '0' <= r && r <= '9' {
		return int(r) - '0', nil
	}
	return 0, ErrRuneNotInt
}

func ParseNum(s string) []int {
	nLen := 0
	for i := 0; i < len(s); i++ {
		if b := s[i]; '0' <= b && b <= '9' {
			nLen++
		}
	}
	var n = make([]int, 0, nLen)
	for i := 0; i < len(s); i++ {
		if b := s[i]; '0' <= b && b <= '9' {
			n = append(n, int(b)-'0')
		}
	}
	return n
}
