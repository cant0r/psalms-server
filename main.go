package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/cant0r/psalms-server/arts"
	"github.com/cant0r/psalms-server/configurations"
	"github.com/cant0r/psalms-server/psalms"
)

var (
	Port             int
	TargetPlayerName string
)

var artPaletteCache map[url.URL]arts.ArtPalette = make(map[url.URL]arts.ArtPalette)

func main() {
	logger := configurations.NewLogger()

	flag.IntVar(&Port, "port", 16666, "On what port the server should listen on?")
	flag.StringVar(&TargetPlayerName, "target-player", "spotify", "Your preferred media player to attach to.")
	flag.Parse()

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/get-playing-psalm", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("Querying psalm metadata from %s", TargetPlayerName)
		psalmist, err := psalms.New(logger, TargetPlayerName)

		if err != nil {
			logger.Fatal("Failed to initialize Psalmist!", "err", err)
			os.Exit(1)
		}

		playingPsalmMetadata, err := psalmist.GetPlayingPsalmMetadata()

		artManager := arts.NewArtManager(logger, arts.Kmeans)
		if artPalette, ok := artPaletteCache[playingPsalmMetadata.ArtUrl]; ok {
			logger.Info("Retrieved from cache!")
			playingPsalmMetadata.ArtPalette = artPalette
		} else {
			playingPsalmMetadata.ArtPalette, err = artManager.GetArtPaletteForImage(playingPsalmMetadata.ArtUrl)

			if err != nil {
				logger.Error("Failed to determine prominent colors for the artwork!", "err", err)
			}

			artPaletteCache[playingPsalmMetadata.ArtUrl] = playingPsalmMetadata.ArtPalette
		}

		if err != nil {
			logger.Error("Failed to retrieve currently playing psalm metadata!", "err", err)
			os.Exit(2)
		}

		logger.Info("Successfully retrieved psalm metadata!")
		json, _ := json.MarshalIndent(playingPsalmMetadata, "", "  ")
		w.Write(json)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", Port), httpMux)

	if err != nil {
		logger.Error("\t", "err", err)
		logger.Fatalf("Failed to start server on %s!", fmt.Sprintf(":%d", Port))
	}
}
