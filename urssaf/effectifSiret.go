package urssaf

import (
	"bufio"
	"encoding/csv"
	"github.com/signaux-faibles/fakeData/common"
	"io"
	"log"
	"os"
)

func ReadAndRandomEffectifSiret(source string, outputFileName string, outputSize int, mapping map[string]string) error {
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
		compte := row[0]
		siret := row[1]
		if newCompte, found := mapping[compte]; found {
			if newSiret, found := mapping[siret]; found {
				output := randomizeEffectifSiretLine(newCompte, newSiret, row)

				err = writer.Write(output[0:])
				if err != nil {
					return err
				}
				if outputSize > 100 {
					if mod := wrote % (outputSize / 100); mod == 0 {
						log.Default().Println("(effectif-siret) wrote ", wrote/(outputSize/100), "%")
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

func randomizeEffectifSiretLine(compte string, siret string, input []string) []string {
	var output []string
	for i, value := range input[:] {
		switch i {
		case 0:
			output = append(output, compte)
		case 1:
			output = append(output, siret)
		case 2:
			output = append(output, common.RandRaisonSociale(value))
		case 3:
			output = append(output, input[i])
		case 4:
			output = append(output, input[i])
		case len(input) - 2:
			output = append(output, input[len(input)-2])
		case len(input) - 1:
			output = append(output, input[len(input)-1])
		default:
			newEffectif, _ := common.RandEffectif(value)
			output = append(output, newEffectif)
		}

	}
	return output
}
