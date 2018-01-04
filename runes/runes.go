package runes

import (
	"fmt"

	"golang.org/x/text/width"
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

func IsWide(r rune) bool {
	k := width.LookupRune(r).Kind()

	switch k {
	case width.Neutral:
		fallthrough
	case width.EastAsianNarrow:
		fallthrough
	case width.EastAsianHalfwidth:
		return false
	default:
		return true
	}
}
