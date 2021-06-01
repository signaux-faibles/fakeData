package urssaf

import (
	"bufio"
	"encoding/csv"
	"github.com/signaux-faibles/fakeData/lib"
	"io"
	"log"
	"os"
)

func ReadAndRandomCCSF(source string, outputFileName string, outputSize int, mapping map[string]string) error {
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("error when closing", file)
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
	var dates []string

	wrote := 0
	for wrote <= outputSize {
		row, err := reader.Read()
		if err == io.EOF {
			log.Default().Println("EOF", source)
			break
		} else if err != nil {
			return err
		}

		if k, ok := mapping[row[2]]; ok {

			log.Default().Println("dates", dates)
			var output [6]string

			output[0] = row[0]
			output[1] = row[1]
			output[2] = k
			output[3] = lib.RandItemFrom(dates)
			output[4] = row[4]
			output[5] = row[5]

			err = writer.Write(output[0:6])
			if err != nil {
				return err
			}
			if outputSize > 100 {
				if mod := wrote % (outputSize / 100); mod == 0 {
					log.Default().Println("(ccsf) wrote ", wrote/(outputSize/100), "%")
				}
			}
			wrote++
		}
	}
	return nil
}
