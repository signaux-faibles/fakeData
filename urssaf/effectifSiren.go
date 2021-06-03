package urssaf

import (
	"bufio"
	"encoding/csv"
	"github.com/signaux-faibles/fakeData/common"
	"io"
	"log"
	"os"
)

func ReadAndRandomEffectifSiren(source string, outputFileName string, outputSize int, mapping map[string]string) error {
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
		siren := row[0]
		if k, ok := mapping[siren]; ok {

			output := randomizeEffectifLine(k, row)

			err = writer.Write(output[0:])
			if err != nil {
				return err
			}
			if outputSize > 100 {
				if mod := wrote % (outputSize / 100); mod == 0 {
					log.Default().Println("(effectif) wrote ", wrote/(outputSize/100), "%")
					common.SkipSomeLines(reader, 3.33)
				}
			}
			wrote++
		}
	}
	return nil
}

func randomizeEffectifLine(key string, input []string) []string {
	var output []string
	output = append(output, key)
	for i, effectif := range input[:] {
		if i == 0 {
			continue
		}
		newEffectif, _ := common.RandEffectif(effectif)
		output = append(output, newEffectif)
	}
	output[len(output)-1] = common.RandRaisonSociale(input[len(input)-1])
	return output
}
