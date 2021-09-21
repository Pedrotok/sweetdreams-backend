package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ParsePriceFromStringToInt(price string) int {
	re, _ := regexp.Compile(`[^0-9]`)
	cents := re.ReplaceAllString(price, "")
	fmt.Println(cents)
	centsInInt, _ := strconv.Atoi(cents)
	return centsInInt
}
