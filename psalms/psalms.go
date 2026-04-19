package psalms

import (
	"encoding/json"
	"net/url"

	"github.com/charmbracelet/log"
)

type Psalmist interface {
	GetPlayingPsalmMetadata() (PsalmMetadata, error)
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
	title  string
	album  string
	artUrl url.URL
}

func (psalmMetadata PsalmMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"title":  psalmMetadata.title,
		"album":  psalmMetadata.album,
		"artUrl": psalmMetadata.artUrl.String(),
	})
}
