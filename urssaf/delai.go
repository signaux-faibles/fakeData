package urssaf

import (
	"bufio"
	"encoding/csv"
	"github.com/pkg/errors"
	"github.com/signaux-faibles/fakeData/common"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
)

//Urssaf_gestion;Dep;Numero_compte_externe;Numero_structure;Date_creation;Date_echeance;Duree_delai;Denomination_premiere_ligne;Indic_6M;Annee_creation;Montant_global_echeancier;Code_externe_stade;Code_externe_action
//311;31;090048281379700073;3110201303256;24/03/2021;28/09/2021;188;OCCITANIA;SUP;2021;5697.00;APPROB;SUR PO

func ReadAndRandomDelais(source string, outputFileName string, outputSize int, mapping map[string]string) error {
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

		if k, ok := mapping[row[2]]; ok {

			var output [13]string

			c := buildOrGetCoef(coef, row[2])

			output[0] = row[0]
			output[1] = row[1]
			output[2] = k
			output[3] = "aSiren-" + strconv.Itoa(rand.Int())
			output[4] = "31/05/2021"
			output[5] = "31/05/2021"
			output[6] = row[6]
			output[7] = common.RandRaisonSociale(row[7])
			output[8] = row[8]
			output[9] = row[9]
			output[10], err = newMontant(row[10], c)
			if err != nil {
				return errors.Cause(err)
			}
			output[11] = row[11]
			output[12] = row[12]

			err = writer.Write(output[0:13])
			if err != nil {
				return err
			}
			if outputSize > 100 {
				if mod := wrote % (outputSize / 100); mod == 0 {
					log.Default().Println("(delai) wrote ", wrote/(outputSize/100), "%")
				}
			}
			wrote++
		}
	}
	return nil
}

func buildOrGetCoef(coef map[string]float64, value string) float64 {
	var c float64
	var found bool

	if c, found = coef[value]; !found {
		c = common.RandCoeff()
		coef[value] = c
	}
	return c
}

func newMontant(value string, c float64) (string, error) {
	montant, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", err
	}
	newMontant := int(montant * c)
	return strconv.Itoa(newMontant), nil
}
