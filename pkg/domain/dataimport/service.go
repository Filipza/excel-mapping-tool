package dataimport

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/crud"
	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/product"
	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/tariff"
	"github.com/Filipza/excel-mapping-tool/internal/settings"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type MappingService interface {
	ReadFile(UploadData) (MappingOptions, error)
	WriteMapping(MappingInstruction) (MappingResult, error) // kann ich hier nicht direkt den custom Error type nutzen oder spricht was dagegen?
}

type mappingService struct {
	chanMap       sync.Map
	tariffAdapter crud.CRUDService[tariff.TariffCRUD, tariff.TariffLookup]
}

var DROPDOWN_OPTIONS = map[string]map[string]string{
	"tariff": {
		"ebootisId":          "EbootisId",
		"basicCharge":        "Preis monatlich",
		"basicChargeRenewal": "Preis monatlich nach Aktionszeitraum",
		"leadType":           "Lead Type",
		"provision":          "Marktprämie",
		"xProvision":         "Onlineprämie",
		"connectionFee":      "Anschlussgebühr (ohne EUR-Zeichen)",
		"dataVolume":         "Inkl. Datenvolumen in GB",
		"legalNote":          "Legalnote",
		"pibLink":            "Pib-URL",
		"highlight1":         "Highlight 1",
		"highlight2":         "Highlight 2",
		"highlight3":         "Highlight 3",
		"highlight4":         "Highlight 4",
		"highlight5":         "Highlight 5",
		"bullet1":            "Inklusiv-Benefit 1",
		"bullet2":            "Inklusiv-Benefit 2",
		"bullet3":            "Inklusiv-Benefit 3",
		"bullet4":            "Inklusiv-Benefit 4",
		"bullet5":            "Inklusiv-Benefit 5",
		"bullet6":            "Inklusiv-Benefit 6",
		"supplierWkz":        "Supplier WKZ",
		"tariffWkz":          "Tariff WKZ",
	},
	"hardware": {
		"ebootisId":             "EbootisId",
		"externalArticleNumber": "Exerterne Artikelnr.",
		"ek":                    "EK",
		"manufactWkz":           "Manufacturer WKZ",
		"ek24Wkz":               "ek24 WKZ",
	},
	"stocks": {
		"ebootisId":             "EbootisId",
		"externalArticleNumber": "Exerterne Artikelnr.",
		"currentStock":          "Stock aktuell",
		"originalStock":         "Stock original",
	},
}

func (svc *mappingService) ReadFile(ud *UploadData) (*MappingOptions, error) {
	if ud.Uuid == "" {
		ud.Uuid = uuid.New().String()
	}

	// creation of dir named after uuid
	dirPath := "/tmp/" + ud.Uuid + "/"

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, &Error{
			ErrTitle: "Verzeichnisfehler",
			ErrMsg:   "Verzeichnis konnte nicht erstellt werden",
		}
	}

	// TODO: Write tests for goroutine/channels
	// removal of dir after timeout
	cleanupCh := make(chan bool)
	svc.chanMap.Store(ud.Uuid, cleanupCh)

	go func(dirPath string, ch <-chan bool) {
		timer := time.NewTimer(1800 * time.Second)

		select {
		case <-timer.C:
			os.RemoveAll(dirPath)
		case <-ch:
			return
		}
	}(dirPath, cleanupCh)

	// TODO: Test timer channel and find out if Stop() actually doesnt stop timer
	// type mappingService struct {
	// 	timerMap sync.Map
	// }
	// timer := time.NewTimer(1800 * time.Second)
	//svc.timerMap.Store(ud.Uuid, timer.C)
	// go func(dirPath string, ch <-chan time.Time) {
	// 	<-timer.C

	// 	os.RemoveAll(dirPath)
	// }(dirPath, timer.C)
	// timer.Stop()

	file, err := excelize.OpenReader(ud.UploadedFile)
	if err != nil {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Datei konnte nicht verarbeitet werden und möglicherweise korrupt.",
		}
	}
	defer file.Close()

	err = file.SaveAs(dirPath + "data.xlsx")
	if err != nil {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Speicherfehler",
			ErrMsg:   "Datei konnte nicht abespeichert werden.",
		}
	}

	mappingOptions := MappingOptions{
		TableSummary: make([][]string, 0),
		Uuid:         ud.Uuid,
	}

	var exists bool
	mappingOptions.DropdownOptions, exists = DROPDOWN_OPTIONS[ud.UploadType]
	if !exists {
		return nil, &Error{
			ErrTitle: "Fehlender/falscher Uploadtyp",
			ErrMsg:   fmt.Sprintf("Der Uploadtype %s ist unbekannt", ud.UploadType),
		}
	}

	sheetLists := file.GetSheetList()
	if len(sheetLists) == 0 {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Fehlerhafte Excel-Datei",
			ErrMsg:   "Datei enthält keine Arbeitsblätter",
		}
	}

	rows, err := file.Rows(sheetLists[0])
	if err != nil {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Es is ein Fehler beim Lesen der Reihen aufgetreten. Überprüfe die Datei.",
		}
	}

	for i := 0; i < 4 && rows.Next(); i++ {
		cols, err := rows.Columns()

		if err != nil {
			log.Debug(err)
			return nil, &Error{
				ErrTitle: "Parsingfehler",
				ErrMsg:   "Es is ein Fehler beim Lesen der Reihen aufgetreten. Überprüfe die Datei.",
			}
		}

		if len(cols) == 0 {
			log.Debug(errors.New("first row is empty"))
			return nil, &Error{
				ErrTitle: "Leerzeile",
				ErrMsg:   "Die erste Zeile der Datei ist leer. Diese muss für den Import die Tabellenköpfe enthalten.",
			}
		}

		if i == 0 {
			mappingOptions.TableHeaders = cols
			continue
		}

		mappingOptions.TableSummary = append(mappingOptions.TableSummary, cols)
	}

	return &mappingOptions, nil
}

func (svc *mappingService) WriteMapping(mi *MappingInstruction) (*MappingResult, error) {

	exists, idIndex, _ := mi.GetIdentifierIndex()
	idCol := idIndex + 1

	fmt.Println(idCol)

	if !exists {
		return nil, &Error{
			ErrTitle: "Fehlende EbootisID / externe Artikelnummer",
			ErrMsg:   "Keine der Spalten wurde der EbootisID / externen Artikelnummer zugewiesen",
		}
	}

	file, err := excelize.OpenFile("/tmp/" + mi.Uuid + "/data.xlsx")
	if err != nil {
		return nil, &Error{
			ErrTitle: "Fehler beim Öffnen der Datei",
			ErrMsg:   "Die zu bearbeitende Excel-Datei konnte nicht geöffnet werden",
		}
	}
	defer file.Close()

	sheetLists := file.GetSheetList()
	if len(sheetLists) == 0 {
		log.Debug(err)
		return nil, &Error{
			ErrTitle: "Fehlerhafte Excel-Datei",
			ErrMsg:   "Datei enthält keine Arbeitsblätter",
		}
	}

	sh := sheetLists[0]
	rows, _ := file.Rows(sh)

	for row := 1; rows.Next(); row++ {

		if row == 1 {
			continue // skips header
		}

		idCoords, err := excelize.CoordinatesToCellName(idCol, row)
		if err != nil {
			// TODO: Adjustments necessary. Skip and log row when error occurs
			return nil, &Error{
				ErrTitle: "Koordinatenfehler",
				ErrMsg:   fmt.Sprintf("Identifikationskoordinate konnte in Zeile %v nicht in Zellname umgewandelt werden", row),
			}
		}
		identifierValue, err := file.GetCellValue(sh, idCoords)
		if err != nil {
			// TODO: Adjustments necessary. Skip and log row when error occurs
			return nil, &Error{
				ErrTitle: "Zellen-Lesefehler",
				ErrMsg:   fmt.Sprintf("Der Zelleninhalt der Zelle %s konnte nicht gelesen werden", idCoords),
			}
		}

		svc.updateTariff(mi, file, identifierValue, row, sh)
	}

	return nil, nil
}

func (svc *mappingService) updateTariff(mi *MappingInstruction, file *excelize.File, identifierValue string, row int, sh string) {
	listResult, _ := svc.tariffAdapter.List(settings.Option{Name: "ebootis_id", Value: identifierValue})
	for _, lookupObj := range listResult {
		tariffObj, _ := svc.tariffAdapter.Read(lookupObj.Id)

		highlightArr := make([]string, 5)
		copy(highlightArr, tariffObj.Highlights)

		for _, inst := range mi.Mapping {
			coords, _ := excelize.CoordinatesToCellName(inst.ColIndex, row)
			cellVal, _ := file.GetCellValue(sh, coords)

			switch inst.MappingValue {
			case "basicCharge":
				// TODO: getPeriodFloat ok?
				tariffObj.BasicCharge, _ = getPeriodFloat(cellVal)

				if len(tariffObj.PricingIntervals) > 0 {
					tariffObj.PricingIntervals[0].Price, _ = getPeriodFloat(cellVal)
				}

				// ! writebullet für tariff_monthly_price > problem, da in tariff_monthly_price auch
				// ! strings wie "34.99€ (ab dem 13. Monat 69.99€)" stehen
				cellVal += " €"
				tariffObj.Bullets = writeOptionArr(tariffObj.Bullets, "tariff_monthly_price", cellVal)
			case "basicChargeRenewal":
				tariffObj.BasicChargeRenewal, _ = getPeriodFloat(cellVal)

				if len(tariffObj.PricingIntervals) > 1 {
					tariffObj.PricingIntervals[1].Price, _ = getPeriodFloat(cellVal)
				}
			case "leadype":
				tariffObj.LeadType, _ = strconv.Atoi(cellVal)
			case "provision":
				tariffObj.Provision, _ = getPeriodFloat(cellVal)
			case "xProvision":
				tariffObj.XProvision, _ = getPeriodFloat(cellVal)
			case "connectionFee":
				tariffObj.ConnectionFee, _ = getPeriodFloat(cellVal)

				// connectionFee schreiben mit neuer writeBullet func() key: tariff_connection_fee (Eurozeichen hinzufügen)
				cellVal += " €"
				tariffObj.Bullets = writeOptionArr(tariffObj.Bullets, "tariff_connection_fee", cellVal)

				ok, i := containsKey(tariffObj.Bullets, "tariff_monthly_price")

				if !ok {
					newOpt := product.Option{Key: "tariff_monthly_price", Value: cellVal}
					tariffObj.Bullets = append(tariffObj.Bullets, &newOpt)
					continue
				}
				tariffObj.Bullets[i].Value = cellVal
			case "dataVolume":
				tariffObj.DataVolume, _ = getPeriodFloat(cellVal)
			case "legalnote":
				tariffObj.LegalNote = cellVal
			case "pibLink":
				tariffObj.PibLink = cellVal
			case "highlight1", "highlight2", "highlight3", "highlight4", "highlight5":
				// Extracts last digit of "highlightx" key, to get array index
				lastdigit := int(inst.MappingValue[len(inst.MappingValue)-1]) - '0'
				highlightArr[lastdigit-1] = cellVal
			case "bullet1", "bullet2", "bullet3", "bullet4", "bullet5", "bullet6":
				// Extracts last digit of "bulletx" key, to search for according "tariff_inclusive_benefitsx" key
				lastdigit := string(inst.MappingValue[len(inst.MappingValue)-1])
				key := fmt.Sprintf("tariff_inclusive_benefit%s", lastdigit)
				//
				ok, i := containsKey(tariffObj.Bullets, key)

				if !ok {
					newOpt := product.Option{Key: key, Value: cellVal}
					tariffObj.Bullets = append(tariffObj.Bullets, &newOpt)
					continue
				}
				tariffObj.Bullets[i].Value = cellVal
			case "supplierWkz", "tariffWkz":
				key := strings.TrimRight(inst.MappingValue, "Wkz")
				ok, i := containsKey(tariffObj.Wkz, key)

				if !ok {
					newOpt := product.Option{Key: key, Value: cellVal}
					tariffObj.Wkz = append(tariffObj.Wkz, &newOpt)
					continue
				}
				tariffObj.Wkz[i].Value = cellVal
			}
		}

		// reducing array to minimum length and writing to tariffCRUD
		lenDiff := 0
		for i := len(highlightArr) - 1; i >= 0; i-- {
			if highlightArr[i] != "" {
				break
			}
			lenDiff++
		}
		highlightArr = highlightArr[:len(highlightArr)-lenDiff]
		copy(tariffObj.Highlights, highlightArr)
	}
}

func writeOptionArr(arr []*product.Option, key string, cellVal string) []*product.Option {
	for i, b := range arr {
		if b.Key == key {
			arr[i].Value = cellVal
			return arr
		}
	}

	newOpt := product.Option{Key: key, Value: cellVal}
	arr = append(arr, &newOpt)
	return arr
}

func getPeriodFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
}
