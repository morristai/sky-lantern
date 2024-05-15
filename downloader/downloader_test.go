package downloader_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"sky-lantern/downloader"
)

func TestDownloadFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "downloader-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a temporary output file
	outputFilename := filepath.Join(tempDir, "output.txt")

	// Test with valid chunk URLs
	chunkURLs := []string{
		"https://morristai.github.io/data/60MB_file.txt_chunk_aa",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ab",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ac",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ad",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ae",
		"https://morristai.github.io/data/60MB_file.txt_chunk_af",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ag",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ah",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ai",
		"https://morristai.github.io/data/60MB_file.txt_chunk_aj",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ak",
	}
	err = downloader.DownloadFile(chunkURLs, false, outputFilename)
	assert.NoError(t, err)

	// Test with an invalid chunk URL
	invalidChunkURLs := []string{
		"https://morristai.github.io/data/60MB_file.txt_chunk_aa",
		"invalid-url",
		"https://morristai.github.io/data/60MB_file.txt_chunk_ac",
	}
	err = downloader.DownloadFile(invalidChunkURLs, false, outputFilename)
	assert.Error(t, err)

	// Test with persistent chunk disabled
	err = downloader.DownloadFile(chunkURLs, false, outputFilename)
	assert.NoError(t, err)

	// Test with persistent chunk enabled
	err = downloader.DownloadFile(chunkURLs, true, outputFilename)
	assert.NoError(t, err)

	// Verify that the output file exists
	_, err = os.Stat(outputFilename)
	assert.NoError(t, err)
}
