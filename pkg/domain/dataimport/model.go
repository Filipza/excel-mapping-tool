package dataimport

import (
	"fmt"
	"io"
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

type Error struct {
	ErrTitle string
	ErrMsg   string
}

func (err *Error) Error() string {
	return fmt.Sprintf("Error: %s", err.ErrMsg)
}

// Returns according column index and identifier type (ebootisId or externalArticleNumber) if either is found.
// Returns false, 0 and "" if none identifier is found.
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
