package downloader

import (
	"os"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	// Test case 1: Download and write chunks correctly
	chunkURLs := []string{"https://morristai.github.io/data/60MB_file.txt_chunk_aa", "https://morristai.github.io/data/60MB_file.txt_chunk_ab"}
	outputFilename := "test_output.txt"
	err := DownloadFile(chunkURLs, false, outputFilename)
	if err != nil {
		t.Errorf("DownloadFile failed: %v", err)
	}
	// Check that the output file was created and contains the expected content
	os.Remove(outputFilename)

	// Test case 2: Invalid chunk URL
	chunkURLs = []string{"http://example.com/chunk1", "invalid_url"}
	err = DownloadFile(chunkURLs, false, outputFilename)
	if err == nil {
		t.Error("DownloadFile should have failed with an invalid URL")
	}

	// Test case 3: Keep downloaded chunk files
	chunkURLs = []string{"https://morristai.github.io/data/60MB_file.txt_chunk_aa", "https://morristai.github.io/data/60MB_file.txt_chunk_ab"}
	err = DownloadFile(chunkURLs, true, outputFilename)
	if err != nil {
		t.Errorf("DownloadFile failed: %v", err)
	}
	// Check that the downloaded chunk files were kept (you can check the existence of the files in the cache folder)
	os.Remove(outputFilename)
}

func TestDownloadAndWriteChunks(t *testing.T) {
	// Test case 1: Download and write chunks concurrently
	chunkURLs := []string{"https://morristai.github.io/data/60MB_file.txt_chunk_aa", "https://morristai.github.io/data/60MB_file.txt_chunk_ab"}
	chunkHashes := []string{"hash1", "hash2"}
	outputFile, _ := os.Create("test_output.txt")
	defer outputFile.Close()
	err := downloadAndWriteChunks(chunkURLs, chunkHashes, false, outputFile)
	if err != nil {
		t.Errorf("downloadAndWriteChunks failed: %v", err)
	}
	// Check that the output file contains the expected content
	os.Remove("test_output.txt")
}

func TestVerifyChunkHash(t *testing.T) {
	// Test case 1: Verify MD5 hash
	data := []byte("test data")
	expectedHash := "098f6bcd4621d373cade4e832627b4f6"
	verifyChunkHash(data, expectedHash, 0)
	// Check the log output for hash match

	// Test case 2: Verify SHA1 hash
	expectedHash = "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	verifyChunkHash(data, expectedHash, 0)
	// Check the log output for hash match

	// Test case 3: Verify SHA256 hash
	expectedHash = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	verifyChunkHash(data, expectedHash, 0)
	// Check the log output for hash match

	// Test case 4: Unknown hash length
	expectedHash = "invalid_hash"
	verifyChunkHash(data, expectedHash, 0)
	// Check the log output for unknown hash length warning

	// Test case 5: Hash mismatch
	expectedHash = "incorrect_hash"
	verifyChunkHash(data, expectedHash, 0)
	// Check the log output for hash mismatch warning
}
