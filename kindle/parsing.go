package kindle

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/timefmt-go"
)

/*A combination of a rant and important information on how My Clippings.txt is layed out. Partially necessary to understand the code, and partially because
I need someone to share in my pain.

Every single format the Kindle uses for ebooks has the location of any highlights and notes saved differently.
-MOBI:
	-Location of highlights is saved with two locations, location of notes has one
	-Both have a page number
-AZW3:
	-Highlights have two line locations, notes have one
	-Neither has a page number
-PDF:
	-Highlights have TWO PAGE NUMBERS, NOTES HAVE ONE PAGE NUMBER
	-NEITHER HAS A LOCATION

Bookmarks of any kind only have a page location

All of this combined makes this file a nightmare to parse, even when you're hardcoding everything, which is why this file is so long and complicated
It also means this program will never properly support PDFs, and at best will let you see highlights and notes without connecting them. This is also going to make
writing proper display functionality a pain in the ass.

The worst bit is that the Kindle has an internal database of these things. You can delete a clipping from
this file and it'll stay displayed on the Kindle. That means this fucking stupid file isn't even a proper interface with the data, and is essentially a log of reads
and writes to the ACTUAL DATA. Who at fucking Amazon decided that they should keep a proper database of the data, but make the users only interface with that data
outside of the book it correlates to a fucking TEXT FILE??? WHO THE FUCK??? WHY??? Why not store it as a JSON file that gets updated when the database of a
particular book is changed, and write a display application that reads from that JSON file? If the database is centralized, you could even let users delete
clippings from certain books easily without having to navigate to that particular location in the book (which by the way, you have to do manually, because your
only interface with your clippings is a fucking TEXT FILE WHICH ISNT EVEN SAVED IN A STANDARDIZED FORMAT ACROSS EBOOK TYPES).

IN ADDITION TO THAT, BECAUSE THE MY CLIPPINGS FILE IS ESSENTIALLY A CHANGELOG, IT ISNT FUCKING SYNCHRONIZED. IF A USER DELETES A CLIPPING IN BOOK, THE FILE IS
UNCHANGED. IF SOMEONE DELETES A CLIPPING IN FILE, THE BOOK IS UNCHANGED.

WHO THE FUCK DESIGNED THIS TRAVESTY????

I HOPE A TRUCK FALLS ON THAT ENGINEERS HEAD

Anyway, enjoy 150+ lines of text parsing because Amazon half assed a central feature of their e-reader
*/

const EntrySeperator = "\r\n=========="

const KindleStrfTime = "%A, %d %B %Y %H:%M:%S"

const (
	CLIPPING = iota
	NOTE
	BOOKMARK
)

// before being sorted into a note or a Section, data is read into a single struct which contains all information an entry in the file can possibly contain
type Section struct {
	Content     string
	SectionType int

	//struct nested within struct for storing Notes
	Notes []Section

	//both needed to correctly match notes to highlights
	Page          int
	BookLineStart int
	BookLineEnd   int

	Date time.Time

	BookAuthor string
	BookTitle  string
}

// as of current only supports one author
type book struct {
	authors    string
	name       string
	highlights []Section
}

func init() {
	//make logger print line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Return clippings as an array of strings seperated by ENTRY_SEPERATOR
func splitClippings(clippings string) (sortedClippings []string) {
	clippings = trimEmptyLines([]byte(clippings))
	sortedClippings = strings.Split(clippings, EntrySeperator)
	return sortedClippings
}

func ReadClippingsFileAsSectionArray(file string) []Section {
	return getSections(file)
}

func getSections(content string) []Section {
	clippings := splitClippings(content)
	sectionArray := []Section{}
	for _, c := range clippings {
		trimmedContent := trimEmptyLines([]byte(c))
		contentLines := strings.Split(trimmedContent, "\n")

		//last line will have a length of 1
		if len(contentLines) < 3 {
			continue
		}

		section := readSection(contentLines)
		if section.SectionType == BOOKMARK {
			continue
		}

		sectionArray = append(sectionArray, section)
	}

	return sectionArray
}

// returns the read section data and a location from which to read the next section
// this is a little horrendous
func readSection(contentLines []string) (output Section) {

	//not yet supporting bookmarks, see no need to
	sectionType := parseClippingType(contentLines[1])
	if sectionType == BOOKMARK {
		return
	}

	book, author := parseTitleLine(contentLines[0])
	page, locationStart, locationEnd, date := parseInfoLine(contentLines[1])
	content := ""

	for _, l := range contentLines[2:] {
		content += l
	}

	return Section{Content: content, BookAuthor: author, BookTitle: book, Page: page, BookLineStart: locationStart, BookLineEnd: locationEnd, Date: date, SectionType: sectionType}
}

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

// - Your Highlight on page 68 | location 1040-1044 | Added on Friday, 3 February 2023 21:36:09
func parseInfoLine(line string) (page int, locationStart int, locationEnd int, date time.Time) {
	trimmedLine, page := parsePage(line)
	trimmedLine, locationStart, locationEnd = parseLocation(trimmedLine)
	date = parseDate(trimmedLine)

	return page, locationStart, locationEnd, date

}

// this is inelegant and hard to read, and it assumes the format of My Clippings better not change
// but if some fucking amazon engineer changes the format of this goddamn text file and DOESN'T MAKE IT JSON OR YAML PLEASE i will give up
func parsePage(line string) (trimmedLine string, page int) {
	pageTextLocation := strings.Index(line, "page")
	dividerLocation := strings.Index(line, "|")
	//cut out all of string before relevant part

	//some formats dont have a page number, catch this
	if pageTextLocation == -1 {
		return line, -1
	}

	pageString := line[pageTextLocation+5 : dividerLocation-1]

	//pdfs have two page numbers, take the first
	hyphenLocation := strings.Index(pageString, "-")
	if hyphenLocation != -1 {
		pageString = pageString[:hyphenLocation-1]
	}

	page, err := strconv.Atoi(pageString)
	if err != nil {
		log.Fatal(err)
	}

	return line[dividerLocation+2:], page
}

func parseLocation(line string) (trimmedLine string, locationStart int, locationEnd int) {
	locationTextLocation := strings.Index(line, "location")
	//PDF file, only page number is given
	if locationTextLocation == -1 {
		return line, -1, -1
	}

	dividerLocation := strings.Index(line, "|")
	substr := line[locationTextLocation:dividerLocation]

	hyphenLocation := strings.Index(substr, "-")

	//handle case with only one location (notes)
	if hyphenLocation == -1 {
		locationStart, err := strconv.Atoi(substr[len("location")+1 : len(substr)-1])
		if err != nil {
			log.Fatal(err)
		}
		return line[dividerLocation+2:], locationStart, -1
	}

	//handle case with two locations (highlights)
	locationStart, err := strconv.Atoi(substr[len("location")+1 : hyphenLocation])
	if err != nil {
		log.Fatal(err)
	}

	locationEnd, err = strconv.Atoi(substr[hyphenLocation+1 : len(substr)-1])
	if err != nil {
		log.Fatal(err)
	}

	return line[dividerLocation+2:], locationStart, locationEnd
}

func parseDate(line string) time.Time {
	dateSubstr := line[len("Added on "):]

	time, err := timefmt.Parse(strings.TrimSpace(dateSubstr), KindleStrfTime)
	if err != nil {
		log.Fatal(err)
	}

	return time
}

// thanks to Hasan Yousef (https://forum.golangbridge.org/t/removing-first-and-last-empty-lines-from-a-string/24285)
func trimEmptyLines(b []byte) string {
	strs := strings.Split(string(b), "\n")
	str := ""
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		str += s + "\n"
	}
	str = strings.TrimSuffix(str, "\n")

	return str
}
