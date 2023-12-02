package csvtoexcel

import "testing"

func TestServiceImportCsvToExcel_Run(t *testing.T) {
	var service ServiceImportCsvToExcel
	service.TemplateFileLink = "LabourReportTemplate.xlsm"
	service.TemplateFileLink2 = "LabourReportTemplate2.xlsm"
	service.CsvFileLink = "labour_logs.csv"
	service.Run()
}
