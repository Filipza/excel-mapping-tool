package dataimport

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFilePositive(t *testing.T) {
	xlsx, err := os.ReadFile("../../../test/positive.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(xlsx),
		UploadType:   "wkz",
	}

	svc := mappingService{}

	result, err := svc.ReadFile(mockUploadData)

	assert.NotNil(t, result, "Result should not be nil")
	assert.NoError(t, err, "ReadFile should not return an error")
	assert.NotEmpty(t, result, "Result should not be empty")
}

func TestReadFileNegative(t *testing.T) {
	xlsx, err := os.ReadFile("../../../test/corrupted.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(xlsx),
		UploadType:   "wkz",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.Error(t, err, "Corrupt file was correctly denied")
}

func TestReadFileEmpty(t *testing.T) {
	emptyFile := bytes.NewReader([]byte{})

	mockUploadData := &UploadData{
		UploadedFile: emptyFile,
		UploadType:   "wkz",
	}

	svc := mappingService{}

	_, err := svc.ReadFile(mockUploadData)

	assert.Error(t, err, "Empty file should be rejected")

	if customErr, ok := err.(*Error); ok {
		assert.Equal(t, "Parsingfehler", customErr.ErrTitle)
	}
}

func TestReadFileEmptyHeaders(t *testing.T) {
	xlsx, err := os.ReadFile("../../../test/empty_headers.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(xlsx),
		UploadType:   "wkz",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	if customErr, ok := err.(*Error); ok {
		assert.Equal(t, customErr.ErrTitle, "Leerzeile")
	}

	assert.Error(t, err)
}

func TestReadFileNoUploadType(t *testing.T) {
	xlsx, err := os.ReadFile("../../../test/corrupted.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(xlsx),
		UploadType:   "",
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.Error(t, err, "Corrupt file was correctly denied")
}
