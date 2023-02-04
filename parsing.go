package main

import (
	"strings"
	"time"
)

const ENTRY_SEPERATOR = "\r\n=========="

const (
	CLIPPING = iota
	NOTE
	BOOKMARK
)

// before being sorted into a note or a section, data is read into a single struct which contains all information an entry in the file can possibly contain
type section struct {
	content     string
	sectionType int

	//struct nested within struct for storing notes
	notes []section

	//both needed to correctly match notes to highlights
	bookLineStart int
	bookLineEnd   int

	date time.Time

	bookAuthor string
	bookTitle  string
}

// as of current only supports one author
type book struct {
	authors    string
	name       string
	highlights []section
}

// Return clippings as an array of strings seperated by ENTRY_SEPERATOR
func splitClippings(clippings string) (sortedClippings []string) {
	clippings = trimEmptyLines([]byte(clippings))
	sortedClippings = strings.Split(clippings, ENTRY_SEPERATOR)
	return sortedClippings
}

/* func readClippingsFile(content []string) []book {
	sections := getSections(content)

	return clippingsToBookArray(sections)
} */

func getSections(content string) []section {
	clippings := splitClippings(content)
	sectionArray := []section{}
	for _, c := range clippings {
		section := readSection(c)
		if section.sectionType == BOOKMARK {
			continue
		}

		sectionArray = append(sectionArray, readSection(c))
	}

	return sectionArray
}

// returns the read section data and a location from which to read the next section
// this is a little horrendous
func readSection(content string) (output section) {
	returnSection := section{}
	content = trimEmptyLines([]byte(content))
	contentLines := strings.Split(content, "\n")

	//not yet supporting bookmarks, see no need to
	output.sectionType = parseClippingType(contentLines[1])
	if output.sectionType == BOOKMARK {
		return
	}

	book, author := parseTitleLine(contentLines[0])
	println(book + " " + author)

	nums := ParseNum(contentLines[1])
	println(contentLines[1])
	for _, n := range nums {
		println(n)
	}

	/* for i := 0; i < len(contentLines); i++ {
		line := contentLines[i]

		if i == 0 {
			bookTitle, authorName := parseTitleLine(line)
			returnSection.bookAuthor = authorName
			returnSection.bookTitle = bookTitle
			continue
		} else if i == 1 {
			page, bookLocationStart, bookLocationEnd, sectionType, date := parseInfoLine(line)
			returnSection.page = page
			returnSection.bookLineEnd = bookLocationStart
			returnSection.bookLineEnd = bookLocationEnd
			returnSection.sectionType = sectionType
			returnSection.date = date
			continue
		}

		returnSection.content += line

	} */

	return returnSection
}

/* // - Your Highlight on page 94 | location 1428-1429 | Added on Thursday, 26 August 2021 20:34:17
func parseInfoLine(line string) (page int, bookLocationStart int, bookLocationEnd int, sectionType int, date time.Time) {
	nums := ParseNum(line)
	println(line)
	println(nums)

} */

func parseClippingType(infoLine string) (sectionType int) {
	if strings.Contains(infoLine, "Highlight") {
		sectionType = CLIPPING
	} else if strings.Contains(infoLine, "Note") {
		sectionType = NOTE
	} else if strings.Contains(infoLine, "Bookmark") {
		sectionType = BOOKMARK
	}
	return sectionType
}

func parseTitleLine(line string) (book string, author string) {
	var bracketIndex int = -1
	for i, c := range line {
		if c == '(' {
			bracketIndex = i
		}
	}

	if bracketIndex == -1 {
		return "", ""
	}

	book = strings.TrimSpace(line[0 : bracketIndex-1])
	author = strings.TrimSpace(line[bracketIndex+1 : len(line)-2])

	return book, author
}

/* func readClippingsAsLineArray() []string {
	location := getMountedKindle()

	return openFileAsStringArray(location)
}

func clippingsToBookArray(section []section) []book {
	return []book{}
}
*/
