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
	Mapping []MappingObject
	Uuid    string
}

type MappingObject struct {
	Index        int
	MappingValue string
}

func (mi *MappingInstruction) GetIdentifierIndex() (exists bool, index int, idType string) {
	for i, m := range mi.Mapping {
		switch m.MappingValue {
		case "EbootisId":
			return true, i, "EbootisId"
		case "externalArticleNumber":
			exists = true
			index = i
			idType = "externalArticleNumber"
		}
	}
	return exists, index, idType
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
