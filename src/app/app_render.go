package app

import (
	"fmt"
	"log"
	"os"

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
	IsOnPSWDPopUp      bool
	IsOnGeneralActions bool
	Archive            *PWMan_Archive
}

func (this_app *PWMan_App) Init() *PWMan_App {
	var err error
	this_app.App = tview.NewApplication()
	this_app.IsOnPSWDList = false
	this_app.IsOnPSWDPopUp = false
	this_app.IsOnGeneralActions = false
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

	// create a slice of grids for each entry
	var pages *tview.Pages = tview.NewPages()

	this_app.App = tview.NewApplication()
	list := tview.NewList()
	pages.AddPage(" Password List ", list, true, true)

	var pswd_modals []*tview.Modal
	// var pswd_ok_buttons []*tview.Button

	action_forms := map[string]*tview.Form{"add_entry": nil, "change_key": nil, "change_entry": nil}
	action_forms["add_entry"] = tview.NewForm().
		AddInputField("Entry name :", "", 35, nil, nil).
		AddPasswordField("Entry password :", "", 35, '*', nil).
		AddButton("Submit", func() {
			if len(action_forms["add_entry"].GetFormItem(0).(*tview.InputField).GetText()) > 1 && len(action_forms["add_entry"].GetFormItem(1).(*tview.InputField).GetText()) > 1 {
				this_app.AddToArchive([]string{action_forms["add_entry"].GetFormItem(0).(*tview.InputField).GetText(), action_forms["add_entry"].GetFormItem(1).(*tview.InputField).GetText()})
			}
			this_app.App.Stop() // stop the app to refresh the list. The main loop will restart it
		})

	if entries != nil {
		var i rune = 'a'
		for p, entry := range entries {
			list.AddItem(entry[0], "Press enter to see the "+entry[0]+" associated password.", i, nil) // the selected entry is handled by the function below
			pswd_modals = append(pswd_modals, tview.NewModal().AddButtons([]string{"OK"}).SetDoneFunc(nil).SetText("Password : "+entry[1]))
			pages.AddPage(fmt.Sprintf(" %d ", p), pswd_modals[p], true, false)
			i++
		}
	} else {
		list.AddItem("No entry found", "Press CTRL + A to add one.", 'a', nil)
	}

	var delete_pswd_popup_text string
	if entries == nil {
		delete_pswd_popup_text = "Are you sure you want to delete this entry ?"
	} else {
		delete_pswd_popup_text = "Are you sure you want to delete " + entries[list.GetCurrentItem()][0] + " ?"
	}

	delete_pswd_popup := tview.NewModal().AddButtons([]string{"Yes", "No"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			if !this_app.DeleteFromArchive(entries[list.GetCurrentItem()][0]) {
				log.Println("Could not delete entry")
			}
			this_app.App.Stop()
		} else {
			this_app.IsOnPSWDList = true
			pages.SwitchToPage(" Password List ")
			this_app.App.Stop()
		}
		pages.HidePage(" Delete Entry ")
		this_app.App.Stop()
	}).SetText(delete_pswd_popup_text)

	pages.AddPage(" Delete Entry ", delete_pswd_popup, true, false)

	custom_key_event_handler := func(input *tcell.EventKey) *tcell.EventKey {
		if input.Key() == tcell.KeyCtrlA {
			this_app.IsOnPSWDList = false
			pages.SwitchToPage(" Add Entry ")
			return nil
		} else if input.Key() == tcell.KeyEscape {
			if this_app.IsOnPSWDList {
				this_app.Quit = true
				this_app.App.Stop()
			} else {
				this_app.IsOnPSWDList = true
				pages.SwitchToPage(" Password List ")
			}
			return nil
		} else if input.Key() == tcell.KeyCtrlC {
			this_app.Quit = true
			this_app.App.Stop()
			this_app.ClearAppResources()
			log.Println("SIGINT received, exiting...")
			os.Exit(0)
			return nil
		} else if input.Key() == tcell.KeyCtrlD {
			if entries != nil {
				this_app.IsOnPSWDList = false
				pages.ShowPage(" Delete Entry ")
				this_app.IsOnPSWDList = false
				return nil
			} else {
				return input
			}
		} else if input.Key() == tcell.KeyCtrlK { // TODO: change key
			// TO BE IMPLEMENTED
			return nil
		} else if input.Key() == tcell.KeyCtrlE { // TODO: change entry
			// TO BE IMPLEMENTED
			return nil
		} else if input.Key() == tcell.KeyLeft && this_app.IsOnPSWDList {
			list.SetCurrentItem(list.GetCurrentItem() - 10)
			return nil
		} else if input.Key() == tcell.KeyRight && this_app.IsOnPSWDList {
			list.SetCurrentItem(list.GetCurrentItem() + 10)
			return nil
		} else if input.Key() == tcell.KeyEnter {
			if this_app.IsOnPSWDList {
				selected_item := list.GetCurrentItem()
				pages.ShowPage(fmt.Sprintf(" %d ", selected_item))
				this_app.IsOnPSWDList = false
				this_app.IsOnPSWDPopUp = true
				return nil
			} else if this_app.IsOnPSWDPopUp {
				this_app.IsOnPSWDPopUp = false
				this_app.IsOnPSWDList = true
				pages.HidePage(fmt.Sprintf(" %d ", list.GetCurrentItem()))
			} else {
				return input
			}
		}
		return input
	}
	this_app.App.SetInputCapture(custom_key_event_handler)
	/* previously_selected := 0 */
	custom_mouse_event_handler := func(event *tcell.EventMouse, mouse_action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		/* if this_app.IsOnPSWDList && mouse_action == tview.MouseLeftDoubleClick && previously_selected != list.GetCurrentItem() {
			previously_selected = list.GetCurrentItem()
			pages.ShowPage(fmt.Sprintf(" %d ", previously_selected))
			this_app.IsOnPSWDList = false
			this_app.IsOnPSWDPopUp = true
			return nil, mouse_action
		} */
		if this_app.IsOnPSWDPopUp {
			for p := range entries {
				ok_button_x, ok_button_y, ok_button_w, ok_button_h := pswd_modals[p].GetInnerRect()
				event_x, event_y := event.Position()
				if (mouse_action == tview.MouseLeftClick || mouse_action == tview.MouseLeftDoubleClick) && event_x >= ok_button_x && event_x <= ok_button_x+ok_button_w && event_y >= ok_button_y && event_y <= ok_button_y+ok_button_h {
					this_app.IsOnPSWDPopUp = false
					this_app.IsOnPSWDList = true
					pages.HidePage(fmt.Sprintf(" %d ", list.GetCurrentItem()))
					return nil, mouse_action
				}
			}
		} /* else if this_app.IsOnPSWDList {
			for p, _ := range entries {
			list_element_x, list_element_y, list_element_w, list_element_h := list.Get
			// TO BE IMPLEMENTED
		} */
		return event, mouse_action
	}
	this_app.App.SetMouseCapture(custom_mouse_event_handler)

	pages.AddPage(" Add Entry ", action_forms["add_entry"], true, false)

	if err := this_app.App.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
