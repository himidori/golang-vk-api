package vkapi

import "strings"

func ArrayToStr(a []int) string {
	return strings.Join(a, ",")
}
