package kindle

import "log"

//This file will handle sorting notes into clippings, the removal of duplicates, and the sorting of clippings into arrays

func SortSections(sections []Section) (books []Book) {
	sections = removeDuplicates(sections)
	sections = sortNotesIntoHighlights(sections)
	books = sortIntoBooks(sections)
	return books
}

// takes unsorted sections and returns a list of top level sections with their children sorted into them (highlights and parentless notes)
func sortNotesIntoHighlights(sections []Section) (topLevelSections []Section) {
	//seperate highlights into seperate array
	for i, s := range sections {
		if s.SectionType == HIGHLIGHT {
			topLevelSections = append(topLevelSections, s)
			sections = RemoveUnordered(sections, i)
		}
	}

	for _, s := range sections {
		if s.SectionType == NOTE {
			sectionIndex, err := matchHighlightToNote(topLevelSections, s)
			if err == -1 {
				log.Fatal(err)
			}

			topLevelSections[sectionIndex].Notes = append(topLevelSections[sectionIndex].Notes, s)
		}
	}

	return topLevelSections

}

// Go through the array of highlights to find one which matches the note, and return its index
func matchHighlightToNote(highlights []Section, note Section) (index int, err int) {
	for i, h := range highlights {
		if (h.BookLineEnd == note.BookLineStart) && h.BookTitle == note.BookTitle {
			return i, 0
		}
	}

	return -1, -1
}

func removeDuplicates(sections []Section) (output []Section) {
	for _, s := range sections {
		if !isDuplicate(sections, s) {
			output = append(output, s)
		}
	}

	return output
}

// checks if content of section is contained by any other section at the same location
func isDuplicate(sections []Section, section Section) bool {
	sectionLength := section.BookLineEnd - section.BookLineStart

	//if there is another section at the same start lovation of a longer length, the section is an unecessary duplicate
	for _, s := range sections {
		if s.BookLineStart == section.BookLineStart {
			if s.BookLineEnd-s.BookLineStart > sectionLength {
				return true
			}
		}
	}

	return false
}

// sorts highlights and top level notes into a book array
func sortIntoBooks(sections []Section) (books []Book) {
	for _, s := range sections {
		bookIndex := bookExists(s.BookTitle, s.BookAuthor, books)

		if bookIndex == -1 {
			books = append(books, Book{Author: s.BookAuthor, Title: s.BookTitle, Highlights: []Section{s}})
			continue
		}

		books[bookIndex].Highlights = append(books[bookIndex].Highlights, s)
	}

	return books
}

// check if a book with a given author and title already exists and return the index
func bookExists(title string, author string, books []Book) int {
	for i, b := range books {
		if b.Author == author && b.Title == title {
			return i
		}
	}
	return -1
}

// removes an elelment from the array without caring about the order the array is in
// much faster than maintaining array order (no shifting of elements)
func RemoveUnordered[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
