package main

import (
	/* "golang.org/x/crypto/ssh/terminal" */
	/* crypt "pwman/src/encryption"
	encrypt "pwman/src/encryption" */
	app "pwman/src/app"
)

func main() {
	pwman := app.PWMan_App{}
	pwman.Init()
	pwman.RunEntryForm()
}
