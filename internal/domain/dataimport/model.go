package dataimport

import (
	"errors"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// Contains file and type of import
type UploadData struct {
	UploadedFile io.Reader
	UploadType   string
}

type MappingOptions struct {
	DropdownOptions map[string]string
	TableHeaders    []string
	TableSummary    [][]string
}

// Instructions sent to BE after mapping in FE
type MappingInstruction struct {
	Mapping []MappingObject
}

type MappingObject struct {
	Index        int
	MappingValue string
}

type MappingResult struct {
	SuccessfulRows   int
	UnsuccessfulRows int
	FailedRows       []string
}

type Error struct {
	ErrTitle string
	ErrMsg   string
}

func (err *Error) Error() string {
	return fmt.Sprintf("Error: %s", err.ErrMsg)
}

func (uploadData UploadData) ReadFile() (*MappingOptions, *Error) {
	xlsx, err := excelize.OpenReader(uploadData.UploadedFile)
	if err != nil {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Die hochgeladene Datei konnte nicht verarbeitet werden. Überprüfe das Dateiformat.",
		}
	}

	ddOptions := make(map[string]string)
	ddOptions["wkz"] = "Webekostenzuschuss"
	ddOptions["stocks"] = "Lagerbestand"

	mappingOptions := MappingOptions{
		DropdownOptions: ddOptions,
		TableSummary:    make([][]string, 0),
	}

	rows, err := xlsx.Rows(xlsx.GetSheetName(0))
	if err != nil {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Es is ein Fehler beim Lesen der Reihen aufgetreten. Überprüfe die Datei.",
		}
	}

	for i := 0; i < 4 && rows.Next(); i++ {
		col, err := rows.Columns()

		if err != nil {
			log.Debug(err)
			return nil, &Error{
				ErrTitle: "Parsingfehler",
				ErrMsg:   "Es is ein Fehler beim Lesen der Reihen aufgetreten. Überprüfe die Datei.",
			}
		}

		if len(col) == 0 {
			log.Debug(errors.New("first row is empty"))
			return nil, &Error{
				ErrTitle: "Dateifehler",
				ErrMsg:   "Die erste Zeile der Datei ist leer. Diese muss für den Import die Tabellenköpfe enthalten.",
			}
		}

		if i == 0 {
			mappingOptions.TableHeaders = col
			continue
		}

		mappingOptions.TableSummary = append(mappingOptions.TableSummary, col)

	}
	fmt.Println(mappingOptions)
	return &mappingOptions, nil
}
