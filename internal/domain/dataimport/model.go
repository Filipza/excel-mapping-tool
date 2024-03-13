package dataimport

import (
	"fmt"
	"io"
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

// {
// 	"mapping": [
// 	  {"index": 1, "mappingValue": "value1"},
// 	  {"index": 2, "mappingValue": "value2"},
// 	  {"index": 3, "mappingValue": "value3"}
// 	]
// }

type MappingInstructionV2 struct {
	MappingV2 map[int]string
}

// {
// 	"mappingV2": {
// 	  "1": "value1",
// 	  "2": "value2",
// 	  "3": "value3"
// 	}
// }

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
