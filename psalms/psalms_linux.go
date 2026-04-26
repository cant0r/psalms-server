//go:build linux

package psalms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/godbus/dbus/v5"
)

type Psalmist struct {
	CommonPsalmist
	mprisPlayerPrettyName string
	mprisPlayerId         string
}

func (psalmist *Psalmist) attachMprisPlayer(conn *dbus.Conn) error {
	var names []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)

	if err != nil {
		return fmt.Errorf("Failed to query owned DBus objects: %s\n", err)
	}

	// Query the player
	var activePlayer string
	for _, name := range names {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
			if strings.Contains(name, psalmist.mprisPlayerPrettyName) {
				activePlayer = name
				break
			}
		}
	}

	if len(activePlayer) == 0 {
		return fmt.Errorf("Found no player with the given name: %s\n", psalmist.mprisPlayerPrettyName)
	}

	psalmist.mprisPlayerId = activePlayer
	psalmist.logger.Infof("Using MRPIS compatible player: %s", psalmist.mprisPlayerId)

	return nil
}

func (psalmist Psalmist) GetPlayingPsalmMetadata() (PsalmMetadata, error) {
	conn, err := dbus.ConnectSessionBus()

	if err != nil {
		return PsalmMetadata{}, fmt.Errorf("Failed to connect to session bus: %s\n", err)
	}
	defer conn.Close()

	err = psalmist.attachMprisPlayer(conn)

	var metadata map[string]dbus.Variant
	call := conn.Object(psalmist.mprisPlayerId, "/org/mpris/MediaPlayer2").Call("org.freedesktop.DBus.Properties.Get", 0, "org.mpris.MediaPlayer2.Player", "Metadata")
	err = call.Store(&metadata)

	if err != nil {
		return PsalmMetadata{}, fmt.Errorf("Failed to get currently playing music: %s\n", err)
	}

	var mprisArtUrl string
	_ = metadata["mpris:artUrl"].Store(&mprisArtUrl)
	artUrl, err := url.Parse(mprisArtUrl)

	if err != nil {
		psalmist.logger.Warn("Failed to parse mpris:artUrl as net.url type!", "mprisArtUrl", artUrl)
		artUrl, _ = url.Parse("http://localhost:666/")
	}

	return PsalmMetadata{
		Title:  metadata["xesam:title"].String(),
		Album:  metadata["xesam:album"].String(),
		ArtUrl: *artUrl,
	}, nil
}

func New(logger *log.Logger, playerName string) (Psalmist, error) {
	return Psalmist{
		CommonPsalmist:        CommonPsalmist{logger: *logger},
		mprisPlayerPrettyName: strings.ToLower(playerName),
	}, nil
}
