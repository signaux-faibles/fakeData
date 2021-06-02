package common

import (
	"github.com/Pallinder/go-randomdata"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var onlyDigitCharacter *regexp.Regexp

func init() {
	onlyDigitCharacter = regexp.MustCompile(`\w`)
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

func RandAround(input int) int {
	return int(float64(input) * RandCoeff())
}

func RandEffectif(input string) (string, error) {
	toInt := strings.Join(onlyDigitCharacter.FindAllString(input, -1), "")
	if len(toInt) == 0 {
		return "", nil
	}
	val, err := strconv.Atoi(toInt)
	if err != nil {
		return input, err
	}
	around := RandAround(val)
	if around == 0 {
		return "", nil
	}
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
