package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

func main() {

	hygraphEndpointInput := ""
	inputFilePath := ""
	flag.StringVar(&hygraphEndpointInput, "hygraphEndpoint", "", "Public content URL for Hygraph")
	flag.StringVar(&inputFilePath, "filePath", "", "Path to the file to use as an input")

	flag.Parse()

	hygraphEndpoint, inputPath := validateAndPromptForInputs(hygraphEndpointInput, inputFilePath)
	inputData := getInputData(inputPath)
	existingGameData := getExistingGameData(hygraphEndpoint)

	ok := validateProcessingData(&inputData, existingGameData)
	if !ok {
		return
	}

	for idx, inputGame := range inputData {
		ok := validateInputGame(idx, inputGame)
		if !ok {
			continue
		}
		existingGame := findExistingGameData(existingGameData, inputGame.GameTitle)
		id := ""
		if existingGame != nil {
			id = existingGame.ID
		}
		pushGameData(id, inputGame, hygraphEndpoint)
	}
}

func validateAndPromptForInputs(endpoint string, filePath string) (string, string) {
	url := endpoint
	inputPath := filePath

	if url == "" {
		fmt.Println("You did not provide an URL to connect to. What is the Public Content URL for your Hygraph Project?")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to parse user input")
			log.Fatal(err.Error())
			return "", ""
		}
		url = input
	}

	if inputPath == "" {
		fmt.Println("You did not provide a path for your input csv. Please enter it now: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to parse user input")
			log.Fatal(err.Error())
			return "", ""
		}
		inputPath = input
	}

	return strings.ReplaceAll(url, "\n", ""), strings.ReplaceAll(inputPath, "\n", "")
}

func validateProcessingData(inputData *[]*InputData, existingGameData []*BoardGame) bool {
	if inputData == nil || existingGameData == nil {
		return false
	}

	if len(existingGameData) == 0 {
		log.Println("No existing games found in the CMS. Would you like to continue? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to parse user input")
			log.Fatal(err.Error())
			return false
		}
		if strings.ToLower(input) == "n" {
			log.Println("Exiting")
			return false
		}
	} else {
		log.Printf("Found %v existing games in the CMS.\n", len(existingGameData))
	}

	return true
}

func validateInputGame(line int, inputGame *InputData) bool {
	if inputGame.GameTitle == "" {
		log.Printf("\nskipping line %v because it is missing a game title\n", line+2)
		return false
	}

	// Remove " characters from Description and Notes fields to resolve update/creation issues
	inputGame.Description = strings.ReplaceAll(inputGame.Description, "\"", "")
	inputGame.Notes = strings.ReplaceAll(inputGame.Notes, "\"", "")

	return true
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

	log.Printf("Found %v games in the input file for updates\n", len(games))

	return games
}

func findExistingGameData(games []*BoardGame, title string) *BoardGame {
	for _, game := range games {
		if game.GameTitle == title {
			return game
		}
	}
	return nil
}
