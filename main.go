package main

import (
	"fmt"
	"ghoulseek/downloader/slsk"
	"ghoulseek/library"
	"ghoulseek/musicbrainz"

	"github.com/fatih/color"
)

func scanTest() {
	libraryScan, err := library.ReadLibrary()

	if err != nil {
		fmt.Println("Error scanning library: " + color.RedString(err.Error()))
		return
	}

	fmt.Println(libraryScan)
}

var testId string = "f6beac20-5dfe-4d1f-ae02-0b0a740aafd6"

func getTest() {

	fmt.Println("Testing artist add with id " + color.GreenString(testId))

	test, err := musicbrainz.GetArtistFull(testId)

	if err != nil {
		fmt.Println("Error getting artist: " + color.RedString(err.Error()))
		return
	}

	fmt.Println(test)

	library.WriteArtistFile(test)
}

func main() {

	library, _ := library.ReadLibrary()

	result, _ := slsk.ReleaseSearch(library[0].Releases[0])

	bestRelease := []slsk.File{}
	var bestScore float64 = -99

	for _, r := range result {
		score := slsk.EvaluateFileList(r.Files)

		if score > float64(bestScore) {
			bestRelease = r.Files
			bestScore = score
		}
	}

	fmt.Println(bestRelease)
}
