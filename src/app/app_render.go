package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"pwman/src/encryption"

	bz "github.com/dsnet/compress/bzip2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PWMan_Archive struct {
	Location         string
	IO_Reader        *os.File
	IO_Writer        *os.File
	BZ2_Reader       *bz.Reader
	BZ2_Writer       *bz.Writer
	EncryptedContent []byte
	DecryptedContent string
	PswEntryCount    uint64
	WasHere          bool
}

type PWMan_App struct {
	App                *tview.Application
	Key                string
	Byte_key           []byte
	Authorized_key     bool
	Previous_key_false bool
	Quit               bool
	IsOnPSWDList       bool
	IsOnPSWDPage       bool
	Archive            *PWMan_Archive
}

func (this_app *PWMan_App) Init() *PWMan_App {
	var err error
	this_app.App = tview.NewApplication()
	this_app.IsOnPSWDList = false
	this_app.IsOnPSWDPage = false
	this_app.Archive = new(PWMan_Archive)
	this_app.Archive.DecryptedContent = ""
	this_app.Archive.EncryptedContent = nil
	this_app.Archive.PswEntryCount = 0
	this_app.Key = ""
	this_app.Byte_key = nil
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
		this_app.Archive.WasHere = false

		_, err = os.Create(this_app.Archive.Location)
		if err != nil {
			this_app.Quit = true
			this_app.Authorized_key = true // To avoid entering the key form again
			this_app.App.Stop()
			log.Panic(err.Error())
			goto func_end
		}

	} else {
		log.Println("Archive found! Loading it...")
		this_app.Archive.WasHere = true
		if err != nil {
			this_app.Quit = true
			this_app.Authorized_key = true // To avoid entering the key form again
			this_app.App.Stop()
			log.Panic(err)
			goto func_end
		}
	}
	this_app.Archive.IO_Reader, err = os.OpenFile(this_app.Archive.Location, os.O_RDONLY, 0755)
	if err != nil {
		this_app.Quit = true
		this_app.Authorized_key = true // To avoid entering the key form again
		this_app.App.Stop()
		log.Panic(err.Error())
		goto func_end
	}
	this_app.Archive.IO_Writer, err = os.OpenFile(this_app.Archive.Location, os.O_WRONLY, 0755)
	if err != nil {
		this_app.Quit = true
		this_app.Authorized_key = true // To avoid entering the key form again
		this_app.App.Stop()
		log.Panic(err.Error())
		goto func_end
	}
	this_app.Archive.BZ2_Reader = nil
	this_app.Archive.BZ2_Writer = nil

func_end:
	return this_app
}

func (this_app *PWMan_App) ModifyKey(key string) {
	this_app.Key = ""
	this_app.Byte_key = nil
	this_app.Key = key
	this_app.CheckKey()
}

func (this_app *PWMan_App) CheckKey() bool {
	this_app.Authorized_key = this_app.VerifyKey()
	return this_app.Authorized_key
}

func (this_app *PWMan_App) RunEntryForm() {

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
			log.Println("Key accepted !")
			//this_app.RunPswdList()
		} else {
			this_app.Previous_key_false = true
			this_app.App.Stop()
			log.Println("Wrong key !")
		}
	}

	/* submit_on_enter := func(input *tcell.EventKey) *tcell.EventKey {
		if input.Key() == tcell.KeyEnter {
			submit_func()
		}
		return input
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
	if !this_app.Archive.WasHere {
		this_app.CreateHeader()
	}
	entries := this_app.ParseArchive()
	this_app.IsOnPSWDList = true

	// create a slice of boxes for each entry
	var pswd_boxes []*tview.Box = nil

	this_app.App = tview.NewApplication()
	list := tview.NewList().AddItem("General actions", "Press a to add or delete a new entry", 'a', nil)
	pswd_boxes = append(pswd_boxes, tview.NewBox().SetBorder(true).SetTitle("General actions").SetTitleAlign(tview.AlignCenter))
	if entries != nil {
		var i rune = 'b'
		for _, entry := range entries {
			list.AddItem(entry[0], "Press enter to see the "+entry[0]+" associated password and the possible actions.", i, nil)
			// create a box for each entry
			pswd_boxes = append(pswd_boxes, tview.NewBox().SetBorder(true).SetTitle(entry[0]).SetTitleAlign(tview.AlignCenter))
			i++
		}
	}

	pages := tview.NewPages().AddPage(" Password List ", list, true, true)

	for i, box := range pswd_boxes {
		pages.AddPage(fmt.Sprintf(" %d ", i), box, false, false)
	}

	custom_event_handler := func(input *tcell.EventKey) *tcell.EventKey {
		if input.Key() == tcell.KeyEscape {
			if this_app.IsOnPSWDList {
				this_app.Quit = true
				this_app.App.Stop()
			} else {
				this_app.IsOnPSWDList = true
				this_app.App.SetRoot(pages, true)
			}
		}
		return input
	}
	this_app.App.SetInputCapture(custom_event_handler)

	if err := this_app.App.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (this_app *PWMan_App) ParseArchive() [][]string {
	var err error
	returned_entries := make([][]string, 0)
	// slurp the whole encrypted content and store it in archive struct
	if this_app.Archive.BZ2_Writer != nil {
		err = this_app.Archive.BZ2_Writer.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}
	if this_app.Archive.BZ2_Reader == nil {
		this_app.Archive.BZ2_Reader, err = bz.NewReader(this_app.Archive.IO_Reader, &bz.ReaderConfig{})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	} else if this_app.Archive.BZ2_Reader != nil { // close and reopen the reader to be sure
		err = this_app.Archive.BZ2_Reader.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		this_app.Archive.BZ2_Reader, err = bz.NewReader(this_app.Archive.IO_Reader, &bz.ReaderConfig{})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}
	this_app.Archive.BZ2_Reader.Reset(this_app.Archive.IO_Reader)
	_, err = this_app.Archive.IO_Reader.Seek(0, io.SeekStart)
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	this_app.Archive.EncryptedContent, err = io.ReadAll(this_app.Archive.BZ2_Reader) // maybe to change
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	this_app.Archive.DecryptedContent, err = encryption.Decrypt(this_app.Archive.EncryptedContent, this_app.Byte_key)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			log.Panic("Wrong key !")
			log.Println("\n\x1b[34;1mError :\n\x1b[0m" + err.Error())
			recover()
		} else {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}

	runed_decrypted_content := []rune(this_app.Archive.DecryptedContent)
	runed_decrypted_content_len := len(runed_decrypted_content)
	no_more_entries := false
	global_cursor := uint64(0)
	entry_count := uint64(0)

	check_if_no_more_entries := func() bool {
		return global_cursor >= uint64(runed_decrypted_content_len)-1
	}
	get_next_entry := func() []string {
		next_entry := make([]string, 2)
		no_more_entries = check_if_no_more_entries()
		for !IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		for IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			next_entry[0] += string(runed_decrypted_content[global_cursor])
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		for !IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		for IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			next_entry[1] += string(runed_decrypted_content[global_cursor])
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		entry_count++
		if entry_count > 1 { // skip the header
			return next_entry
		} else {
			return nil
		}

	}

	// split the decrypted content into entries
	/*
	* Each entries has the following format :
	* 1.	Name of the entry (website, app, etc.)
	* 2.	Associated password
	* 3.	2 newlines
	 */

	for global_cursor < uint64(runed_decrypted_content_len) && !no_more_entries {
		returned_entries = append(returned_entries, get_next_entry())
	}
	this_app.Archive.PswEntryCount = entry_count
	err = this_app.Archive.BZ2_Reader.Close()
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	if entry_count > 1 { // skip the header
		return returned_entries
	} else {
		return nil
	}
}

/* func (this_app *PWMan_App) ReshapeArchive() { */

func (this_app *PWMan_App) AddToArchive(entry []string) {
	var err error
	if this_app.Archive.BZ2_Reader != nil {
		err = this_app.Archive.BZ2_Reader.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}
	if this_app.Archive.BZ2_Writer == nil {
		err = os.Truncate(this_app.Archive.Location, 0)
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		this_app.Archive.BZ2_Writer, err = bz.NewWriter(this_app.Archive.IO_Writer, &bz.WriterConfig{Level: bz.BestCompression})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	} else if this_app.Archive.BZ2_Writer != nil {
		err = this_app.Archive.BZ2_Writer.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		err = os.Truncate(this_app.Archive.Location, 0)
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		this_app.Archive.BZ2_Writer, err = bz.NewWriter(this_app.Archive.IO_Writer, &bz.WriterConfig{Level: bz.BestCompression})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}

	this_app.Archive.BZ2_Writer.Reset(this_app.Archive.IO_Writer)
	_, err = this_app.Archive.IO_Writer.Seek(0, io.SeekStart)
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	var new_entry string = "\n" + entry[0] + "\n" + entry[1] + "\n"
	this_app.Archive.DecryptedContent += new_entry
	this_app.Archive.EncryptedContent, err = encryption.Encrypt([]byte(this_app.Archive.DecryptedContent), this_app.Byte_key)
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}

	_, err = io.WriteString(this_app.Archive.BZ2_Writer, string(this_app.Archive.EncryptedContent))
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	this_app.Archive.PswEntryCount++
	err = this_app.Archive.BZ2_Writer.Close()
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
}
