package urssaf

import (
	"bufio"
	"encoding/csv"
	"github.com/signaux-faibles/fakeData/common"
	"io"
	"log"
	"os"
	"strings"
)

var pcollDateFormat = "02Jan2006"
var procedures []string
var cat_v2 []string

func ReadAndRandomPcoll(source string, outputFileName string, outputSize int, mapping map[string]string) error {
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

	wrote := 0
	for wrote <= outputSize {
		row, err := reader.Read()
		if err == io.EOF {
			log.Default().Println("EOF", source)
			break
		} else if err != nil {
			return err
		}
		siret := row[0]
		siren := row[1]
		if newSiret, found := mapping[siret]; found {
			if newSiren, found := mapping[siren]; found {

				cat_v2 = append(cat_v2, row[3])
				procedures = append(procedures, row[4])
				output := randomizePColl(newSiret, newSiren, row)

				err = writer.Write(output[0:])
				if err != nil {
					return err
				}
				if outputSize > 100 {
					if mod := wrote % (outputSize / 100); mod == 0 {
						log.Default().Println("(pcoll) wrote ", wrote/(outputSize/100), "%")
						common.SkipSomeLines(reader, 3.33)
					}
				}
				wrote++
			} else {
				log.Default().Println("siret not found", siret)
			}
		}
	}
	return nil
}

func randomizePColl(siret string, siren string, input []string) []string {
	var output [5]string

	output[0] = siret
	output[1] = siren

	randomized, err := common.RandDate(pcollDateFormat, input[2])
	if err != nil {
		output[2] = input[2]
	} else {
		output[2] = strings.ToUpper(randomized)
	}
	output[3] = common.RandItemFrom(cat_v2)
	output[4] = common.RandItemFrom(procedures)
	return output[:]
}
