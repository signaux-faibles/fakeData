package main

import (
	"bufio"
	"encoding/csv"
	"github.com/signaux-faibles/fakeData/common"
	"io"
	"log"
	"os"
	"strings"
)

func ReadAndRandomEtablissements(source string, outputFileName string, outputSize int) (map[string]string, error) {
	mapping := make(map[string]string)
	sirens := make(map[string]string)

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
		siret := row[5]
		if len(siren) < 4 || len(siret) < len(siren) {
			continue
		}

		var newSiren string
		if _, ok := sirens[siren]; ok {
			newSiren = sirens[siren]
		} else {
			for {
				newSiren, err = common.FalsifyNumber(siren)
				if err != nil {
					break
				}
				if _, ok := sirens[newSiren]; !ok && newSiren != siren {
					break
				}
			}
		}

		compte := row[2]
		var newSiret, newCompte string

		for {
			suffix := strings.SplitAfter(siret, siren)[1]
			newSuffix, err := common.FalsifyNumber(suffix)
			if err != nil {
				continue
			}
			newSiret = newSiren + newSuffix
			if _, ok := mapping[newSiret]; !ok && newSiret != siret {
				break
			}
		}
		for {
			newCompte, err = common.FalsifyNumber(compte)
			if err != nil {
				break
			}
			if _, ok := mapping[newCompte]; !ok && newCompte != compte {
				break
			}
		}
		mapping[compte] = newCompte
		mapping[siret] = newSiret
		mapping[siren] = newSiren
		sirens[siren] = newSiren

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
