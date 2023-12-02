package ocr

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"os"
)

//  sudo apt-get install libleptonica-dev
/*
export CGO_CFLAGS="-I/usr/include/leptonica"
export CGO_LDFLAGS="-llept"

sudo apt-get install tesseract-ocr

*/
type ServiceConvertImage struct {
	ImageLink string // Mozambique_plate_02.jpg
}

func (obj *ServiceConvertImage) Run() {
	os.Setenv("CGO_CFLAGS", "-I/usr/include/leptonica")
	os.Setenv("CGO_LDFLAGS", "-llept")
	// Initialize Tesseract OCR client
	client := gosseract.NewClient()
	defer client.Close()

	// Set the path to the Tesseract executable (if not in PATH)
	// client.SetTessExecutable("/path/to/tesseract")

	// Set the language data directory (if not in the default location)
	// client.SetLanguage("eng") // Replace "eng" with the language you need

	// Load an image from file
	err := client.SetImage(obj.ImageLink) // Replace "input.png" with your image file path
	if err != nil {
		fmt.Println("Error setting image:", err)
		return
	}

	// Perform OCR
	text, err := client.Text()
	if err != nil {
		fmt.Println("Error performing OCR:", err)
		return
	}

	// Print the extracted text
	fmt.Println("Extracted Text:")
	fmt.Println(text)

}
