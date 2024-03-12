package dataimport

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFileNegative(t *testing.T) {
	xlsx, err := os.ReadFile("../../../test/png-test.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(xlsx),
	}

	svc := mappingService{}

	_, err = svc.ReadFile(mockUploadData)

	assert.Error(t, err, "Corrupt file was correctly denied.")
}

func TestReadFilePositive(t *testing.T) {
	xlsx, err := os.ReadFile("../../../test/test.xlsx")
	if err != nil {
		t.Fatalf("Loading test .xlsx failed: %v", err)
	}

	mockUploadData := &UploadData{
		UploadedFile: bytes.NewReader(xlsx),
	}

	svc := mappingService{}

	result, err := svc.ReadFile(mockUploadData)

	assert.NotNil(t, result, "Result should not be nil")
	assert.NoError(t, err, "ReadFile should not return an error")
	assert.NotEmpty(t, result, "Result should not be empty")
}
