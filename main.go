package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"path"
	"strings"

	"github.com/divan/qrlogo"
)

var (
	input    = flag.String("i", "data.csv", "CSV file in format ID;data")
	logoPath = flag.String("l", "logo.png", "Logo in PNG format")
	output   = flag.String("o", "out/", "Output folder")
	size     = flag.Int("size", 1024, "Image size in pixels")
)

func main() {
	flag.Usage = Usage
	flag.Parse()

	data := readCsvFile(*input)
	err := os.MkdirAll(*output, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(*logoPath)
	errcheck(err, "Failed to open logo:")
	defer file.Close()

	logo, _, err := image.Decode(file)
	errcheck(err, "Failed to decode PNG with logo:")

	for _, text := range data {
		fmt.Println(text)
		text = strings.Split(text[0], ";")
		dataString := padString(text[1])
		fmt.Println(dataString)
		createQRCode(dataString, logo, *size, path.Join(*output, text[0]+".png"))
	}

	//createQRCode(text)
}

func padString(input string) string {
	toAdd := 20 - len(input)
	if toAdd <= 0 {
		return input
	}
	input = input + ";"
	for i := 0; i < toAdd; i++ {
		input = input + "0"
	}
	return input
}

func createQRCode(text string, logo image.Image, qrSize int, outputFile string) {

	qr, err := qrlogo.Encode(text, logo, qrSize)
	errcheck(err, "Failed to encode QR:")

	out, err := os.Create(outputFile)
	errcheck(err, "Failed to open output file:")
	out.Write(qr.Bytes())
	out.Close()

	fmt.Println("Done! Written QR image to", outputFile)
}
func Usage() {
	fmt.Fprintln(os.Stderr, "Usage: qrlogo [options] text")
	flag.PrintDefaults()
}

func errcheck(err error, str string) {
	if err != nil {
		fmt.Println(str, err)
		os.Exit(1)
	}
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	return records
}
