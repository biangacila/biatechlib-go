package pdfs

import "testing"

func TestServiceMergePdf_Run(t *testing.T) {
	fileList := []string{
		"./docs/ID-Certified.pdf",
		"./docs/Proof-Qualification.pdf",
		"./docs/proof-tax-number.pdf",
		"./docs/csv.pdf",
		"./docs/bank-account-confirmation.pdf",
	}
	outputFile := "SUPPORTING-DOCUMENTATION-REQUIRED.pdf"
	hub := ServiceMergePdf{
		FileList:   fileList,
		OutputFile: outputFile,
	}
	hub.Run()
}
