package common

import (
	"encoding/csv"
	"golang.org/x/text/encoding"
	"math/rand"
	"strings"
)

func SkipSomeLines(reader *csv.Reader, magnitude float32) {
	var skip = (rand.Int()%10 + 1) * int(magnitude*(1+rand.Float32()))
	for j := 0; j < skip; j++ {
		_, _ = reader.Read()
		continue
	}
}

func EncodeToCsv(record []string, separator rune, encoder *encoding.Encoder) string {
	row := "\"" + strings.Join(record, "\""+string(separator)+"\"") + "\"\n"
	encoded, err := encoder.String(row)
	if err != nil {
		panic(err)
	}
	return encoded
}
