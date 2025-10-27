package slsk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghoulseek/globals"
	"net/http"

	"github.com/fatih/color"
)

func StartSearch(query string) (SearchResponse, error) {
	url := globals.SlskdEndpoint + "/searches"

	type Params struct {
		SearchText string `json:"searchText"`
	}

	paramData := Params{
		SearchText: query,
	}

	// no need to handle err, it will never fail
	encode, _ := json.Marshal(paramData)

	resp, err := http.Post(url, "application/json", bytes.NewReader(encode))

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

	fmt.Println("Created search request: " + color.GreenString(result.SearchId))

	return result, nil
}
