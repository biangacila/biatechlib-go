package pdfs

import (
	"fmt"
	bia_dropbox "github.com/biangacila/biatechlib-go/dropbox"
	"io/ioutil"
	"os"
	"os/exec"
)

type GeneratePdfFromHtml struct {
	HtmlBody string
	//Reference for file name
	AppName string
	Org     string
	Module  string
	Ref     string

	DropboxRoot     string
	DropboxClientId string
	DropboxSecret   string
	DropboxToken    string

	ActionSaveToDropbox bool
	ActionSendEmail     bool
	OwnBucket           bool

	fileNameHtml string
	fileNamePdf  string
	Result       struct {
		FileName     string
		LinkRevision string
		LinkDownload string
	}
}

func (s *GeneratePdfFromHtml) Run() {
	if !s.OwnBucket {
		s.DropboxRoot = fmt.Sprintf("companies/files/pdf")
	}
	s.createWorkingFileNames()
	s.createHtmlFileFromStringTemporally()
	s.generatePdfFileFromHtml()
	if s.ActionSaveToDropbox {
		s.uploadToDropbox()
	}
}
func (s *GeneratePdfFromHtml) uploadToDropbox() {
	ref := fmt.Sprintf("%v/%v", s.DropboxRoot, s.fileNamePdf)
	src := s.fileNamePdf
	bucket := s.DropboxRoot
	myFilter := ref
	dst := bucket + "" + myFilter
	fmt.Println("destination file name: ", dst)
	bdb := bia_dropbox.HubDropbox{}
	//todo authentication to dropbox
	bdb.DROPBOX_CLIENTID = s.DropboxClientId
	bdb.DROPBOX_SECRET = s.DropboxSecret
	bdb.DROPBOX_TOKEN = s.DropboxToken

	bdb.Company = s.Org
	bdb.UploadFilename = dst
	bdb.UploadBusket = bucket
	bdb.UploadSrc = src
	entry := bdb.Upload(src, dst)
	s.Result.LinkRevision = entry.Revision
	fmt.Println("destination  entry.Revision: ", entry.Revision)
	s.Result.LinkDownload = fmt.Sprintf("https://cloudcalls.easipath.com/backend-pmis/api/files/download/%v", s.Result.LinkRevision)
}
func (s *GeneratePdfFromHtml) generatePdfFileFromHtml() {
	fileUrlHtml := s.fileNameHtml
	fileUrlPdf := s.fileNamePdf

	err := ioutil.WriteFile(fileUrlHtml, []byte(s.HtmlBody), 0666)
	if err != nil {
		fmt.Print("ERRER ioutil.WriteFile html :> ", err)
	}
	// converting "compiled-html.html" file to "my.pdf"
	//strCommand:=fmt.Sprintf("%v %v %v %v %v","xvfb-run","wkhtmltopdf","-O landscape",fileUrlHtml,fileUrlPdf)
	err = exec.Command("xvfb-run", "wkhtmltopdf", fileUrlHtml, fileUrlPdf).Run()
	if err != nil {
		fmt.Printf("\n\nError on trying to generate PDF: %s", err)
	}

}
func (s *GeneratePdfFromHtml) createHtmlFileFromStringTemporally() {
	f, err := os.OpenFile(s.fileNameHtml, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(s.HtmlBody); err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
}
func (s *GeneratePdfFromHtml) createWorkingFileNames() {
	str := fmt.Sprintf("%v--%v--%v--%v",
		s.AppName, s.Org, s.Module, s.Ref)
	s.fileNameHtml = str + ".html"
	s.fileNamePdf = str + ".pdf"
}
