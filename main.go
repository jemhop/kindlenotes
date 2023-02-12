package main

import (
	"kindlenotes/kindle"
)

func main() {

	file := kindle.GetClippingsFileContent()

	sections := kindle.ReadClippingsFileAsSectionArray(file)

	books := kindle.SortSections(sections)

	for _, b := range books {
		PrintBookAsJson(b)
	}

}
