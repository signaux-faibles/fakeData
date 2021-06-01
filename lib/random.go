package lib

import "math/rand"

const numberBytes = "0123456789"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numberBytes[rand.Int63()%int64(len(numberBytes))]
	}
	return string(b)
}

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ "

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func RandBool() bool {
	return rand.Float32() > 0.5
}
