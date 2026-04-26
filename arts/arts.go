package arts

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"

	"github.com/charmbracelet/log"
)

type ArtManagerType string

const (
	Kmeans ArtManagerType = "kmeans"
)

type ArtManager interface {
	GetArtPaletteForImage(url.URL) (ArtPalette, error)
}

type CommonArtManager struct {
	logger log.Logger
}

type ArtColor struct {
	ColorHex  string
	Intensity int
}

type ArtPalette struct {
	ArtColors []ArtColor
}

func (artPalette *ArtPalette) AddArtColor(artColor ArtColor) {
	artPalette.ArtColors = append(artPalette.ArtColors, artColor)
	sort.SliceStable(artPalette.ArtColors, func(i, j int) bool { return artPalette.ArtColors[i].Intensity > artPalette.ArtColors[j].Intensity })
}

func NewArtManager(logger *log.Logger, artManagerType ArtManagerType) ArtManager {
	switch artManagerType {
	case Kmeans:
		return NewKmeansArtManager(logger)
	default:
		panic(fmt.Errorf("We don't support any %s based ArtManagers yet!\n", artManagerType))
	}
}

func (artPalette ArtPalette) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"artColors": artPalette.ArtColors,
	})
}

func (artColor ArtColor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"hex":       artColor.ColorHex,
		"intensity": artColor.Intensity,
	})
}
