package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
)

func ReadGzipFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, gzipReader); err != nil {
		return nil, fmt.Errorf("failed to read compressed data: %v", err)
	}

	return buffer.Bytes(), nil
}
