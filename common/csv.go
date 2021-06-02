package common

import (
	"encoding/csv"
	"math/rand"
)

func SkipSomeLines(reader *csv.Reader, magnitude float32) {
	var skip = (rand.Int()%10 + 1) * int(magnitude*(1+rand.Float32()))
	for j := 0; j < skip; j++ {
		_, _ = reader.Read()
		continue
	}
}
