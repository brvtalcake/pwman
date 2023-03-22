package app

import (
	"io"
	"log"
	"os"

	bz "github.com/dsnet/compress/bzip2"
	"github.com/rivo/tview"
)

type PWMan_Archive struct {
	Location         string
	IO_Reader        io.Reader
	IO_Writer        io.Writer
	BZ2_Reader       *bz.Reader
	BZ2_Writer       *bz.Writer
	EncryptedContent []byte
	EntryCount       uint64
}

type PWMan_App struct {
	App                *tview.Application
	Key                string
	Authorized_key     bool
	Previous_key_false bool
	Quit               bool
	Archive            *PWMan_Archive
}

func (this_app *PWMan_App) Init() *PWMan_App {
	this_app.App = tview.NewApplication()
	this_app.Key = ""
	this_app.Authorized_key = false
	this_app.Previous_key_false = false
	this_app.Quit = false
	this_app.Archive = new(PWMan_Archive)
	ret, err_path := this_app.Archive.GetPhysicalArchivePath()
	Check(err_path)
	if ret == "" {
		log.Panic("Could not get archive path")
		this_app.Quit = true
		this_app.Authorized_key = true // To avoid entering the key form again
		this_app.App.Stop()
		goto func_end
	}
	if have_archive, err := LookForArchive(); !have_archive {
		log.Println("Archive not found! Creating it...")
		if err != nil {
			this_app.Quit = true
			this_app.Authorized_key = true // To avoid entering the key form again
			this_app.App.Stop()
			log.Panic(err.Error())
			goto func_end
		}

		_, err = os.Create(this_app.Archive.Location)
		if err != nil {
			this_app.Quit = true
			this_app.Authorized_key = true // To avoid entering the key form again
			this_app.App.Stop()
			log.Panic(err.Error())
			goto func_end
		}
		os.OpenFile(this_app.Archive.Location, os.O_RDWR, 0755)
	} else {
		log.Println("Archive found! Loading it...")
		if err != nil {
			this_app.Quit = true
			this_app.Authorized_key = true // To avoid entering the key form again
			this_app.App.Stop()
			log.Panic(err)
			goto func_end
		}
	}

func_end:
	return this_app
}

func (this_app *PWMan_App) ModifyKey(key string) {
	this_app.Key = key
	this_app.CheckKey()
}

func (this_app *PWMan_App) CheckKey() bool {
	this_app.Authorized_key = this_app.VerifyKey()
	return this_app.Authorized_key
}

func (this_app *PWMan_App) RunEntryForm() {
	this_app.Key = ""

	var form *tview.Form = nil

	quit_func := func() {
		this_app.Authorized_key = true
		this_app.Quit = true
		this_app.App.Stop()
	}

	submit_func := func() {
		this_app.ModifyKey(form.GetFormItem(0).(*tview.InputField).GetText())
		/* this_app.App.SetInputCapture(nil) */
		if this_app.CheckKey() {
			this_app.App.Stop()
			this_app.RunPswdList()
		} else {
			println("Wrong key !")
			this_app.Previous_key_false = true
			this_app.App.Stop()
		}
	}

	/* submit_on_enter := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			submit_func()
		}
		return event
	} */

	form = tview.NewForm().
		AddPasswordField("Enter your key :", "", 35, '*', nil).
		SetFocus(0).
		AddButton("Submit", submit_func).
		AddButton("Quit", quit_func).
		AddTextView("NOTE :", "The key won't be saved anywhere, so you have to remember it.", 0, 0, false, false)

	if this_app.Previous_key_false {
		form.AddTextView("ERROR :", "Wrong key !", 0, 0, false, false)
	}

	form.SetBorder(true).SetTitle(" Welcome to PWMan ! ").SetTitleAlign(tview.AlignCenter)

	/* this_app.App.SetInputCapture(submit_on_enter) */
	this_app.Authorized_key = false
	if err := this_app.App.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (this_app *PWMan_App) RunPswdList() {
	// TODO
}
