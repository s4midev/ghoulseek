package main

import (
	"fmt"
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

var testId string = "4e2b6413-ebc1-4ade-bca7-6d078e476c97"

func getTest() {

	fmt.Println("Testing artist add with id " + color.GreenString(testId))

	test, err := musicbrainz.GetArtistFull(testId)

	if err != nil {
		fmt.Println("Error getting artist: " + color.RedString(err.Error()))
		return
	}

	fmt.Println(test)
}

func main() {
	scanTest()
}
