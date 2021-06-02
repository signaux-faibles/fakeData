package main

import (
	"bufio"
	"encoding/csv"
	"github.com/signaux-faibles/fakeData/common"
	"io"
	"log"
	"os"
)

func ReadAndRandomComptes(source string, outputFileName string, outputSize int) (map[string]string, error) {
	mapping := make(map[string]string)
	sirets := make(map[string]string)

	// source
	sourceFile, err := os.Open(source)
	if err != nil {
		return nil, err
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
		return nil, err
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
	if err != nil {
		return nil, err
	}
	err = writer.Write(row)
	if err != nil {
		return nil, err
	}

	wrote := 0
	for wrote < outputSize {

		row, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				log.Default().Println("EOF for ", source)
				break
			}
			if perr, ok := err.(*csv.ParseError); ok && perr.Err == csv.ErrFieldCount {
				log.Default().Println("(skip line) ", perr)
				continue
			}
			return nil, err
		}

		siren := row[4]
		if len(siren) != 9 {
			continue
		}

		var newSiren string
		if _, ok := sirets[siren]; ok {
			newSiren = sirets[siren]
		} else {
			for {
				newSiren = common.RandStringBytesRmndr(9)
				if _, ok := sirets[newSiren]; !ok && newSiren != siren {
					break
				}
			}
		}

		siret := row[5]

		compte := row[2]
		var newSiret, newCompte string

		for {
			newSiret = newSiren + common.RandStringBytesRmndr(5)
			if _, ok := mapping[newSiret]; !ok && newSiret != siret {
				break
			}
		}
		for {
			newCompte = common.RandStringBytesRmndr(len(compte))
			if _, ok := mapping[newCompte]; !ok && newCompte != compte {
				break
			}
		}
		mapping[compte] = newCompte
		mapping[siren] = newSiren
		sirets[siret] = newSiret

		row[2] = newCompte
		row[4] = newSiren
		row[5] = newSiret

		err = writer.Write(row)
		if err != nil {
			return nil, err
		}
		if mod := wrote % (outputSize / 100); mod == 0 {
			log.Default().Println("(comptes) wrote ", wrote/(outputSize/100), "%")
			common.SkipSomeLines(reader, 777)
		}
		wrote++
	}
	return mapping, nil
}
