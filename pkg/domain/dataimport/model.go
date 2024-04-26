package dataimport

import (
	"fmt"
	"io"

	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/hardware"
)

// Contains file and type of import
type UploadData struct {
	UploadedFile io.Reader
	UploadType   string
	Uuid         string
}

type MappingOptions struct {
	DropdownOptions map[string]string
	TableHeaders    []string
	TableSummary    [][]string
	Uuid            string
}

// Instructions sent to BE after mapping in FE
type MappingInstruction struct {
	Mapping    []MappingObject
	Uuid       string
	UploadType string
}

type MappingObject struct {
	ColIndex     int
	MappingValue string
}

type MappingResult struct {
	SuccessfulRows   int
	UnsuccessfulRows int
	FailedRows       []Error
}

type editedCRUDobj struct {
	hardwareCRUD *hardware.HardwareCRUD
	hasError     bool
}

type Error struct {
	ErrTitle string
	ErrMsg   string
}

func (err *Error) Error() string {
	return fmt.Sprintf("Error: %s", err.ErrMsg)
}

// Returns according column index and identifier type (ebootisId or externalArticleNumber) if either is found.
// Returns false, 0 and "" if no identifier found.
func (mi *MappingInstruction) GetIdentifierIndex() (exists bool, idIndex int, idType string) {
	for i, m := range mi.Mapping {
		switch m.MappingValue {
		case "ebootisId":
			return true, i, "ebootisId"
		case "externalArticleNumber":
			exists = true
			idIndex = i
			idType = "externalArticleNumber"
		}
	}
	return
}
