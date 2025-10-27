package slsk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghoulseek/globals"
	"ghoulseek/metadata"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

var (
	socket        *websocket.Conn
	SocketRunning bool
)

type WsMessage struct {
	Type      string           `json:"type"`
	Target    string           `json:"target"`
	Arguments []SearchResponse `json:"arguments"`
}

func InitSocket() {
	if SocketRunning {
		return
	}

	u := url.URL{Scheme: "ws", Host: globals.SlskdIp, Path: "/hub/search"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	socket = conn
	SocketRunning = true

	go func() {
		defer func() {
			SocketRunning = false
			if socket != nil {
				socket.Close()
			}
		}()

		for {
			_, msg, err := socket.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) ||
					strings.Contains(err.Error(), "use of closed network connection") {
					log.Println("socket closed")
					return
				}
				log.Println("read error:", err)
				return
			}

			if strings.Contains(string(msg), "\x1e") {
				continue
			}

			parseMessage := WsMessage{}
			if err := json.Unmarshal(msg, &parseMessage); err != nil {
				log.Println("failed to unmarshal ws message:", err, "raw:", string(msg))
				continue
			}

			if err := json.Unmarshal(msg, &parseMessage); err != nil {
				log.Println("failed to unmarshal ws message:", err, "raw:", string(msg))
				continue
			}
		}
	}()
}

// this is the worst piece of code i have ever written. but it works. so. you know.
func WaitForMessage(checker func(WsMessage) bool) bool {
	if !SocketRunning {
		log.Println("socket is not running")
		return false
	}

	resultCh := make(chan bool)

	go func() {
		for {
			_, msg, err := socket.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) ||
					strings.Contains(err.Error(), "use of closed network connection") {
					fmt.Println("socket closed")
					resultCh <- false
					return
				}
				resultCh <- false
				return
			}

			parseMessage := WsMessage{}
			if err := json.Unmarshal(msg, &parseMessage); err != nil {
				fmt.Println(err.Error())
				continue
			}

			if checker(parseMessage) {
				resultCh <- true
				return
			}
		}
	}()

	return <-resultCh
}

func StartSearch(query string) (SearchResponse, error) {
	if SocketRunning == false {
		InitSocket()
	}

	url := globals.SlskdEndpoint + "/searches"

	paramData := struct {
		SearchText string `json:"searchText"`
	}{
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

	ok := WaitForMessage(func(wm WsMessage) bool {
		return wm.Arguments[0].IsComplete && wm.Arguments[0].SearchId == result.SearchId
	})

	if !ok {
		fmt.Println(color.RedString("oh well, we'll probably be fine"))
	}

	return result, nil
}

func GetResponses(id string) ([]SearchResponses, error) {
	if SocketRunning == false {
		InitSocket()
	}

	url := globals.SlskdEndpoint + "/searches/" + id + "/responses"

	resp, err := http.Get(url)

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
