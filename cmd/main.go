package main

import (
	"errors"
	"os"

	"github.com/Filipza/excel-mapping-tool/internal/domain/dataimport"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func main() {

	file := loadFile("./test/test.xlsx")
	defer file.Close()

	uploadData := dataimport.UploadData{
		UploadedFile: file,
		UploadType:   "wkz",
	}

	ReadFile(uploadData)
}

func loadFile(path string) *os.File {
	file, err := os.Open(path)

	if err != nil {
		log.Debug(err)
		return nil
	}
	return file
}

func ReadFile(uploadData dataimport.UploadData) (*dataimport.MappingOptions, error) {

	xlsx, err := excelize.OpenReader(uploadData.UploadedFile)
	if err != nil {
		log.Debug(err)
		return nil, &dataimport.Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Die hochgeladene Datei konnte nicht verarbeitet werden. Überprüfe das Dateiformat.",
		}
	}

	ddOptions := make(map[string]string)
	ddOptions["wkz"] = "Webekostenzuschuss"
	ddOptions["stocks"] = "Lagerbestand"

	mappingOptions := dataimport.MappingOptions{
		DropdownOptions: ddOptions,
		TableSummary:    make([][]string, 0),
	}

	rows, err := xlsx.Rows(xlsx.GetSheetName(0))
	if err != nil {
		log.Debug(err)
		return nil, &dataimport.Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Es is ein Fehler beim Lesen der Reihen aufgetreten. Überprüfe die Datei.",
		}
	}

	for i := 0; i < 4 && rows.Next(); i++ {
		col, err := rows.Columns()

		if err != nil {
			log.Debug(err)
			return nil, &dataimport.Error{
				ErrTitle: "Parsingfehler",
				ErrMsg:   "Es is ein Fehler beim Lesen der Reihen aufgetreten. Überprüfe die Datei.",
			}
		}

		if len(col) == 0 {
			log.Debug(errors.New("first row is empty"))
			return nil, &dataimport.Error{
				ErrTitle: "Dateifehler",
				ErrMsg:   "Die erste Zeile der Datei ist leer. Diese muss für den Import die Tabellenköpfe enthalten.",
			}
		}

	}
	return &mappingOptions, nil
}
