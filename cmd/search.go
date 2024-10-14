/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"cobraCli/models"

	"cobraCli/exportExcel"

	"github.com/spf13/cobra"
)

var exportExcelFlag bool

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("search called")
		const greenBackgroundBlackText = "\033[1;32m"
		const reset = "\033[0m"
		word := args[0]
		wordData, err := fetchDictionaryData(word)
		if err != nil {
			fmt.Println(err)
		}

		for _, word := range wordData {
			fmt.Printf("\n%sword%s: %s, POS: %s, \nMeanings:\n\t%s,\nExample: %s\nSynonyms: %s \nAntonyms: %s\n",
				greenBackgroundBlackText, reset,
				word.Word, word.Pos, word.Meaning, word.Example, word.Synonyms, word.Antonyms)
		}

		if exportExcelFlag {
			fmt.Println("export in progress...")
			// exportExcel.export(wordData)
			exportExcel.Export(wordData)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVarP(&exportExcelFlag, "export", "e", false, "Set this flag to export to excel")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func fetchDictionaryData(word string) ([]models.ParsedResponse, error) {
	// Make the API call
	apiURL := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into the struct
	var dictionaryResponse []models.DictionaryResponse
	err = json.Unmarshal(body, &dictionaryResponse)
	if err != nil {
		return nil, err
	}

	// var d []models.ParsedResponse
	d := make([]models.ParsedResponse, 0)

	for _, entry := range dictionaryResponse {
		meanings := entry.Meanings

		for _, meaning := range meanings {
			var p models.ParsedResponse

			// firstDefinition := meaning.Definitions[0]

			p.Word = entry.Word
			p.Pos = meaning.PartOfSpeech

			var mList []string
			var eList []string
			var sList []string
			var aList []string

			for _, definition := range meaning.Definitions {
				mList = append(mList, definition.Definition)

				if definition.Example != "" {
					eList = append(eList, definition.Example)
				}

				if len(definition.Synonyms) != 0 {
					sList = append(sList, strings.Join(definition.Synonyms, ","))
				}

				if len(definition.Antonyms) != 0 {
					aList = append(sList, strings.Join(definition.Antonyms, ","))
				}

			}

			p.Meaning = strings.Join(mList, "; \n\t")

			p.Example = strings.Join(eList, ", ")

			if len(sList) != 0 {
				p.Synonyms = strings.Join(sList, ", ")
			} else {
				p.Synonyms = strings.Join(meaning.Synonyms, ", ")
			}

			if len(aList) != 0 {
				p.Antonyms = strings.Join(aList, ", ")
			} else {
				p.Antonyms = strings.Join(meaning.Antonyms, ", ")
			}

			// p.Meaning = firstDefinition.Definition
			// p.Example = firstDefinition.Example
			// p.Synonyms = GetFirstOrNil(firstDefinition.Synonyms)
			// p.Antonyms = GetFirstOrNil(firstDefinition.Antonyms)

			d = append(d, p)
		}
	}

	// Return the response data and any errors
	return d, nil
}

func GetFirstOrNil[T any](slice []T) interface{} {
	if len(slice) > 0 {
		return slice[0]
	}
	return nil
}
