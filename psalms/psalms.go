package psalms

import (
	"encoding/json"
	"net/url"

	"github.com/cant0r/psalms-server/arts"
	"github.com/charmbracelet/log"
)

type Psalmist interface {
	GetPlayingPsalmMetadata() (PsalmMetadata, error)
}

type CommonPsalmist struct {
	logger log.Logger
}

func New(logger *log.Logger, playerName string) (Psalmist, error) {
	psalmist, err := newPsalmist(logger, playerName)

	if err != nil {
		logger.Error("Failed to get any active player!", "err", err)
		return nil, err
	}

	return psalmist, nil
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
