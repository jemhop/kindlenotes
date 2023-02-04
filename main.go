package main

func main() {

	content := openClippingsFile(getMountedKindle())
	getSections(content)
	//clippings := openFile(getMountedKindle())
	//book := readClippingsFile(clippings)

	//println(len(book))
}
