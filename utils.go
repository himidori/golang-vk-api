package vkapi

import "strconv"

func ArrayToStr(a []int) string {
	var s string
	for _, num := range a {
		s += strconv.Itoa(num) + ","
	}

	return s[:len(s)-1]
}
