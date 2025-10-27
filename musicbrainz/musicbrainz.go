package musicbrainz

import (
	"encoding/json"
	"fmt"
	"ghoulseek/metadata"
	"net/http"

	"github.com/fatih/color"
)

func GetArtistReleases(id string, artistName string) ([]metadata.Release, error) {
	url := "https://musicbrainz.org/ws/2/release-group/?artist=" + id + "&fmt=json"

	resp, err := http.Get(url)
	if err != nil {
		return []metadata.Release{}, err
	}
	defer resp.Body.Close()

	var result struct {
		ReleaseGroups []ReleaseGroup `json:"release-groups"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []metadata.Release{}, err
	}

	releases := []metadata.Release{}

	for _, r := range result.ReleaseGroups {
		releases = append(releases, metadata.Release{
			Title:         r.Title,
			Information:   "",
			ReleaseDate:   r.FirstReleaseDate,
			Tracks:        []metadata.Track{},
			MusicBrainzId: r.ID,
			ReleaseType:   r.PrimaryType,
			ArtistName:    artistName,
		})
	}

	return releases, nil
}

func GetArtistBase(id string) (metadata.Artist, error) {
	url := "https://musicbrainz.org/ws/2/artist/" + id + "?fmt=json"

	resp, err := http.Get(url)
	if err != nil {
		return metadata.Artist{}, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return metadata.Artist{}, err
	}

	return metadata.Artist{
		Name:          result["name"].(string),
		MusicBrainzId: result["id"].(string),
		Releases:      []metadata.Release{},
	}, nil
}

func GetArtistFull(id string) (metadata.Artist, error) {
	artistData, err := GetArtistBase(id)

	if err != nil {
		fmt.Println("Error getting raw artist " + color.YellowString(id) + ": " + color.RedString(err.Error()))
		return metadata.Artist{}, err
	}

	artistReleases, err := GetArtistReleases(id, artistData.Name)

	if err != nil {
		fmt.Println("Error getting releases for " + color.YellowString(id) + ": " + color.RedString(err.Error()))
		return metadata.Artist{}, err
	}

	for _, r := range artistReleases {
		artistData.Releases = append(artistData.Releases, r)
	}

	return artistData, nil
}
