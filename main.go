package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"

	"sky-lantern/downloader"
)

func main() {
	keepChunkFiles := flag.Bool("keep-chunks", false, "Keep the downloaded chunk files")
	outputFilename := flag.String("output", "output.txt", "Output filename")
	debugMode := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	if *debugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	chunkURLs := flag.Args()
	if len(chunkURLs) < 2 {
		log.Error().Msg("Usage: go run main.go [flags] <chunk1_url> <chunk2_url> <chunk3_url> ...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := downloader.DownloadFile(chunkURLs, *keepChunkFiles, *outputFilename)
	if err != nil {
		log.Error().Err(err).Msg("Error")
		os.Exit(1)
	}

	log.Info().Msg("File downloaded successfully!")
}
