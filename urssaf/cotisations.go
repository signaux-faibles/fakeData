package urssaf

import (
	"bufio"
	"encoding/csv"
	"github.com/signaux-faibles/fakeData/common"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func ReadAndRandomCotisations(source string, outputFileName string, outputSize int, mapping map[string]string) error {
	// source
	sourceFile, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("error when closing", file)
		}
	}(sourceFile)
	reader := csv.NewReader(bufio.NewReader(sourceFile))
	reader.Comma = ';'

	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("error when closing", file)
		}
	}(outputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	writer.Comma = reader.Comma
	coef := make(map[string]float64)

	// ligne de titre
	row, err := reader.Read()
	if err != nil {
		return err
	}
	err = writer.Write(row)
	if err != nil {
		return err
	}
	wrote := 0
	for wrote <= outputSize {
		row, err := reader.Read()
		if err == io.EOF {
			log.Default().Println("EOF", source)
			break
		}
		row[2] = mapping[row[2]]
		if len(row[2]) == 0 { // mapping not found
			continue // skip line
		}
		mer, _ := strconv.ParseFloat(row[4], 64)
		encDirect, _ := strconv.ParseFloat(row[5], 64)
		cotisDue, _ := strconv.ParseFloat(row[6], 64)

		if c, ok := coef[row[2]]; ok {
			row[4] = strconv.Itoa(int(mer * c))
			row[5] = strconv.Itoa(int(encDirect * c))
			row[6] = strconv.Itoa(int(cotisDue * c))
		} else {
			coef[row[2]] = rand.Float64() * rand.Float64() / 150
			row[4] = strconv.Itoa(int(mer * coef[row[2]]))
			row[5] = strconv.Itoa(int(encDirect * coef[row[2]]))
			row[6] = strconv.Itoa(int(cotisDue * coef[row[2]]))
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
		if outputSize > 100 {
			if mod := wrote % (outputSize / 100); mod == 0 {
				log.Default().Println("(cotisations) wrote ", wrote/(outputSize/100), "%")
				common.SkipSomeLines(reader, 9.87)
			}
		}
		wrote++
	}

	return nil
}
