package library

import (
	"ghoulseek/musicbrainz"
)

func LoadArtist(musicBrainzId string) error {
	artistData, err := musicbrainz.GetArtistFull(musicBrainzId)

	if err != nil {
		return err
	}

	WriteArtistFile(artistData)

	return nil
}
