# Hygraph Board Game Data Importer
The purpose of this tool is to provide an easy way to accept a CSV file with headers containing board game data and import it into a Hygraph CMS instance. 

## Installation
1. Download the most recent version from the Releases tab for your specific operating system and CPU architecture.
2. Extract the archive into a directory.
3. Open a terminal and navigate to the directory where this tool was extracted.
4. Run `./hygraph-importer` to get started

## Usage
Running `hygraph-importer` will start the CLI for this tool. If you supply no flags to the CLI you will be prompted for them instead.

### Flags
* \-filePath This is the path to the file you wish to use as an input for data. The data should be provided in a CSV format with properly formatted headers. Failure to match the input headers schema will produce erratic results.
* \-hygraphEndpoint This is the public content URL for your Hygraph instance. 
* \-\-help Will display a simple help page for this CLI

### Input File Schema
This tool expects the input file CSV to have the following headers named exactly as listed
* Display
* Game Title
* Game Type 1
* Game Type 2
* Game Type 3
* Number of Player (Min)
* Number of Player (Max)
* Playing Time (Min)
* Playing Time (Max)
* Age
* Complexity Rating (Out of 5)
* Average BGG Rating (Out of 10)
* Location
* Description
* Link to BBG
* Notes

## Disclaimers
While I have taken care to handle as many edge cases and foot guns as possible there are possibly scenarios which I could not predict. If you encounter any such errors please open an issue in this github repository with as much detail as possible. There are plenty of validation scenarios that I did not account for that will break execution; failure to provide an endpoint or input path, non CSV file provided as input, and unreachable or incorrect endpoint url for example.
