package psalms

import (
	"encoding/json"
	"net/url"

	"github.com/cant0r/psalms-server/arts"
	"github.com/charmbracelet/log"
)

type CommonPsalmist struct {
	logger log.Logger
}

type PsalmMetadata struct {
	Title      string
	Album      string
	ArtUrl     url.URL
	ArtPalette arts.ArtPalette
}

func (psalmMetadata PsalmMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"title":      psalmMetadata.Title,
		"album":      psalmMetadata.Album,
		"artUrl":     psalmMetadata.ArtUrl.String(),
		"artPalette": psalmMetadata.ArtPalette,
	})
}
