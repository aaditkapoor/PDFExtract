// Extract text from a pdf file.
// and get the github.com url.
// Future: Ability to get recommendation based from list of papers.
// author: aadit

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"code.sajari.com/docconv"
	"github.com/cheggaaa/pb"
	"mvdan.cc/xurls/v2"
)

// Flags
var pdfURL = flag.String("url", "", "PDF URL (required)")
var download = flag.Bool("download", false, "Download the Source.")

// Downloader downloader interface
type Downloader interface {
	Download() error
	Get() string
}

// FileDownload download a file
type FileDownload struct {
	fileName string
	url      string
	response *http.Response
	text     string
}

// NewFileDownload create a new FileDownload object
func NewFileDownload(url string, filename string) FileDownload {
	fd := FileDownload{url: url, fileName: filename}
	return fd
}

// FileParser parse the file for code links.
type FileParser struct {
	fileDownloader *FileDownload
	codeLinks      []string
}

// NewFileParser Create a new FileParser object.
func NewFileParser(fileDownloader *FileDownload) FileParser {
	fp := FileParser{fileDownloader: fileDownloader}
	return fp
}

// PrintLinks print the gathered links.
func (fp *FileParser) PrintLinks() {
	log.Println("------------LINKS------------")
	for i, v := range fp.codeLinks {
		log.Printf("%d - %s", i, v)
	}
}

// Parse parses the file for code links
// Downloads the file
// Extracts texts
func (fp *FileParser) Parse() error {

	err := fp.fileDownloader.Download()
	if err != nil {
		fmt.Println(err)
	}

	// Conversion of pdf to text (pdftotext should exist)
	log.Printf("CONVERTING %s", fp.fileDownloader.fileName)
	res, err := docconv.ConvertPath(fp.fileDownloader.fileName)
	if err != nil {
		log.Fatal(err)
	}
	fp.fileDownloader.text = res.Body
	log.Printf("DONE")

	rx := xurls.Strict()

	for _, v := range rx.FindAllString(fp.fileDownloader.text, -1) {
		fp.codeLinks = append(fp.codeLinks, v)
	}
	log.Printf("Gathering links...")
	log.Printf("DONE")
	fp.PrintLinks()
	return nil
}

/* Implementation: Download */

// Download downloads a file (saves the content of the file in text)
func (fileDownload *FileDownload) Download() error {
	// Download the file and save as fileDownload.filename
	// check for if the file exists and if exists then use the downloaded version.
	// If file does not exist then download file and convert.
	if _, err := os.Stat(fileDownload.fileName); os.IsNotExist(err) {
		u := fileDownload.url
		// Downloading the file
		downloadError := fileDownload.DownloadFile(fileDownload.fileName, u)
		if downloadError != nil {
			panic(downloadError)
		}
		log.Printf("Downloaded in %s.", fileDownload.fileName)
		return nil
	}
	log.Printf("Using %s", fileDownload.fileName)
	return nil
}

// Get return a file
func (fileDownload *FileDownload) Get() string {
	return fileDownload.text
}

// checks for a url
func checkForURL(urlString string) bool {
	// check for url structure here.
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	return true
}

/* end */

// DownloadFile downloads the file.
func (fileDownload *FileDownload) DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	go func() {
		n, err := io.Copy(out, resp.Body)
		if n != fileSize {
			log.Fatal("Truncated")
		}
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

	countSize := int(fileSize / 1000)
	bar := pb.StartNew(countSize)
	var fi os.FileInfo
	for fi == nil || fi.Size() < fileSize {
		fi, _ = out.Stat()
		bar.SetCurrent(int64(fi.Size() / 1000))
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	return nil
}

// check if /pdfs exist.
func checkIfPdfsExist() bool {
	if _, err := os.Stat("pdfs"); os.IsNotExist(err) {
		errDir := os.MkdirAll("pdfs", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		return false
	}
	return true
}

// Check if executable is in bin.
func execExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func main() {
	if !execExists("pdftotext") {
		log.Fatal("pdftotext does not exists.Please install pdftotext from: http://www.xpdfreader.com/download.html\n")
	}
	if !checkIfPdfsExist() {
		log.Print("/pdfs did not exist! Created /pdfs...")

	}

	flag.Parse()
	if *pdfURL == "/" || *pdfURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !checkForURL(*pdfURL) {
		fmt.Fprintf(os.Stderr, "Error: %s not a valid URL.\n", *pdfURL)
	}

	// Create a new fileDownload object and run FileParser on it.

	last := strings.Split(*pdfURL, "/")
	fileName := fmt.Sprintf("pdfs/%s", last[len(last)-1])
	fd := NewFileDownload(*pdfURL, fileName)
	fp := NewFileParser(&fd)

	// Parsing file
	err := fp.Parse()

	if err != nil {
		fmt.Println(err)
	}
}
