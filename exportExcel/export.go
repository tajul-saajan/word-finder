package exportExcel

import (
	"cobraCli/models"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func Export(words []models.ParsedResponse) {
	filename := "word.xlsx"

	// Check if the file exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, append the new data
		fmt.Println("File exists, appending data.")
		appendToExcel(filename, words)
	} else {
		// File does not exist, create a new one
		fmt.Println("File does not exist, creating new file and adding data.")
		createNewExcel(filename, words)
	}
}

func appendToExcel(filename string, words []models.ParsedResponse) {
	// Open the existing file
	file, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Get the existing sheet or create a new one
	sheetName := "Sheet1"
	if index, _ := file.GetSheetIndex(sheetName); index == 0 {
		file.NewSheet(sheetName)
	}

	// Find the next empty row
	rows, err := file.GetRows(sheetName)
	if err != nil {
		fmt.Println("Error reading rows:", err)
		return
	}
	startRow := len(rows) + 1 // Start after the last row

	// Write data to the sheet
	// for i, word := range words {
	// 	file.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow+i), word.Word)
	// 	file.SetCellValue(sheetName, fmt.Sprintf("B%d", startRow+i), word.Pos)
	// 	file.SetCellValue(sheetName, fmt.Sprintf("C%d", startRow+i), word.Meaning)
	// 	file.SetCellValue(sheetName, fmt.Sprintf("D%d", startRow+i), word.Example)
	// 	file.SetCellValue(sheetName, fmt.Sprintf("E%d", startRow+i), word.Synonyms)
	// 	file.SetCellValue(sheetName, fmt.Sprintf("F%d", startRow+i), word.Antonyms)
	// }

	appendWordsToSheet(file, words, startRow, sheetName)

	// Save the file
	if err := file.Save(); err != nil {
		fmt.Println("Error saving file:", err)
	}
}

// Function to create a new Excel file
func createNewExcel(filename string, words []models.ParsedResponse) {
	// Create a new Excel file
	f := excelize.NewFile()

	// Add a new sheet
	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	// Write headers
	f.SetCellValue(sheetName, "A1", "Word")
	f.SetCellValue(sheetName, "B1", "POS")
	f.SetCellValue(sheetName, "C1", "Meaning")
	f.SetCellValue(sheetName, "D1", "Example")
	f.SetCellValue(sheetName, "E1", "Synonyms")
	f.SetCellValue(sheetName, "F1", "Meaning")

	// Write the data starting from the second row
	// for i, word := range words {
	// 	f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), word.Word)
	// 	f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), word.Meaning)
	// }

	appendWordsToSheet(f, words, 2, sheetName)

	// Save the file with the given name
	if err := f.SaveAs(filename); err != nil {
		fmt.Println("Error saving file:", err)
	}
}

func appendWordsToSheet(file *excelize.File, words []models.ParsedResponse, startRow int, sheetName string) {
	for i, word := range words {
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow+i), word.Word)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", startRow+i), word.Pos)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", startRow+i), word.Meaning)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", startRow+i), word.Example)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", startRow+i), word.Synonyms)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", startRow+i), word.Antonyms)
	}
}
