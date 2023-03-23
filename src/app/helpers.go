package app

import (
	"log"
	"unicode"
)

func Check(e error) {
	if e != nil {
		log.Panic(e.Error())
	}
}

func (this_app *PWMan_App) ClearAppResources() {
	if this_app.Archive.BZ2_Reader != nil {
		this_app.Archive.BZ2_Reader.Close()
	}
	if this_app.Archive.BZ2_Writer != nil {
		this_app.Archive.BZ2_Writer.Close()
	}
	if this_app.Archive.IO_Reader != nil {
		this_app.Archive.IO_Reader.Close()
	}
	if this_app.Archive.IO_Writer != nil {
		this_app.Archive.IO_Writer.Close()
	}
	// TODO: Clear other resources
}

func IsValidEntryChar(char rune) bool {
	return !unicode.IsControl(char)
}
