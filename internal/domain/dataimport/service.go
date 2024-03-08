package dataimport

import "io"

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

type MappingService interface {
	ReadFile(UploadData) (MappingObject, error)
	WriteMapping(MappingInstruction) (MappingResult, error)
}
