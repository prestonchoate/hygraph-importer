package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/machinebox/graphql"
)

type HygraphResponse struct {
  Data struct {
    Games []*BoardGames `json:"boardGameDatabases"`
  } `json:"data"`
}

type BoardGames struct {
  ID string `json:"id"`
  Display bool `json:"display"`
  GameTitle string `json:"gameTitle"`
  GameType1 []string `json:"gameType1"`
  GameType2 []string `json:"gameType2"`
  GameType3 []string `json:"gameType3"`
  Players string `json:"numberOfPlayers"`
  MinPlayTime string `json:"PlayingTimeMin"`
  Age string `json:"age"`
  Complexity string `json:"complexityRatingOutOf5"`
  AverageRating string `json:"averageBbgRatingOutOf10"`
  Location string `json:"location"`
  Description string `json:"description"`
  BBGLink string `json:"linkToBbg"`
  Notes string `json:"notes"`
}

type InputData struct {
  Display bool `csv:"Display"`
  GameTitle string `csv:"Game Title"`
  GameType1 string `csv:"Game Type 1"`
  GameType2 string `csv:"Game Type 2"`
  GameType3 string `csv:"Game Type 3"`
  Players string `csv:"Number of Players"`
  MinPlayTime string `csv:"Playing Time (Min)"`
  Age string `csv:"Age"`
  Complexity float32 `csv:"Complexity Rating (Out of 5)"`
  AverageRating float32 `csv:"Average BGG Rating (Out of 10)"`
  Location string `csv:"Location"`
  Description string `csv:"Description"`
  BBGLink string `csv:"Link to BBG"`
  Notes string `csv:"Notes"`
}

func main() {
  hygraphEndppoint := "https://api-us-east-1-shared-usea1-02.hygraph.com/v2/cltttqazc0cc907uwzeapa71t/master"
  inputData := getInputData("input.csv")
  existingGameData := getExistingGameData(hygraphEndppoint)
  
  if inputData == nil || existingGameData == nil {
    return
  }

  for _, inputGame := range inputData {
    existingGame := findExistingGameData(existingGameData.Data.Games, inputGame.GameTitle)
    id := ""
    if existingGame != nil {
      id = existingGame.ID
    }
    pushGameData(id, inputGame, hygraphEndppoint)
  }
}

//TODO: There is a mutation `updateManyBoardGameDatabasesConnection` that can accept multiple payloads at once
//      we should switch to that at some point
func pushGameData(id string, gameData *InputData, hygraphEndpoint string) {
  if id != "" {
    fmt.Printf("Trying to update game ID: %v\n", id)
    req := fmt.Sprintf(`
        mutation {
        updateBoardGameDatabase(where: {
          id: %v
        },
        data: { 
          display: %v,
          gameTitle: "%v",
          gameType1: ["%v"],
          gameType2: ["%v"],
          gameType3: ["%v"],
          numberOfPlayers: "%v",
          playingTimeMin: "%v",
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
      gameData.Players,
      gameData.MinPlayTime,
      gameData.Age,
      gameData.Complexity,
      gameData.AverageRating,
      gameData.Location,
      gameData.Description,
      gameData.BBGLink,
      gameData.Notes,
      )

    data := sendGqlRequest(hygraphEndpoint, req)
    fmt.Println(string(data))
  } else {
    fmt.Printf("Trying to create new game with Title: %v\n", gameData.GameTitle)
    req := fmt.Sprintf(`
      mutation {
        createBoardGameDatabase(data: { 
          display: %v,
          gameTitle: "%v",
          gameType1: ["%v"],
          gameType2: ["%v"],
          gameType3: ["%v"],
          numberOfPlayers: "%v",
          playingTimeMin: "%v",
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
      gameData.Players,
      gameData.MinPlayTime,
      gameData.Age,
      gameData.Complexity,
      gameData.AverageRating,
      gameData.Location,
      gameData.Description,
      gameData.BBGLink,
      gameData.Notes,
      )

    data := sendGqlRequest(hygraphEndpoint, req)
    fmt.Println(string(data))
  }
}


func getInputData(inputPath string) []*InputData {
 input, err := os.Open(inputPath)
  if err != nil {
    log.Fatal("failed to open input file")
    return nil
  }
  defer input.Close()

  games := []*InputData{}

  err = gocsv.UnmarshalFile(input, &games)
  if err != nil {
    log.Println(err.Error())
    log.Fatal("failed to unmarshal csv data into struct")
  }

  return games
}

func getExistingGameData(hygraphEndpoint string) *HygraphResponse {
  req := `query {
  boardGameDatabases {
    id
    stage
    gameTitle
    display
    gameType1
    gameType2
    gameType3
    numberOfPlayers
    playingTimeMin
    age
    complexityRatingOutOf5
    averageBbgRatingOutOf10
    location
    description
    linkToBbg
    notes
  }
}`

  data := sendGqlRequest(hygraphEndpoint, req)

  response := &HygraphResponse{}
  err := json.Unmarshal(data, &response)
  if err != nil {
    log.Fatal(err.Error())
    return nil
  }

  for _, game := range response.Data.Games {
    fmt.Printf("%+v\n", game)
  }

  return response
}

func findExistingGameData(games []*BoardGames, title string) *BoardGames {
  for _, game := range games {
    if game.GameTitle == title {
      return game
    }
  }
  return nil
}

func __sendGqlRequest(endpoint string, reqBody string) any {
  client := graphql.NewClient(endpoint)
  req := graphql.NewRequest(reqBody)
  ctx := context.Background()
  var res any
  err := client.Run(ctx, req, &res)
  if err != nil {
    log.Fatal(err.Error())
    return nil
  }
  return res
}

func sendGqlRequest(hygraphEndpoint string, req string) []byte {
  
  reqBody := map[string]string{"query": req}
  jsonBody, err := json.Marshal(reqBody)
  if err != nil {
    log.Println("failed to marshal request to json")
    log.Fatal(err.Error())
  }

  log.Printf("Attempting to send the following request:\n\n\n\n\"%+v\"\n\n\n", string(jsonBody))
  jsonBodyIo := strings.NewReader(string(jsonBody))

  request, _ := http.NewRequest("POST", hygraphEndpoint, jsonBodyIo)
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
