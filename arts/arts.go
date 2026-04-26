package arts

import (
	"encoding/json"
	"sort"

	"github.com/charmbracelet/log"
)

type ArtManagerType string

const (
	Kmeans ArtManagerType = "kmeans"
)

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
