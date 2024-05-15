package utils_test

import (
	"context"
	"sky-lantern/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHTTPRequest(t *testing.T) {
	ctx := context.Background()

	// Test with a valid URL
	validURL := "https://morristai.github.io/data/60MB_file.txt_chunk_aa"
	data, err := utils.MakeHTTPRequest(ctx, validURL)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// Test with an invalid URL
	invalidURL := "invalid-url"
	data, err = utils.MakeHTTPRequest(ctx, invalidURL)
	assert.Error(t, err)
	assert.Nil(t, data)

	// Test with a URL that returns a non-200 status code
	nonOKURL := "https://morristai.github.io/data/60MB_file.txt_chunk_az"
	data, err = utils.MakeHTTPRequest(ctx, nonOKURL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status code")
	assert.Nil(t, data)

	// TODO: Test with a URL that times out
}

func TestGetFileMeta(t *testing.T) {
	ctx := context.Background()

	// Test with a valid URL
	validURL := "https://morristai.github.io/data/60MB_file.txt_chunk_aa"
	meta, err := utils.GetFileMeta(ctx, validURL)
	assert.NoError(t, err)
	assert.NotZero(t, meta.Size)
	assert.NotEmpty(t, meta.Hash)

	// Test with an invalid URL
	invalidURL := "invalid-url"
	meta, err = utils.GetFileMeta(ctx, invalidURL)
	assert.Error(t, err)
	assert.Zero(t, meta.Size)
	assert.Empty(t, meta.Hash)

	// Test with a URL that returns a non-200 status code
	nonOKURL := "https://morristai.github.io/data/60MB_file.txt_chunk_az"
	meta, err = utils.GetFileMeta(ctx, nonOKURL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status code")
	assert.Zero(t, meta.Size)
	assert.Empty(t, meta.Hash)
}

func TestCalculateFileSha256(t *testing.T) {
	// Test with sample data
	data := []byte("Hello, World!")
	expectedHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	hash, err := utils.CalculateFileSha256(data)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)

	// Test with empty data
	emptyData := []byte{}
	expectedEmptyHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	hash, err = utils.CalculateFileSha256(emptyData)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmptyHash, hash)
}

func TestCalculateFileSha1(t *testing.T) {
	// Test with sample data
	data := []byte("Hello, World!")
	expectedHash := "0a0a9f2a6772942557ab5355d76af442f8f65e01"
	hash, err := utils.CalculateFileSha1(data)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)

	// Test with empty data
	emptyData := []byte{}
	expectedEmptyHash := "da39a3ee5e6b4b0d3255bfef95601890afd80709"
	hash, err = utils.CalculateFileSha1(emptyData)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmptyHash, hash)
}

func TestCalculateFileMD5(t *testing.T) {
	// Test with sample data
	data := []byte("Hello, World!")
	expectedHash := "65a8e27d8879283831b664bd8b7f0ad4"
	hash, err := utils.CalculateFileMD5(data)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)

	// Test with empty data
	emptyData := []byte{}
	expectedEmptyHash := "d41d8cd98f00b204e9800998ecf8427e"
	hash, err = utils.CalculateFileMD5(emptyData)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmptyHash, hash)
}

func TestCreateOutputFile(t *testing.T) {
	// Test with a valid output filename
	validFilename := "output.txt"
	file, err := utils.CreateOutputFile(validFilename)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()

	// Test with a non-existent directory in the output filename
	nestedFilename := "nested/dir/output.txt"
	file, err = utils.CreateOutputFile(nestedFilename)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()

	// Test with an invalid output filename (e.g., a directory instead of a file)
	invalidFilename := "invalid/"
	file, err = utils.CreateOutputFile(invalidFilename)
	assert.Error(t, err)
	assert.Nil(t, file)
}
