package csv

import (
	"encoding/csv"
	"os"
)

// WriteIntoCSVFile -
func WriteIntoCSVFile(fileName string, datas [][]string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, data := range datas {
		w.Write(data)
	}

	return nil
}
