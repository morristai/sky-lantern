package downloader

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"sync"

	"sky-lantern/utils"
)

const chunkFolder = "./cache"

func DownloadFile(chunkURLs []string, persistentChunk bool, outputFilename string) error {
	numChunks := len(chunkURLs)
	chunkHashes := make([]string, numChunks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type result struct {
		index int
		meta  utils.FileMeta
		err   error
	}

	resultCh := make(chan result, numChunks)

	// NOTE: This operation simulates file info requests for each chunk, which seems to be unnecessary in this project.
	// But in real-world scenarios, some download behaviors may differ based on the file info.
	for i, chunkURL := range chunkURLs {
		go func(i int, chunkURL string) {
			chunkMeta, err := utils.GetFileMeta(ctx, chunkURL)
			resultCh <- result{index: i, meta: chunkMeta, err: err}
		}(i, chunkURL)
	}

	for i := 0; i < numChunks; i++ {
		select {
		case res := <-resultCh:
			if res.err != nil {
				log.Error().Err(res.err).Int("chunk", res.index).Msg("Failed to get size of chunk")
				cancel()
				return res.err
			}
			chunkHashes[res.index] = res.meta.Hash
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	outputFile, err := utils.CreateOutputFile(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = downloadAndWriteChunks(ctx, chunkURLs, chunkHashes, persistentChunk, outputFile)
	if err != nil {
		return err
	}

	return nil
}

// TODO: outputFile may use interface
func downloadAndWriteChunks(ctx context.Context, chunkURLs []string, chunkHashes []string, persistentChunk bool, outputFile *os.File) error {
	numChunks := len(chunkURLs)
	var wg sync.WaitGroup
	// NOTE: like Arc<Mutex<Vec<Bytes>>> in Rust but LOCK-FREE. Can use channel to guarantee thread-safe. (But in this case, all i are unique)
	chunkData := make([][]byte, numChunks)
	errChan := make(chan error, numChunks)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i, chunkURL := range chunkURLs {
		wg.Add(1)
		go func(i int, chunkURL string) {
			defer wg.Done()

			chunkFilename := filepath.Base(chunkURL)
			chunkFilePath := filepath.Join(chunkFolder, chunkFilename)

			if persistentChunk {
				if data, err := os.ReadFile(chunkFilePath); err == nil {
					chunkData[i] = data
					return
				}
			}

			log.Debug().Str("url", chunkURL).Int("chunk", i).Msg("Downloading chunk")

			data, err := utils.MakeHTTPRequest(ctx, chunkURL)
			if err != nil {
				log.Error().Err(err).Int("chunk", i).Msg("Failed to download chunk")
				errChan <- fmt.Errorf("failed to download chunk %d: %v", i, err)
				cancel() // Cancel the context to stop all other goroutines
				return
			}

			chunkData[i] = data

			if persistentChunk {
				if err := os.WriteFile(chunkFilePath, data, 0644); err != nil {
					log.Error().Err(err).Str("chunk", chunkFilename).Msg("Error saving chunk file")
				}
			}
		}(i, chunkURL)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			log.Error().Err(err).Msg("Error occurred during chunk download")
			return err
		}
	}

	for i, data := range chunkData {
		if data == nil {
			continue
		}

		if chunkHashes[i] != "" {
			verifyChunkHash(data, chunkHashes[i], i)
		}

		_, err := outputFile.Write(data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to write chunk data")
			return err
		}
	}

	return nil
}

func verifyChunkHash(data []byte, expectedHash string, chunkIndex int) {
	var (
		hash string
		err  error
	)

	// Note: rough check, should be improved
	switch len(expectedHash) {
	// MD5
	case 32:
		hash, err = utils.CalculateFileMD5(data)
	// SHA1
	case 40:
		hash, err = utils.CalculateFileSha1(data)
	// SHA256
	case 64:
		hash, err = utils.CalculateFileSha256(data)
	default:
		log.Warn().Int("chunk", chunkIndex).Str("hash", expectedHash).Msg("Unknown hash length")
		return
	}

	if err != nil {
		log.Warn().Err(err).Int("chunk", chunkIndex).Msg("Failed to calculate hash")
		return
	}

	if hash != expectedHash {
		log.Warn().Int("chunk", chunkIndex).Str("expected", expectedHash).Str("actual", hash).Msg("Hash mismatch")
		return
	}

	log.Info().Int("chunk", chunkIndex).Msg("Hash matched")
}
