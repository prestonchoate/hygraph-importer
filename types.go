package main

type HygraphResponse struct {
	Data struct {
		Games []*BoardGame `json:"boardGameDatabases"`
	} `json:"data"`
}

type HygraphCreateResponse struct {
	Data struct {
		CreatedGameData BoardGame `json:"createBoardGameDatabase"`
	} `json:"data"`
}

type HygraphUpdateResponse struct {
	Data struct {
		UpdatedGameData BoardGame `json:"updateBoardGameDatabase"`
	} `json:"data"`
}

type BoardGame struct {
	ID            string   `json:"id"`
	Display       bool     `json:"display"`
	GameTitle     string   `json:"gameTitle"`
	GameType1     []string `json:"gameType1"`
	GameType2     []string `json:"gameType2"`
	GameType3     []string `json:"gameType3"`
	MinPlayers    string   `json:"numberOfPlayers"`
	MaxPlayers    string   `json:"numberOfPlayersMax"`
	MinPlayTime   string   `json:"playingTimeMin"`
	MaxPlayTime   string   `json:"playingTimeMax"`
	Age           string   `json:"age"`
	Complexity    string   `json:"complexityRatingOutOf5"`
	AverageRating string   `json:"averageBbgRatingOutOf10"`
	Location      string   `json:"location"`
	Description   string   `json:"description"`
	BBGLink       string   `json:"linkToBbg"`
	Notes         string   `json:"notes"`
}

type InputData struct {
	Display       bool    `csv:"Display"`
	GameTitle     string  `csv:"Game Title"`
	GameType1     string  `csv:"Game Type 1"`
	GameType2     string  `csv:"Game Type 2"`
	GameType3     string  `csv:"Game Type 3"`
	MinPlayers    string  `csv:"Number of Player (Min)"`
	MaxPlayers    string  `csv:"Number of Player (Max)"`
	MinPlayTime   string  `csv:"Playing Time (Min)"`
	MaxPlayTime   string  `csv:"Playing Time (Max)"`
	Age           string  `csv:"Age"`
	Complexity    float32 `csv:"Complexity Rating (Out of 5)"`
	AverageRating float32 `csv:"Average BGG Rating (Out of 10)"`
	Location      string  `csv:"Location"`
	Description   string  `csv:"Description"`
	BBGLink       string  `csv:"Link to BBG"`
	Notes         string  `csv:"Notes"`
}
