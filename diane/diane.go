package diane

import "C"
import (
	"encoding/csv"
	"github.com/pkg/errors"
	"github.com/signaux-faibles/fakeData/common"
	"golang.org/x/text/encoding/unicode"
	"io"
	"log"
	"os"
	"strconv"

	//"strings"
	"github.com/TomOnTime/utfutil"
)

var utf16le = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()

func ReadAndRandomDiane(source string, outputFileName string, outputSize int, mapping map[string]string) error {
	// source
	sourceFile, err := utfutil.OpenFile(source, utfutil.UTF16LE)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	//buffer := bufio.NewReader(sourceFile)
	//err = dealingWithUTF16LEBom(buffer)
	if err != nil {
		return err
	}
	reader := csv.NewReader(sourceFile)
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
	if err != nil {
		return errors.Wrap(err, "Error reading header row")
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
		} else if err != nil {
			return err
		}
		siren := row[2]

		if newSiren, found := mapping[siren]; found {

			output := randomizeDiane(newSiren, row)

			err = writer.Write(output[0:])
			if err != nil {
				return err
			}
			if outputSize > 100 {
				if mod := wrote % (outputSize / 100); mod == 0 {
					log.Default().Println("(diane) wrote ", wrote/(outputSize/100), "%")
					//common.SkipSomeLines(reader, 3.33)
				}
			}
			wrote++
		}
	}
	return nil
}

// randomizeDiane see columns 15 (P), 35 (AJ)
//, unknown example value
func randomizeDiane(siren string, input []string) []string {
	var output []string

	for i, value := range input[:] {
		switch i {
		case 0: // marquee : 1
			output = append(output, strconv.Itoa(i))
		case 1: // nom de l'entreprise
			output = append(output, common.RandRaisonSociale(value))
		case 2: // siren
			output = append(output, siren)
		case 3, 4, 14, 15, 35: // statut juridique, procédure collective, nombre de mois
			output = append(output, value)
		//case 5, 6, 7, 8, 9, 10, 11, 24, 25, 26, 27, 30, 34, 38, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81:
		//	newValue, _ := common.RandInt(value)
		//	output = append(output, newValue)
		case 12, 13:
			newDate, err := common.RandDate("02/01/2006", value)
			if err != nil {
				log.Default().Println("Error when randomize", value, "(Siren is", siren, ") Set old value.", "Error is", err)
				output = append(output, value)
			}
			output = append(output, newDate)
		//case 16, 17, 18, 19, 20, 21, 22, 23, 28, 29, 31, 32, 33, 36, 37, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54:
		//	newValue, err := common.RandDecimal(value)
		//	if err != nil {
		//		log.Default().Println("Error when randomize", value, "(Siren is", siren, ") Set old value", "Error is", err)
		//		output = append(output, value)
		//	}
		//	output = append(output, newValue)
		default:
			newValue, err := common.RandNumber(value)
			if err != nil {
				log.Default().Println("Error when randomize", value, "(Siren is", siren, ") Set old value", "Error is", err)
				output = append(output, value)
			}
			output = append(output, newValue)

		}
	}
	return output
}
