package arts

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"net/http"
	"net/url"

	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/charmbracelet/log"
)

const (
	NumberOfProminentColors = 10
	CroppedImageWidth       = 100
)

type KmeansArtManager struct {
	CommonArtManager
}

func (manager KmeansArtManager) loadImageOverUrl(artworkUrl url.URL) (image.Image, error) {
	manager.logger.Infof("Load artwork from %s\n", artworkUrl.String())

	response, err := http.Get(artworkUrl.String())
	if err != nil || !strings.HasPrefix(response.Status, "2") {
		manager.logger.Error("Failed to download artwork!", "err", err)
		return nil, err
	}

	defer response.Body.Close()

	image, _, err := image.Decode(response.Body)
	if err != nil {
		manager.logger.Error("Failed to parse artwork as image.Image!", "err", err)
		return nil, err
	}

	return image, nil
}

func (manager KmeansArtManager) GetArtPaletteForImage(artworkUrl url.URL) (ArtPalette, error) {
	image, err := manager.loadImageOverUrl(artworkUrl)
	if err != nil {
		return ArtPalette{}, err
	}

	manager.logger.Info("Determine prominent colors.")
	prominentColors, err := prominentcolor.KmeansWithAll(NumberOfProminentColors, image, prominentcolor.ArgumentNoCropping, CroppedImageWidth, prominentcolor.GetDefaultMasks())
	if err != nil {
		manager.logger.Fatal("Oof. Failed to get prominent colors", "err", err)
		return ArtPalette{}, err
	}

	artPalette := ArtPalette{}
	for _, prominentColor := range prominentColors {
		colorHex := fmt.Sprintf("#%02x%02x%02x", prominentColor.Color.R, prominentColor.Color.G, prominentColor.Color.B)

		artPalette.AddArtColor(ArtColor{
			ColorHex:  colorHex,
			Intensity: prominentColor.Cnt,
		})
	}

	return artPalette, nil
}

func NewKmeansArtManager(logger *log.Logger) ArtManager {
	return KmeansArtManager{
		CommonArtManager: CommonArtManager{logger: *logger},
	}
}
