package csvreader

import (
	"encoding/csv"
	"os"
)

func ReadCsvFile(path string) ([]string, error) {
	var urls []string

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, row := range data {
		urls = append(urls, "https://"+row[0])
	}

	return urls, nil
}
