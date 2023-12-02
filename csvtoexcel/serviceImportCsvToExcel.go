package csvtoexcel

import (
	"encoding/csv"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
)

type ServiceImportCsvToExcel struct {
	TemplateFileLink  string
	TemplateFileLink2 string
	CsvFileLink       string //
}

func (s *ServiceImportCsvToExcel) Run() {
	s.removeAllSheet()
	// Open the Excel template file
	excelFile, err := xlsx.OpenFile(s.TemplateFileLink2)
	if err != nil {
		panic(err)
	}

	// Open the CSV file for reading
	csvFile, err := os.Open(s.CsvFileLink)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	// Create a new sheet or get an existing one
	sheet := excelFile.Sheets[16]
	// Remove all rows from the sheet.
	for x := 0; x < len(sheet.Rows); x++ {
		sheet.RemoveRowAtIndex(x)
	}

	//global.DisplayObject("SHeets", sheet.Rows)
	// Read CSV data and populate Excel cells
	csvReader := csv.NewReader(csvFile)
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		row := sheet.AddRow()
		for _, value := range record {
			cell := row.AddCell()
			cell.Value = value
		}
	}

	// Save the modified Excel file
	err = excelFile.Save("output.xlsx")
	if err != nil {
		panic(err)
	}

}
func (s *ServiceImportCsvToExcel) removeAllSheet() {
	// Replace with the path to your Excel file.
	filePath := s.TemplateFileLink

	// Open the Excel file.
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open Excel file: %v", err)
	}

	// Specify the sheet you want to empty.
	sheetIndex := 16 // Change this to the index of the sheet you want to empty.

	if sheetIndex < len(xlFile.Sheets) {
		// Get the sheet.
		sheet := xlFile.Sheets[sheetIndex]

		// Remove all rows from the sheet.
		for len(sheet.Rows) > 0 {
			sheet.RemoveRowAtIndex(0)
		}

		// Save the changes to the Excel file.
		err = xlFile.Save(s.TemplateFileLink2)
		if err != nil {
			log.Fatalf("Failed to save changes to Excel file: %v", err)
		}

		fmt.Printf("Sheet emptied successfully.\n")
	} else {
		fmt.Printf("Invalid sheet index.\n")
	}
}
