package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var onlyNumbersDotAndComma *regexp.Regexp

func init() {
	onlyNumbersDotAndComma = regexp.MustCompile(`-?\d+[.,]?\d+`)
}

func selectOnlyDigits(input string) string {
	digits := strings.Join(onlyNumbersDotAndComma.FindAllString(input, -1), "")
	if len(digits) == 0 {
		return ""
	}
	return digits
}

func parseFloat(input string) (float64, error) {
	withDot := strings.Replace(input, ",", ".", 1)
	return strconv.ParseFloat(withDot, 64)
}

func floatToStringWith2Digits(input float64) string {
	result := fmt.Sprintf("%.2f", input)
	return strings.Replace(result, ".", ",", 1)
}
