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
		"connectionFee":      "Anschlussgebühr in Euro",
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
			// fmt.Println(cols)
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
	row := 2 // assumes that data begins at 2nd row

	for rows.Next() {
		row++

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

		// TODO: implement externalArticleNumber indentifier logic
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
					tariffObj.BasicCharge, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
				case "basicChargeRenewal": // ? richtiges Feld für "nach Aktionszeitraum?"
					tariffObj.BasicChargeRenewal, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
				case "leadype":
					tariffObj.LeadType, _ = strconv.Atoi(cellVal)
				case "provision":
					tariffObj.Provision, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
				case "xProvision":
					tariffObj.XProvision, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
				case "connectionFee":
					tariffObj.ConnectionFee, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
				case "dataVolume":
					tariffObj.DataVolume, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
				case "legalnote":
					tariffObj.LegalNote = cellVal
				case "pibLink":
					tariffObj.PibLink = cellVal
				case "highlight1", "highlight2", "highlight3", "highlight4", "highlight5":
					lastdigit, _ := strconv.Atoi(string(inst.MappingValue[len(inst.MappingValue)-1]))
					highlightArr[lastdigit-1] = cellVal
				case "bullet1", "bullet2", "bullet3", "bullet4", "bullet5", "bullet6":
					lastdigit := string(inst.MappingValue[len(inst.MappingValue)-1])
					key := fmt.Sprintf("tariff_inclusive_benefit%s", lastdigit)
					ok, i := containsKey(tariffObj.Bullets, key)

					if !ok {
						newOpt := product.Option{Key: key, Value: cellVal}
						tariffObj.Bullets = append(tariffObj.Bullets, &newOpt)
						continue
					}
					tariffObj.Bullets[i].Value = cellVal

				case "supplierWkz":
					ok, i := containsKey(tariffObj.Wkz, "supplier")

					if !ok {
						newOpt := product.Option{Key: "supplier", Value: cellVal}
						tariffObj.Wkz = append(tariffObj.Wkz, &newOpt)
						continue
					}
					tariffObj.Wkz[i].Value = cellVal
				case "tariffWkz":
					ok, i := containsKey(tariffObj.Wkz, "tariff")

					if !ok {
						newOpt := product.Option{Key: "tariff", Value: cellVal}
						tariffObj.Wkz = append(tariffObj.Wkz, &newOpt)
						continue
					}
					tariffObj.Wkz[i].Value = cellVal
					// case "supplierWkz", "tariffWkz":
					// 	re := regexp.MustCompile(`(.+?)Wkz`)
					// 	key := re.FindStringSubmatch(inst.MappingValue)

					// 	ok, i := containsKey(tariffObj.Wkz, key[0])

					// 	if !ok {
					// 		newOpt := product.Option{Key: key[0], Value: cellVal}
					// 		tariffObj.Wkz = append(tariffObj.Wkz, &newOpt)
					// 		continue
					// 	}
					// 	tariffObj.Wkz[i].Value = cellVal
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

	return nil, nil
}

func containsKey(opts []*product.Option, s string) (bool, int) {
	for i, b := range opts {
		if b.Key == s {
			return true, i
		}
	}
	return false, -1
}

func updateTariff(svc *mappingService, mi *MappingInstruction, file *excelize.File, identifierValue string, row int, sh string) {
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
				tariffObj.BasicCharge, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
			case "basicChargeRenewal": // ? richtiges Feld für "nach Aktionszeitraum?"
				tariffObj.BasicChargeRenewal, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
			case "leadype":
				tariffObj.LeadType, _ = strconv.Atoi(cellVal)
			case "provision":
				tariffObj.Provision, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
			case "xProvision":
				tariffObj.XProvision, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
			case "connectionFee":
				tariffObj.ConnectionFee, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
			case "dataVolume":
				tariffObj.DataVolume, _ = strconv.ParseFloat(strings.ReplaceAll(cellVal, ",", "."), 64)
			case "legalnote":
				tariffObj.LegalNote = cellVal
			case "pibLink":
				tariffObj.PibLink = cellVal
			case "highlight1", "highlight2", "highlight3", "highlight4", "highlight5":
				lastdigit, _ := strconv.Atoi(string(inst.MappingValue[len(inst.MappingValue)-1]))
				highlightArr[lastdigit-1] = cellVal
			case "bullet1", "bullet2", "bullet3", "bullet4", "bullet5", "bullet6":
				lastdigit := string(inst.MappingValue[len(inst.MappingValue)-1])
				key := fmt.Sprintf("tariff_inclusive_benefit%s", lastdigit)
				ok, i := containsKey(tariffObj.Bullets, key)

				if !ok {
					newOpt := product.Option{Key: key, Value: cellVal}
					tariffObj.Bullets = append(tariffObj.Bullets, &newOpt)
					continue
				}
				tariffObj.Bullets[i].Value = cellVal

			case "supplierWkz":
				ok, i := containsKey(tariffObj.Wkz, "supplier")

				if !ok {
					newOpt := product.Option{Key: "supplier", Value: cellVal}
					tariffObj.Wkz = append(tariffObj.Wkz, &newOpt)
					continue
				}
				tariffObj.Wkz[i].Value = cellVal
			case "tariffWkz":
				ok, i := containsKey(tariffObj.Wkz, "tariff")

				if !ok {
					newOpt := product.Option{Key: "tariff", Value: cellVal}
					tariffObj.Wkz = append(tariffObj.Wkz, &newOpt)
					continue
				}
				tariffObj.Wkz[i].Value = cellVal
				// case "supplierWkz", "tariffWkz":
				// 	re := regexp.MustCompile(`(.+?)Wkz`)
				// 	key := re.FindStringSubmatch(inst.MappingValue)

				// 	ok, i := containsKey(tariffObj.Wkz, key[0])

				// 	if !ok {
				// 		newOpt := product.Option{Key: key[0], Value: cellVal}
				// 		tariffObj.Wkz = append(tariffObj.Wkz, &newOpt)
				// 		continue
				// 	}
				// 	tariffObj.Wkz[i].Value = cellVal
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
