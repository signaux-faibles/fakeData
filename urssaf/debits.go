package urssaf

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func ReadAndRandomDebits(source string, outputFileName string, outputSize int, mapping map[string]string) error {
	// source
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
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
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			log.Fatalln("error when flushing writer on", outputFile)
		}
	}(outputFile)
	writer := csv.NewWriter(outputFile)
	writer.Comma = reader.Comma
	defer writer.Flush()

	// ligne de titre
	row, err := reader.Read()
	err = writer.Write(row)
	if err != nil {
		return err
	}

	// map des coefficients générés
	coef := make(map[string]float64)

	wrote := 0
	for wrote <= outputSize {
		row, err := reader.Read()
		if err == io.EOF {
			log.Default().Println("EOF", source)
			break
		} else if err != nil {
			return err
		}

		row[2] = mapping[row[2]]
		if len(row[2]) == 0 { // mapping not found
			continue // skip line
		}

		partOuvriere, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			return err
		}
		partPatronale, err := strconv.ParseFloat(row[12], 64)
		if err != nil {
			return err
		}
		// est-ce le bon nom ???
		partPenalite, err := strconv.ParseFloat(row[15], 64)
		if err != nil {
			return err
		}

		if c, ok := coef[row[2]]; ok {
			row[11] = strconv.Itoa(int(partOuvriere * c))
			row[12] = strconv.Itoa(int(partPatronale * c))
			row[15] = strconv.Itoa(int(partPenalite * c))
		} else {
			coef[row[2]] = rand.Float64() * rand.Float64() / 150
			row[11] = strconv.Itoa(int(partOuvriere * coef[row[2]]))
			row[12] = strconv.Itoa(int(partPatronale * coef[row[2]]))
			row[15] = strconv.Itoa(int(partPenalite * coef[row[2]]))
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
		if outputSize > 100 {
			if mod := wrote % (outputSize / 100); mod == 0 {
				log.Default().Println("(debits) wrote ", wrote/(outputSize/100), "%")
				skipSomeDebits(reader)
			}
		}
		wrote++
	}
	return nil
}

func skipSomeDebits(reader *csv.Reader) {
	var skip = rand.Int() % 10 * 10
	for j := 0; j < skip; j++ {
		_, _ = reader.Read()
		continue
	}
}
