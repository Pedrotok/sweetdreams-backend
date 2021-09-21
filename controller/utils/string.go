package utils

import (
	"regexp"
	"strconv"
)

func ParsePriceFromStringToInt(price string) int {
	re, _ := regexp.Compile(`[^0-9]`)
	cents := re.ReplaceAllString(price, "")
	centsInInt, _ := strconv.Atoi(cents)
	return centsInInt
}
