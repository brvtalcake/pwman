package main

import (
	/* "golang.org/x/crypto/ssh/terminal" */
	/* crypt "pwman/src/encryption"
	encrypt "pwman/src/encryption" */

	"log"
	app "pwman/src/app"
)

func main() {
	pwman := app.PWMan_App{}
	pwman.Init()
	for !pwman.Quit {
		if !pwman.Authorized_key {
			pwman.RunEntryForm()
		} else {
			pwman.RunPswdList()
		}
	}

	pwman.ClearAppResources()
	log.Println("Quiting...")
}
