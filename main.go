package main

import (
	"flag"
	"fmt"
	"github.com/signaux-faibles/fakeData/urssaf"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Pour déformer des fichiers réels pour créer un dataset consistent avec anonymisation des entreprises

type randomizer func(string, string, int, map[string]string) error

func init() {
	log.Default().Println("init randoÎmization")
	rand.Seed(time.Now().UTC().UnixNano())
}

// run execute une fonction de randomisation
func run(name string, handler randomizer, mapping map[string]string) error {
	source := viper.GetString("input.files." + name)
	outputFile := outputFileNamePrefixer(viper.GetString("output.prefix"), source)
	outputSize := viper.GetInt("output.size")
	log.Default().Print("Fake ", name, ": ")
	err := handler(source, outputFile, outputSize, mapping)
	if err != nil {
		fmt.Println("Fail : " + err.Error())
		fmt.Println("Interruption.")
		return fmt.Errorf("interruption")
	}
	log.Default().Println("OK -> ", outputFile)
	return nil
}

func main() {
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Erreur à la lecture de la configuration:" + err.Error())
	}

	// Génération des numéros de comptes
	comptesFilename := viper.GetString("input.files.comptes")
	outputCompte := outputFileNamePrefixer(viper.GetString("output.prefix"), comptesFilename)
	outputSize := viper.GetInt("output.size")
	log.Default().Print("Fake comptes: ")
	// mapping contains mapping about compte & siret
	mapping, err := ReadAndRandomEtablissements(comptesFilename, outputCompte, outputSize*10)
	if err != nil {
		fmt.Println("Fail : " + err.Error())
		fmt.Println("Interruption.")
	} else {
		log.Default().Println("OK -> ", outputCompte)
	}
	var randomizers = map[string]randomizer{
		"delais":        urssaf.ReadAndRandomDelais,
		"debits":        urssaf.ReadAndRandomDebits,
		"cotisations":   urssaf.ReadAndRandomCotisations,
		"ccsf":          urssaf.ReadAndRandomCCSF,
		"effectifSiren": urssaf.ReadAndRandomEffectifSiren,
		"effectifSiret": urssaf.ReadAndRandomEffectifSiret,
		"pcoll":         urssaf.ReadAndRandomPcoll,
	}
	order := []string{"cotisations", "debits", "delais", "ccsf", "effectifSiren", "effectifSiret", "pcoll"}

	for _, k := range order {
		err := run(k, randomizers[k], mapping)
		if err != nil {
			panic(err)
		}
	}
}

func outputFileNamePrefixer(prefixOutput string, fileName string) string {
	path := strings.Split(fileName, "/")
	path[len(path)-1] = prefixOutput + path[len(path)-1]
	return strings.Join(path, "/")
}
