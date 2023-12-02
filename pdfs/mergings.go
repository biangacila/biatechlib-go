package pdfs

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type ServiceMergePdf struct {
	FileList   []string
	OutputFile string
}

func (s *ServiceMergePdf) Run() {

	inFiles := s.FileList //[]string{"in1.pdf", "in2.pdf"}
	err := api.MergeCreateFile(inFiles, s.OutputFile, nil)
	if err != nil {
		fmt.Printf("Error creating PDF context: %v\n", err)
		return
	}

	fmt.Printf("Merged PDFs saved as %s\n", s.OutputFile)
}
