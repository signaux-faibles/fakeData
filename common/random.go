package common

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"hash/crc32"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var crc32q *crc32.Table

func init() {
	// In this package, the CRC polynomial is represented in reversed notation,
	// or LSB-first representation.
	//
	// LSB-first representation is a hexadecimal number with n bits, in which the
	// most significant bit represents the coefficient of x⁰ and the least significant
	// bit represents the coefficient of xⁿ⁻¹ (the coefficient for xⁿ is implicit).
	//
	// For example, CRC32-Q, as defined by the following polynomial,
	//	x³²+ x³¹+ x²⁴+ x²²+ x¹⁶+ x¹⁴+ x⁸+ x⁷+ x⁵+ x³+ x¹+ x⁰
	// has the reversed notation 0b11010101100000101000001010000001, so the value
	// that should be passed to MakeTable is 0xD5828281.
	crc32q = crc32.MakeTable(0xD5828281)
}

// FalsifyNumber no collision when `len(input) >= 10`
func FalsifyNumber(input string) (string, error) {
	length := len(input)
	checksum := crc32.Checksum([]byte(input), crc32q)
	sizedChecksum := fmt.Sprintf("%"+strconv.Itoa(length)+"v", checksum)
	if len(sizedChecksum) > length {
		result := sizedChecksum[0:length]
		return result, nil
	}
	result := strings.ReplaceAll(sizedChecksum, " ", "0")
	return result, nil
}

func RandStringBytesRmndr(n int) string {
	min := math.Pow(10, float64(n-1))
	max := math.Pow(10, float64(n))
	number := randomdata.Number(int(min), int(max))
	return strconv.Itoa(number)
}

func RandName() string {
	return strings.ToUpper(randomdata.SillyName())
}

func RandRaisonSociale(input string) string {
	if len(input) == 0 {
		return input
	}
	names := []string{RandName(), RandName()}
	return strings.Join(names, " ")
}

func RandCoeff() float64 {
	min := 0.8
	max := 1.2
	for {
		coeff := randomdata.Decimal(2)
		if coeff > min && coeff < max {
			return coeff
		}
	}
}

func RandDate(format string, toChange string) (string, error) {
	if len(toChange) == 0 {
		return toChange, nil
	}
	parse, err := time.Parse(format, toChange)
	if err != nil {
		log.Default().Println("can't parse date from", toChange, "with format", format, ". Error is", err)
		return "", err
	}
	randomized := randDateAround(parse)
	return randomized.Format(format), nil
}

func RandDecimal(input string) (string, error) {
	digits := selectOnlyDigits(input)
	if len(digits) == 0 {
		return "", nil
	}
	val, err := parseFloat(digits)
	if err != nil {
		return input, err
	}
	around := val * RandCoeff()
	//if around == 0 {
	//	return "", nil
	//}
	return floatToStringWith2Digits(around), nil
}

func RandInt(input string) (string, error) {
	toInt := selectOnlyDigits(input)
	if len(toInt) == 0 {
		return "", nil
	}
	val, err := strconv.Atoi(toInt)
	if err != nil {
		return input, err
	}
	around := int(float64(val) * RandCoeff())
	//if around == 0 {
	//	return "", nil
	//}
	return strconv.Itoa(around), nil
}

func RandItemFrom(datas []string) string {
	size := len(datas)
	if size > 0 {
		i := rand.Int() % size
		return datas[i]
	}
	log.Default().Println("can't get an item from empty arrays")
	return ""
}

func randDateAround(toChange time.Time) time.Time {
	var today = time.Now()
	if toChange.After(today) {
		return changeDate(toChange)
	}
	return changeDateButLimited(toChange, today)
}

func changeDate(toChange time.Time) time.Time {
	daysToChange := randomdata.Number(-19, 19)
	r := toChange.AddDate(0, 0, daysToChange)
	return r
}

func changeDateButLimited(toChange time.Time, limit time.Time) time.Time {
	r := limit
	for !r.Before(limit) {
		daysToChange := randomdata.Number(-19, 19)
		r = toChange.AddDate(0, 0, daysToChange)
	}
	return r
}
