package library

import (
	"encoding/json"
	"errors"
	"fmt"
	"ghoulseek/globals"
	"ghoulseek/metadata"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func ParseGhoulArtist(path string) (metadata.Artist, error) {
	_, stat := os.Stat(path)

	if os.IsNotExist(stat) {
		return metadata.Artist{}, errors.New("Invalid path supplied")
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return metadata.Artist{}, errors.New("Failed to read file data (? what ?)")
	}

	parseData := metadata.Artist{}

	err = json.Unmarshal(data, &parseData)

	if err != nil {
		return metadata.Artist{}, errors.New("Invalid contents, failed to parse JSON")
	}

	return parseData, nil
}

func ReadLibrary() ([]metadata.Artist, error) {
	pattern := filepath.Join(globals.MusicDir, "*/"+globals.ArtistFile)

	files, err := filepath.Glob(pattern)

	if err != nil {
		fmt.Println(err.Error())
		return []metadata.Artist{}, err
	}

	artistList := []metadata.Artist{}

	for _, file := range files {
		fmt.Println(file)

		artist, err := ParseGhoulArtist(file)

		if err != nil {
			fmt.Println("Error parsing artist of path " + color.YellowString(file) + ": " + color.RedString(err.Error()))
		} else {
			artistList = append(artistList, artist)
		}
	}

	return artistList, nil
}
