package runes

import (
	"fmt"
)

func ImageName(r rune) string {
	return fmt.Sprintf("image/rune/u%s.png", RuneCode(r))
}

func GetFirstRune(str string) rune {
	rs := []rune(str)
	return rs[0]
}

func RuneCode(r rune) string {
	return fmt.Sprintf("%X", r)
}
