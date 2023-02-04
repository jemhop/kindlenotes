package main

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

//file to detect mounted kindle and find the relevant file

// placeholder func for now
func getMountedKindle() string {
	return "/home/emma/Projects/kindlenotes/My Clippings.txt"
}

func openClippingsFile(path string) string {
	file, _ := readFileUTF16(path)
	return string(file)

}

func readFileUTF16(filename string) ([]byte, error) {

	// Read the file into a []byte:
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, err := ioutil.ReadAll(unicodeReader)
	return decoded, err
}
