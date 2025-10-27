package slsk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghoulseek/globals"
	"ghoulseek/metadata"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

type WsMessage struct {
	Type      string           `json:"type"`
	Target    string           `json:"target"`
	Arguments []SearchResponse `json:"arguments"`
}

type Session struct {
	ConnectionId    string `json:"connectionId"`
	ConnectionToken string `json:"connectionToken"`
}

func NegotiateSession() Session {
	resp, err := http.Post(globals.SlskdBase+"hub/search/negotiate?negotiateVersion=1", "application/json", strings.NewReader(""))
	if err != nil {
		fmt.Println(err.Error())
		return Session{}
	}
	defer resp.Body.Close()

	var result Session

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err.Error())
		return Session{}
	}

	fmt.Println("negotiated session")
	fmt.Println(result)
	return result
}

func StartSearch(query string) (SearchResponse, error) {
	urlToUse := globals.SlskdEndpoint + "/searches"

	paramData := struct {
		SearchText string `json:"searchText"`
	}{
		SearchText: query,
	}

	encode, _ := json.Marshal(paramData)

	resp, err := http.Post(urlToUse, "application/json", bytes.NewReader(encode))
	if err != nil {
		fmt.Println(err.Error())
		return SearchResponse{}, err
	}
	defer resp.Body.Close()

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err.Error())
		return SearchResponse{}, err
	}

	// i have to manually delay.. the horrors

	time.Sleep(2 * time.Second)
	return result, nil
}

func GetResponses(id string) ([]SearchResponses, error) {
	url := globals.SlskdEndpoint + "/searches/" + id + "/responses"

	resp, err := http.Get(url)

	fmt.Println(resp.Status)

	if err != nil {
		fmt.Println(err.Error())
		return []SearchResponses{}, err
	}

	defer resp.Body.Close()

	var result []SearchResponses

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err.Error())
		return []SearchResponses{}, err
	}

	if len(result) < 100 {
		return result, nil
	}

	return result[0:100], nil
}

func FullSearch(query string) ([]SearchResponses, error) {
	search, err := StartSearch(query)

	if err != nil {
		fmt.Println("Failed to search slskd: " + color.RedString(err.Error()))
		return []SearchResponses{}, err
	}

	fmt.Println("Successfully searched: " + color.GreenString(search.SearchId))

	results, err := GetResponses(search.SearchId)

	if err != nil {
		fmt.Println("Failed to get search responses: " + color.RedString(err.Error()))
		return []SearchResponses{}, err
	}

	return results, nil
}

func ReleaseSearch(release metadata.Release) ([]SearchResponses, error) {
	result, err := FullSearch(release.ArtistName + " - " + release.Title)

	return result, err
}
