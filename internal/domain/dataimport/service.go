package dataimport

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type MappingService interface {
	ReadFile(UploadData) (MappingOptions, error)
	WriteMapping(MappingInstruction) (MappingResult, error) // kann ich hier nicht direkt den custom Error type nutzen oder spricht was dagegen?
}

type mappingService struct {
	chanMap sync.Map
}

var DROPDOWN_OPTIONS = map[string]map[string]string{
	"tariff": {
		"monthlyPrice":               "Preis monatlich",
		"monthlyPriceAfterPromotion": "Preis monatlich nach Aktionszeitraum",
		"leadType":                   "Lead Type",
		"storeBonus":                 "Marktprämie",
		"onlineBonus":                "Onlineprämie",
		"connectionFeeEur":           "Anschlussgebühr in Euro",
		"inclDataVolumeGb":           "Inkl. Datenvolumen in GB",
		"legalNote":                  "Legalnote",
		"highlight1":                 "Highlight 1",
		"highlight2":                 "Highlight 2",
		"highlight3":                 "Highlight 3",
		"highlight4":                 "Highlight 4",
		"highlight5":                 "Highlight 5",
		"pibUrl":                     "Pib-URL",
		"connectionFee":              "Anschlusspreis",
		"monthlyBasePrice":           "Monatsgrundpreis",
		"inclBenefit1":               "Inklusiv-Benefit 1",
		"inclBenefit2":               "Inklusiv-Benefit 2",
		"inclBenefit3":               "Inklusiv-Benefit 3",
		"inclBenefit4":               "Inklusiv-Benefit 4",
		"inclBenefit5":               "Inklusiv-Benefit 5",
		"wkz":                        "WKZ",
	},
	"hardware": {
		"ek":          "EK",
		"manufactWkz": "Manufacturer WKZ",
		"ek24Wkz":     "ek24 WKZ",
	},
	"stocks": {
		"currentStock":  "Stock aktuell",
		"originalStock": "Stock original",
	},
}

func (svc *mappingService) ReadFile(ud *UploadData) (*MappingOptions, error) {
	if ud.Uuid == "" {
		ud.Uuid = uuid.New().String()
	}

	// creation of dir named after uuid
	dirPath := "../files/" + ud.Uuid + "/"

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, &Error{
			ErrTitle: "Verzeichnisfehler",
			ErrMsg:   "Verzeichnis konnte nicht erstellt werden",
		}
	}

	// removal of dir after timeout
	// TODO: Create channel as entrypoint for WriteMapping() to kill goroutine
	rfch := make(chan int)
	svc.chanMap.Store(ud.Uuid, rfch)

	go func(dirPath string) {
		time.Sleep(1800 * time.Second)
		os.RemoveAll(dirPath)
	}(dirPath)

	xlsx, err := excelize.OpenReader(ud.UploadedFile)
	if err != nil {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Datei konnte nicht verarbeitet werden und möglicherweise korrupt.",
		}
	}

	if err := xlsx.SaveAs(dirPath + "data.xlsx"); err != nil {
		fmt.Println(err)
	}

	mappingOptions := MappingOptions{
		TableSummary: make([][]string, 0),
	}

	var exists bool
	mappingOptions.DropdownOptions, exists = DROPDOWN_OPTIONS[ud.UploadType]
	if !exists {
		return nil, &Error{
			ErrTitle: "Fehlender/falscher Uploadtyp",
			ErrMsg:   fmt.Sprintf("Der Uploadtype %s ist unbekannt", ud.UploadType),
		}
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

	return &mappingOptions, nil
}

func (svg *mappingService) WriteMapping(mi *MappingInstruction) (*MappingResult, error) {

	return nil, nil
}
