package utils

import (
	"encoding/csv"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
)

func GetTransactionId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func ReadCSVFile(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Could not close file: %s", fileName)
		}
	}()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
