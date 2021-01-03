package utils

import (
	"encoding/csv"
	"os"
	"os/user"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
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

// PreprocessFilepath will check if the user has included the tilde character to signify
// the home directory and replaces it with the user's $HOME dir
func PreprocessFilepath(path string) string {
	// Check if the user has provided the tilde for home directory and expand it
	if !strings.HasPrefix(path, "~/") {
		// return as is
		return path
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Encountered an error while processing database path: %v", err)
	}
	return filepath.Join(usr.HomeDir, path[2:])
}
