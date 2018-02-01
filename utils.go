package vkapi

import (
	"strconv"
	"strings"
)

func ArrayToStr(a []int) string {
	s := []string{}
	for _, num := range a {
		s = append(s, strconv.Itoa(num))
	}
	return strings.Join(s, ",")
}
