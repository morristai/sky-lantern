package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMakeHTTPRequest(t *testing.T) {
	// Test case 1: Successful GET request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	data, err := MakeHTTPRequest(server.URL)
	if err != nil {
		t.Errorf("MakeHTTPRequest failed: %v", err)
	}
	if string(data) != "test response" {
		t.Errorf("Unexpected response: %s", string(data))
	}

	// Test case 2: Invalid URL
	_, err = MakeHTTPRequest("invalid_url")
	if err == nil {
		t.Error("MakeHTTPRequest should have failed with an invalid URL")
	}
}

func TestGetFileMeta(t *testing.T) {
	// Test case 1: Retrieve file size and hash
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1024")
		w.Header().Set("Etag", "test_hash")
	}))
	defer server.Close()

	meta, err := GetFileMeta(server.URL)
	if err != nil {
		t.Errorf("GetFileMeta failed: %v", err)
	}
	if meta.Size != 1024 {
		t.Errorf("Unexpected file size: %d", meta.Size)
	}
	if meta.Hash != "test_hash" {
		t.Errorf("Unexpected file hash: %s", meta.Hash)
	}

	// Test case 2: Invalid URL
	_, err = GetFileMeta("invalid_url")
	if err == nil {
		t.Error("GetFileMeta should have failed with an invalid URL")
	}
}

func TestCalculateFileSha256(t *testing.T) {
	data := []byte("test data")
	expectedHash := "916f0027a575074ce72a331777c3478d6513f786a591bd892da1a577bf2335f9"
	hash, err := CalculateFileSha256(data)
	if err != nil {
		t.Errorf("CalculateFileSha256 failed: %v", err)
	}
	if hash != expectedHash {
		t.Errorf("Unexpected hash: %s", hash)
	}
}

func TestCalculateFileSha1(t *testing.T) {
	data := []byte("test data")
	expectedHash := "f48dd853820860816c75d54d0f584dc863327a7c"
	hash, err := CalculateFileSha1(data)
	if err != nil {
		t.Errorf("CalculateFileSha1 failed: %v", err)
	}
	if hash != expectedHash {
		t.Errorf("Unexpected hash: %s", hash)
	}
}

func TestCalculateFileMD5(t *testing.T) {
	data := []byte("test data")
	expectedHash := "eb733a00c0c9d336e65691a37ab54293"
	hash, err := CalculateFileMD5(data)
	if err != nil {
		t.Errorf("CalculateFileMD5 failed: %v", err)
	}
	if hash != expectedHash {
		t.Errorf("Unexpected hash: %s", hash)
	}
}

func TestCreateOutputFile(t *testing.T) {
	// Test case 1: Create output file
	outputFilename := "test_output.txt"
	outputFile, err := CreateOutputFile(outputFilename)
	if err != nil {
		t.Errorf("CreateOutputFile failed: %v", err)
	}
	defer outputFile.Close()
	// Check that the output file was created
	_, err = os.Stat(outputFilename)
	if os.IsNotExist(err) {
		t.Errorf("Output file '%s' was not created", outputFilename)
	}
	os.Remove(outputFilename)

	// Test case 2: Create output file in non-existent directory
	outputFilename = "nonexistent/test_output.txt"
	outputFile, err = CreateOutputFile(outputFilename)
	if err != nil {
		t.Errorf("CreateOutputFile failed: %v", err)
	}
	defer outputFile.Close()
	// Check that the output file was created
	_, err = os.Stat(outputFilename)
	if os.IsNotExist(err) {
		t.Errorf("Output file '%s' was not created", outputFilename)
	}
	os.RemoveAll("nonexistent")
}
