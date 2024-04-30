package dataimport

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/Filipza/excel-mapping-tool/internal/domain/v1/tariff"
	"github.com/Filipza/excel-mapping-tool/internal/settings"
	"github.com/stretchr/testify/assert"
)

func TestCustomError(t *testing.T) {
	customErr := Error{
		ErrTitle: "custom error title",
		ErrMsg:   "custom error message",
	}

	errorMsg := customErr.Error()
	assert.Equal(t, errorMsg, "Error: custom error message")
}

func TestReadFilePositive(t *testing.T) {
	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "stocks",
	}

	svc := mappingService{}

	result, err := svc.ReadFile(mockUploadData)

	assert.NotNil(t, result, "result should not be nil")
	assert.NoError(t, err, "function should not return an error")
	assert.NotEmpty(t, result, "result should not be empty")
}

func TestDirCreationPositive(t *testing.T) {
	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "stocks",
		Uuid:         "3fc7522d-ed25-40df-9972-333ba8aea504",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.NoError(t, err, "valid UUID should be accepted for directory creation")
	assert.FileExists(t, "/tmp/"+mockUploadData.Uuid+"/data.xlsx")
}

func TestDirCreationNegative(t *testing.T) {
	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "stocks",
		Uuid:         "longdirname" + strings.Repeat("a", 300),
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.Error(t, err, "invalid uuid should not be accepted")
}

func TestReadFileNegative(t *testing.T) {
	file, err := os.ReadFile("../../../test/corrupted.xlsx")
	if err != nil {
		t.Fatalf("loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "stocks",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.Error(t, err, "corrupt file was correctly denied")
}

func TestReadFileEmpty(t *testing.T) {
	emptyFile := bytes.NewReader([]byte{})

	mockUploadData := &UploadData{
		UploadedFile: emptyFile,
		UploadType:   "stocks",
	}

	svc := mappingService{}

	_, err := svc.ReadFile(mockUploadData)

	assert.Error(t, err, "empty file should be rejected")

	if customErr, ok := err.(*Error); ok {
		assert.Equal(t, "Parsingfehler", customErr.ErrTitle)
	}
}

func TestReadFileEmptyHeaders(t *testing.T) {
	file, err := os.ReadFile("../../../test/empty_headers.xlsx")
	if err != nil {
		t.Fatalf("loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "stocks",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	if customErr, ok := err.(*Error); ok {
		assert.Equal(t, customErr.ErrTitle, "Leerzeile")
	}

	assert.Error(t, err)
}

func TestReadFileNoUploadType(t *testing.T) {
	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.Error(t, err, "no upload type given")
}

func TestReadFileUploadType(t *testing.T) {
	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("loading test .xlsx failed: %v", err)
	}

	uploadTypes := []string{"tariff", "hardware", "stocks"}

	for _, v := range uploadTypes {
		mockUploadData := &UploadData{
			UploadedFile: bytes.NewReader(file),
			UploadType:   v,
		}

		svc := mappingService{}

		_, err = svc.ReadFile(mockUploadData)

		assert.NoError(t, err, "function should no return an error")
	}
}
func TestReadFileUnknownUploadType(t *testing.T) {
	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "bananen",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.Error(t, err, "UploadData should contain valid uploadType")
}

func TestGetEbootisIndexPositive(t *testing.T) {
	mi := &MappingInstruction{
		Uuid: "1e1133c1-65cf-46f6-a246-6049234d3447",
		Mapping: []MappingObject{
			{ColIndex: 0, MappingValue: "externalArticleNumber"},
			{ColIndex: 1, MappingValue: "ebootisId"},
			{ColIndex: 2, MappingValue: "pibLink"},
			{ColIndex: 3, MappingValue: "supplierWkz"},
		},
	}

	exists, idIndex, idtype := mi.GetIdentifierIndex()

	assert.Equal(t, exists, true, "exists should equal true")
	assert.Equal(t, idIndex, 1, "idIndex should equal 0")
	assert.Equal(t, idtype, "ebootisId", "idtype should equal 'ebootisId'")
}

func TestGetExternalArticleNumberIndexPositive(t *testing.T) {
	mi := &MappingInstruction{
		Uuid: "2e1133c1-65cf-46f6-a246-6049234d3448",
		Mapping: []MappingObject{
			{ColIndex: 0, MappingValue: "leadType"},
			{ColIndex: 1, MappingValue: "externalArticleNumber"},
			{ColIndex: 2, MappingValue: "pibLink"},
			{ColIndex: 3, MappingValue: "supplierWkz"},
		},
	}

	exists, idIndex, idtype := mi.GetIdentifierIndex()

	assert.Equal(t, exists, true, "exists should equal true")
	assert.Equal(t, idIndex, 1, "idIndex should equal 0")
	assert.Equal(t, idtype, "externalArticleNumber", "idtype should equal 'externalArticleNumber'")
}

func TestGetIdentifierIndexNegative(t *testing.T) {
	mi := &MappingInstruction{
		Uuid: "3e1133c1-65cf-46f6-a246-6049234d3449",
		Mapping: []MappingObject{
			{ColIndex: 0, MappingValue: "leadType"},
			{ColIndex: 1, MappingValue: "provision"},
			{ColIndex: 2, MappingValue: "pibLink"},
			{ColIndex: 3, MappingValue: "supplierWkz"},
		},
	}

	exists, idIndex, idtype := mi.GetIdentifierIndex()

	assert.Equal(t, exists, false, "exists should equal true")
	assert.Equal(t, idIndex, 0, "idIndex should equal 0")
	assert.Equal(t, idtype, "", "idtype should equal ''")
}

func TestWriteMappingGetIdentifierNegative(t *testing.T) {

	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "stocks",
	}

	svc := mappingService{}

	result, _ := svc.ReadFile(mockUploadData)

	mi := &MappingInstruction{
		Uuid: result.Uuid,
		Mapping: []MappingObject{
			{ColIndex: 0, MappingValue: "pibLink"},
			{ColIndex: 1, MappingValue: "supplierWkz"},
		},
	}

	_, err = svc.WriteMapping(mi)

	assert.Error(t, err, "MappingInstruction should contain either 'EbootisID' or 'externalArticleNumber' as MappingValue")
}

// func TestWriteMappingOpenXlsxNegative(t *testing.T) {

// 	file, err := os.ReadFile("../../../test/positive.xlsx")
// 	if err != nil {
// 		t.Fatalf("Loading test .xlsx failed: %v", err)
// 	}

// 	mockUploadData := &UploadData{
// 		UploadedFile: bytes.NewReader(file),
// 		UploadType:   "hardware",
// 	}

// 	svc := mappingService{}

// 	svc.ReadFile(mockUploadData)

// 	mi := &MappingInstruction{
// 		Uuid: "8cae326f-d3de-45d4-8bb8-ded181f44a0e",
// 		Mapping: []MappingObject{
// 			{ColIndex: 0, MappingValue: "pibUrl"},
// 			{ColIndex: 1, MappingValue: "wkz"},
// 		},
// 	}

// 	_, err = svc.WriteMapping(mi)

// 	assert.NoError(t, err, "function should not return an error")
// }

type crudMock[T, L any] struct {
	list   func(...settings.Option) ([]*L, error)
	create func(*T, ...settings.Option) (*T, error)
	read   func(string, ...settings.Option) (*T, error)
	update func(string, *T, ...settings.Option) (*T, error)
	delete func(string, ...settings.Option) (*T, error)
}

func (svc *crudMock[T, L]) List(opts ...settings.Option) ([]*L, error) {
	if svc.list != nil {
		return svc.list(opts...)
	}
	return nil, errors.New("not implememnted")
}

func (svc *crudMock[T, L]) Create(t *T, opts ...settings.Option) (*T, error) {
	if svc.create != nil {
		return svc.create(t, opts...)
	}
	return nil, errors.New("not implememnted")
}

func (svc *crudMock[T, L]) Read(id string, opts ...settings.Option) (*T, error) {
	if svc.read != nil {
		return svc.read(id, opts...)
	}
	return nil, errors.New("not implememnted")
}

func (svc *crudMock[T, L]) Update(id string, t *T, opts ...settings.Option) (*T, error) {
	if svc.update != nil {
		return svc.update(id, t, opts...)
	}
	return nil, errors.New("not implememnted")
}

func (svc *crudMock[T, L]) Delete(id string, opts ...settings.Option) (*T, error) {
	if svc.delete != nil {
		return svc.delete(id, opts...)
	}
	return nil, errors.New("not implememnted")
}

func TestWriteMapping(t *testing.T) {

	file, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(file),
		UploadType:   "tariff",
	}

	tariffAdapter := &crudMock[tariff.TariffCRUD, tariff.TariffLookup]{}
	svc := &mappingService{
		tariffAdapter: tariffAdapter,
	}

	result, err := svc.ReadFile(mockUploadData)

	mi := &MappingInstruction{
		Uuid: result.Uuid,
		Mapping: []MappingObject{
			{ColIndex: 0, MappingValue: "ebootisId"},
			{ColIndex: 1, MappingValue: "externalArticleNumber"},
			{ColIndex: 2, MappingValue: "pibUrl"},
			{ColIndex: 3, MappingValue: "wkz"},
		},
		UploadType: "tariff",
	}

	_, err = svc.WriteMapping(mi)

	assert.Error(t, err, "function should not return an error")

	tariffAdapter.list = func(o ...settings.Option) ([]*tariff.TariffLookup, error) {
		return []*tariff.TariffLookup{}, nil
	}

	// tariffAdapter.read = func(s string, o ...settings.Option) (*tariff.TariffCRUD, error) {

	// }
}
