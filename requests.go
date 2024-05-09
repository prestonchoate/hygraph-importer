package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func getExistingGameData(hygraphEndpoint string) []*BoardGame {
	games := []*BoardGame{}
	lastId := ""

	response := getNextPage(lastId, hygraphEndpoint)

	for len(response.Data.Games) > 0 {
		games = append(games, response.Data.Games...)
		lastId = games[len(games)-1].ID
		response = getNextPage(lastId, hygraphEndpoint)
	}

	return games
}

func getNextPage(lastId string, hygraphEndpoint string) *HygraphResponse {
	req := fmt.Sprintf(`query {
		boardGameDatabases (first: 100 %v) {
    id
    stage
    gameTitle
    display
    gameType1
    gameType2
    gameType3
    numberOfPlayers
    numberOfPlayersMax
    playingTimeMin
    playingTimeMax
    age
    complexityRatingOutOf5
    averageBbgRatingOutOf10
    location
    description
    linkToBbg
    notes
  }
}`, getPaginationIdField(lastId))

	data := sendGqlRequest(hygraphEndpoint, req)

	response := &HygraphResponse{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return response
}

func getPaginationIdField(lastId string) string {
	if lastId != "" {
		return fmt.Sprintf(", after: \"%v\"", lastId)
	} else {
		return ""
	}
}

// TODO: There is a mutation `updateManyBoardGameDatabasesConnection` that can accept multiple payloads at once
// we should switch to that at some point
func pushGameData(id string, gameData *InputData, hygraphEndpoint string) {
	if id != "" {
		log.Printf("\nTrying to update game ID: %v\n", id)
		processUpdateRequest(id, gameData, hygraphEndpoint)
	} else {
		log.Printf("\nTrying to create new game with Title: %v\n", gameData.GameTitle)
		processCreateRequest(gameData, hygraphEndpoint)
	}
}

func processCreateRequest(gameData *InputData, hygraphEndpoint string) {
	req := fmt.Sprintf(`
      mutation {
        createBoardGameDatabase(data: { 
          display: %v,
          gameTitle: "%v",
          gameType1: ["%v"],
          gameType2: ["%v"],
          gameType3: ["%v"],
          numberOfPlayers: "%v",
          numberOfPlayersMax: "%v",
          playingTimeMin: "%v",
          playingTimeMax: "%v",
          age: "%v",
          complexityRatingOutOf5: "%v",
          averageBbgRatingOutOf10: "%v",
          location: "%v",
          description: "%v",
          linkToBbg: "%v",
          notes: "%v",
        }) {
          id
        }
      }
      `,
			gameData.Display,
			gameData.GameTitle,
			gameData.GameType1,
			gameData.GameType2,
			gameData.GameType3,
			gameData.MinPlayers,
			gameData.MaxPlayers,
			gameData.MinPlayTime,
			gameData.MaxPlayTime,
			gameData.Age,
			gameData.Complexity,
			gameData.AverageRating,
			gameData.Location,
			gameData.Description,
			gameData.BBGLink,
			gameData.Notes,
	)

	data := sendGqlRequest(hygraphEndpoint, req)
	response := HygraphCreateResponse{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Println("Failed to parse response on creation. Unsure if new game was created")
		log.Println(err.Error())
	}
	if response.Data.CreatedGameData.ID != "" {
		log.Printf("Game successfully created. New ID is: %v\n", response.Data.CreatedGameData.ID)
	}

}

func processUpdateRequest(id string, gameData *InputData, hygraphEndpoint string) {
	req := fmt.Sprintf(`
        mutation {
        updateBoardGameDatabase(where: {
          id: "%v"
        },
        data: { 
          display: %v,
          gameTitle: "%v",
          gameType1: ["%v"],
          gameType2: ["%v"],
          gameType3: ["%v"],
          numberOfPlayers: "%v",
          numberOfPlayersMax: "%v",
          playingTimeMin: "%v",
          playingTimeMax: "%v",
          age: "%v",
          complexityRatingOutOf5: "%v",
          averageBbgRatingOutOf10: "%v",
          location: "%v",
          description: "%v",
          linkToBbg: "%v",
          notes: "%v",
        }) 
        {
          id
        }
      }
      `,
		id,
		gameData.Display,
		gameData.GameTitle,
		gameData.GameType1,
		gameData.GameType2,
		gameData.GameType3,
		gameData.MinPlayers,
		gameData.MaxPlayers,
		gameData.MinPlayTime,
		gameData.MaxPlayTime,
		gameData.Age,
		gameData.Complexity,
		gameData.AverageRating,
		gameData.Location,
		gameData.Description,
		gameData.BBGLink,
		gameData.Notes,
	)

	data := sendGqlRequest(hygraphEndpoint, req)
	response := HygraphUpdateResponse{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Println("Failed to parse response on update request. Unsure if game was updated correctly")
		log.Println(err.Error())
	}
	if response.Data.UpdatedGameData.ID != "" {
		log.Println("Game successfully updated")
	}

}

func sendGqlRequest(hygraphEndpoint string, req string) []byte {

	reqBody := map[string]string{"query": req}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("failed to marshal request to json")
		log.Fatal(err.Error())
	}
	jsonBodyIo := strings.NewReader(string(jsonBody))

	request, err := http.NewRequest("POST", hygraphEndpoint, jsonBodyIo)
	if err != nil {
		log.Println("failed to connect to supplied URL")
		log.Fatal(err.Error())
		return nil
	}
	request.Header.Set("content-type", "application/json")
	client := &http.Client{Timeout: time.Second * 60}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("failed to execute request to hygraph")
		log.Fatal(err.Error())
		return nil
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to read response body")
		log.Fatal(err.Error())
		return nil
	}
	return data
}
