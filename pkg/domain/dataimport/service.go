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
	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/hardware"
	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/product"
	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/tariff"
	"github.com/Filipza/excel-mapping-tool/internal/settings"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

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
		"price":                 "EK",
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

type MappingService interface {
	ReadFile(*UploadData) (*MappingOptions, error)
	WriteMapping(*MappingInstruction) (*MappingResult, error) // kann ich hier nicht direkt den custom Error type nutzen oder spricht was dagegen?
}

type mappingService struct {
	chanMap         sync.Map
	tariffAdapter   crud.CRUDService[tariff.TariffCRUD, tariff.TariffLookup]
	hardwareAdapter crud.CRUDService[hardware.HardwareCRUD, hardware.HardwareLookup]
	// variantAdapter  crud.CRUDService[hardware.VariantCRUD, hardware.VariantLookup]
}

func NewMappingService(tariffAdapter crud.CRUDService[tariff.TariffCRUD, tariff.TariffLookup], hardwareAdapter crud.CRUDService[hardware.HardwareCRUD, hardware.HardwareLookup]) MappingService {
	return &mappingService{tariffAdapter: tariffAdapter, hardwareAdapter: hardwareAdapter}
}

func (svc *mappingService) ReadFile(ud *UploadData) (*MappingOptions, error) {
	// assign Uuid in case of e.g cli/standalone application upload
	if ud.Uuid == "" {
		ud.Uuid = uuid.New().String()
	}

	// creation of dir in linux temp folder named after uuid
	dirPath := "/tmp/" + ud.Uuid + "/"

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, &Error{
			ErrTitle: "Verzeichnisfehler",
			ErrMsg:   "Verzeichnis konnte nicht erstellt werden",
		}
	}

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

	file, err := excelize.OpenReader(ud.UploadedFile)
	if err != nil {
		//TODO: errorf schreiben wo im code ich mich befinde
		log.Error(err)
		return nil, &Error{
			ErrTitle: "Parsingfehler",
			ErrMsg:   "Datei konnte nicht verarbeitet werden und möglicherweise korrupt.",
		}
	}
	defer file.Close()

	err = file.SaveAs(dirPath + "data.xlsx")
	if err != nil {
		log.Error(err)
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
	var err error
	// TODO: WIE BEKOMME ICH DEN ERROR DIREKT AUS DER FUNC ZUGEWIESEN WENN DEFERED?
	defer func() error {
		err := os.RemoveAll("/tmp/" + mi.Uuid + "/")
		return err
	}()

	if err != nil {
		return nil, &Error{
			ErrTitle: "Bereinigungsfehler",
			ErrMsg:   "Fehler beim Bereinigen der geschriebenen Daten.",
		}
	}

	if cancelInfo, ok := svc.chanMap.Load(mi.Uuid); ok {
		// send signal to cancel the cleanup routine
		if cancelChan, ok := cancelInfo.(chan bool); ok {
			cancelChan <- true
		}
	}

	result := &MappingResult{}

	exists, idIndex, idType := mi.GetIdentifierIndex()
	idCol := idIndex + 1

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

	editedCRUDmap := make(map[string]*editedCRUDobj)

	sh := sheetLists[0]
	rows, _ := file.Rows(sh)

	for row := 1; rows.Next(); row++ {

		if row == 1 {
			continue // skips header row
		}

		idCoords, err := excelize.CoordinatesToCellName(idCol, row)
		if err != nil {
			result.FailedRows = append(result.FailedRows, Error{
				ErrTitle: "Koordinatenfehler",
				ErrMsg:   fmt.Sprintf("Identifikationskoordinate konnte in Zeile %v nicht in Zellname umgewandelt werden", row),
			})
			continue
		}
		identifierValue, err := file.GetCellValue(sh, idCoords)
		if err != nil {
			result.FailedRows = append(result.FailedRows, Error{
				ErrTitle: "Zellen-Lesefehler",
				ErrMsg:   fmt.Sprintf("Der Zelleninhalt der Zelle %s konnte nicht gelesen werden", idCoords),
			})
			continue
		}

		var updateErr *Error

		switch mi.UploadType {
		case "tariff":
			updateErr = svc.updateTariff(mi, file, identifierValue, row, sh)
		case "hardware":
			updateErr = svc.updateHardware(mi, file, identifierValue, idType, row, sh, editedCRUDmap)

			// case "stocks":
			// updateErr = svc.updateStocks()
		}

		if updateErr != nil {
			result.UnsuccessfulRows++
			result.FailedRows = append(result.FailedRows, *updateErr)
			continue
		}
		result.SuccessfulRows++
	}

	if len(editedCRUDmap) > 0 {
		for _, v := range editedCRUDmap {
			if !v.hasError {
				if _, err := svc.hardwareAdapter.Update(v.hardwareCRUD.Id, v.hardwareCRUD); err != nil {
					result.FailedRows = append(result.FailedRows, Error{
						ErrTitle: "Hardware Speicherfehler",
						ErrMsg:   fmt.Sprintf("Update von Eardware %s konnte nicht durchgeführt werden", v.hardwareCRUD.Id),
					})
				}
			}
		}
	}

	return result, err
}

func (svc *mappingService) updateTariff(mi *MappingInstruction, file *excelize.File, identifierValue string, row int, sh string) *Error {
	listResult, err := svc.tariffAdapter.List(settings.Option{Name: "ebootis_id", Value: identifierValue})
	log.Error(err)
	if err != nil {
		return &Error{
			ErrTitle: "Identifizierungs-Fehler",
			ErrMsg:   fmt.Sprintf("Fehler in Zeile %d. Es konnten keine Tarifobjekte mit der EbootisId '%s' ermittelt werden", row, identifierValue),
		}
	}
	for _, lookupObj := range listResult {
		tariffObj, err := svc.tariffAdapter.Read(lookupObj.Id)
		log.Error(err)
		if err != nil {
			return &Error{
				ErrTitle: "Identifizierungs-Fehler",
				ErrMsg:   fmt.Sprintf("Fehler in Zeile %d. Es konnten kein Tarifobjekt mit der Id %s ermittelt werden.", row, lookupObj.Id),
			}
		}

		highlightArr := make([]string, 5)
		copy(highlightArr, tariffObj.Highlights)

		for _, inst := range mi.Mapping {
			coords, err := excelize.CoordinatesToCellName(inst.ColIndex, row)
			if err != nil {
				return &Error{
					ErrTitle: "Koordinatenfehler",
					ErrMsg:   fmt.Sprintf("Fehler in Zeile %d, Spalte %d. Es die dazugehörige Excel-Koordinate konnte nicht konvertiert werden", row, inst.ColIndex),
				}
			}
			cellVal, err := file.GetCellValue(sh, coords)
			if err != nil {
				return &Error{
					ErrTitle: "Lesefehler",
					ErrMsg:   fmt.Sprintf("Wert der Zelle %s konnte nicht ausgelesen werden", coords),
				}
			}

			switch inst.MappingValue {
			case "basicCharge":
				tariffObj.BasicCharge, _ = getPeriodFloat(cellVal)

				if len(tariffObj.PricingIntervals) > 0 {
					tariffObj.PricingIntervals[0].Price, _ = getPeriodFloat(cellVal)
				}

				// ! writebullet für tariff_monthly_price > problem, da in tariff_monthly_price auch
				// ! strings wie "34.99€ (ab dem 13. Monat 69.99€)" stehen
				// Klärung mit Produktmanagement. Siehe Notizen
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

				cellVal += " €" // TODO: Produktmanagement nahebringen
				tariffObj.Bullets = writeOptionArr(tariffObj.Bullets, "tariff_connection_fee", cellVal)
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
				// Extracts last digit of "bulletx" key, to get according "tariff_inclusive_benefitsx" key
				lastdigit := string(inst.MappingValue[len(inst.MappingValue)-1])
				key := fmt.Sprintf("tariff_inclusive_benefit%s", lastdigit)
				tariffObj.Bullets = writeOptionArr(tariffObj.Bullets, key, cellVal)
			case "supplierWkz", "tariffWkz":
				key := strings.TrimRight(inst.MappingValue, "Wkz")
				tariffObj.Wkz = writeOptionArr(tariffObj.Wkz, key, cellVal)
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

		// write into db
		svc.tariffAdapter.Update(lookupObj.Id, tariffObj)
	}
	return nil
}

func (svc *mappingService) updateHardware(mi *MappingInstruction, file *excelize.File, identifierValue string, idType string, row int, sh string, editedCRUDmap map[string]*editedCRUDobj) *Error {
	var err error
	var hardwareLookupList []*hardware.HardwareLookup

	switch idType {
	case "ebootisId":
		hardwareLookupList, err = svc.hardwareAdapter.List(settings.Option{Name: "variants.ebootis_id", Value: identifierValue})
	case "externalArticleNumber":
		hardwareLookupList, err = svc.hardwareAdapter.List(settings.Option{Name: "variants.external_articlenumber", Value: identifierValue})
	}

	log.Error(err)
	if err != nil {
		return &Error{
			ErrTitle: "Identifizierungs-Fehler",
			ErrMsg:   fmt.Sprintf("Fehler in Zeile %d. Es konnten keine Hardware mit dem Identifikator '%s' ermittelt werden", row, identifierValue),
		}
	}

	for _, listResult := range hardwareLookupList {
		var hardwareObj *editedCRUDobj

		// check if hardwareCRUD already in editedCRUDmap to prevent unecessary calls to hardwareAdapter
		// insert into editedCRUDmap if not present
		if _, ok := editedCRUDmap[listResult.Id]; !ok {
			readObj, err := svc.hardwareAdapter.Read(listResult.Id)
			if err != nil {
				return &Error{
					ErrTitle: "Identifizierungs-Fehler",
					ErrMsg:   fmt.Sprintf("Fehler in Zeile %d. Es konnten keine Hardware mit dem Identifikator '%s' ermittelt werden", row, identifierValue),
				}
			}

			hardwareObj = &editedCRUDobj{readObj, false}
			// add editedCRUDobj to editedCRUDmap if not present
			editedCRUDmap[listResult.Id] = hardwareObj
		} else {
			hardwareObj = editedCRUDmap[listResult.Id]
		}

		for _, inst := range mi.Mapping {
			coords, err := excelize.CoordinatesToCellName(inst.ColIndex, row)
			if err != nil {
				hardwareObj.hasError = true
				return &Error{
					ErrTitle: "Koordinatenfehler",
					ErrMsg:   fmt.Sprintf("Fehler in Zeile %d, Spalte %d. Es die dazugehörige Excel-Koordinate konnte nicht konvertiert werden", row, inst.ColIndex),
				}
			}
			cellVal, err := file.GetCellValue(sh, coords)
			if err != nil {
				hardwareObj.hasError = true
				return &Error{
					ErrTitle: "Lesefehler",
					ErrMsg:   fmt.Sprintf("Wert der Zelle %s konnte nicht ausgelesen werden", coords),
				}
			}

			switch inst.MappingValue {
			case "manufactWkz":
				hardwareObj.hardwareCRUD.Wkz = writeOptionArr(hardwareObj.hardwareCRUD.Wkz, "manufacturer", cellVal)
			case "ek24Wkz":
				hardwareObj.hardwareCRUD.Wkz = writeOptionArr(hardwareObj.hardwareCRUD.Wkz, "ek24", cellVal)
			case "price":
				switch idType {
				case "ebootisId":
					if variant, ok := hardwareObj.hardwareCRUD.Variant(identifierValue); ok {
						variant.Price, _ = getPeriodFloat(cellVal)
						continue
					}
					hardwareObj.hasError = true
					return &Error{
						ErrTitle: "Variante unbekannt",
						ErrMsg:   fmt.Sprintf("Es konnte keine Variante mit der Ebootis-ID %s gefunden werden", cellVal),
					}
				case "externalArticleNumber":
					if variant, ok := hardwareObj.hardwareCRUD.VariantViaArticleNo(idType); ok { // added method from project
						variant.Price, _ = getPeriodFloat(cellVal)
						continue
					}
					hardwareObj.hasError = true
					return &Error{
						ErrTitle: "Variante unbekannt",
						ErrMsg:   fmt.Sprintf("Es konnte keine Variante mit der MSD Artikelnummer %s gefunden werden", cellVal),
					}
				}
			}

		}

	}

	return nil
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
