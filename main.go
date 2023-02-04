package main

import (
	"bufio"
	"kindlenotes/kindle"
	"os"
)

func main() {

	file := kindle.GetClippingsFileContent()

	sections := kindle.ReadClippingsFileAsSectionArray(file)
	scanner := bufio.NewScanner(os.Stdin)

	for _, s := range sections {
		PrintClippingAsJson(s)
		scanner.Scan()
	}

}
