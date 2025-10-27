package library

import (
	"encoding/json"
	"fmt"
	"ghoulseek/globals"
	"ghoulseek/metadata"
	"os"
	"path"

	"github.com/fatih/color"
)

// i don't think this needs to return an error?
func WriteArtistFile(artist metadata.Artist) {
	artistPath := path.Join(globals.MusicDir, artist.Name)

	os.MkdirAll(artistPath, 0755)

	filePath := path.Join(artistPath, globals.ArtistFile)

	data, err := json.Marshal(artist)

	if err != nil {
		fmt.Println("Failed to write artist file: " + color.RedString(err.Error()))
		return
	}

	os.WriteFile(filePath, data, 0755)
}
