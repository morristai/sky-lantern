package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
)

func MakeHTTPRequest(url string) ([]byte, error) {
	// Support Get requests for now
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type FileMeta struct {
	Size int64
	Hash string
}

func GetFileMeta(url string) (FileMeta, error) {
	resp, err := http.Head(url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Failed to send HEAD request")
		return FileMeta{}, err
	}
	defer resp.Body.Close()

	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Invalid content length")
		return FileMeta{}, err
	}

	hash := resp.Header.Get("Etag")

	log.Debug().Str("url", url).Int64("size", size).Str("hash", hash).Msg("File metadata")

	return FileMeta{Size: size, Hash: hash}, nil
}

func CalculateFileSha256(data []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

func CalculateFileSha1(data []byte) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

func CalculateFileMD5(data []byte) (string, error) {
	hash := md5.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

func CreateOutputFile(outputFilename string) (*os.File, error) {
	outputFile, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			destinationFolder := filepath.Dir(outputFilename)
			if err := os.MkdirAll(destinationFolder, os.ModePerm); err != nil {
				log.Error().Err(err).Msg("Failed to create destination folder")
				return nil, err
			}
			outputFile, err = os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to create output file")
			return nil, err
		}
	}
	return outputFile, nil
}
