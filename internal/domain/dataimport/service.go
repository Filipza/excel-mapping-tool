package dataimport

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type MappingService interface {
	ReadFile(UploadData) (MappingOptions, error)
	WriteMapping(MappingInstruction) (MappingResult, error)
}

type mappingService struct {
}

func (svc *mappingService) ReadFile(ud *UploadData) (*MappingOptions, error) {
	xlsx, err := excelize.OpenReader(ud.UploadedFile)
	if err != nil {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Datei konnte nicht verarbeitet werden und möglicherweise korrupt.",
		}
	}

	if ud.UploadType == "" {
		return nil, &Error{
			ErrTitle: "Fehlender Uploadtyp",
			ErrMsg:   "Uploadtyp wurde nicht ausgewählt",
		}
	}

	// TODO: anpassen
	ddOptions := make(map[string]string)
	ddOptions["wkz"] = "Webekostenzuschuss"
	ddOptions["stocks"] = "Lagerbestand"

	mappingOptions := MappingOptions{
		DropdownOptions: ddOptions,
		TableSummary:    make([][]string, 0),
	}

	sheetLists := xlsx.GetSheetList()
	if len(sheetLists) == 0 {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Fehlerhafte Excel-Datei",
			ErrMsg:   "Datei enthält keine Arbeitsblätter",
		}
	}

	rows, err := xlsx.Rows(sheetLists[0])
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
				ErrTitle: "Leerzeile",
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
