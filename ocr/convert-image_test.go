package ocr

import "testing"

func TestServiceConvertImage_Run(t *testing.T) {
	var service ServiceConvertImage
	service.ImageLink = "Mozambique_plate_02.jpg"
	service.Run()
}
