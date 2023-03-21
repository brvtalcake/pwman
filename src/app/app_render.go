package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PWMan_App struct {
	App            *tview.Application
	Key            string
	Authorized_key map[string]bool
}

func (this_app *PWMan_App) Init() *PWMan_App {
	this_app.App = tview.NewApplication()
	this_app.Key = ""
	this_app.Authorized_key = make(map[string]bool)
	this_app.Authorized_key["has_been_false"] = false
	this_app.Authorized_key["current"] = false
	return this_app
}

func (this_app *PWMan_App) ModifyKey(key string) {
	this_app.Key = key
	this_app.CheckKey()
}

func (this_app *PWMan_App) CheckKey() bool {
	this_app.Authorized_key["has_been_false"] = this_app.Authorized_key["current"]
	this_app.Authorized_key["current"] = this_app.VerifyKey()
	return this_app.Authorized_key["current"]
}

func (this_app *PWMan_App) RunEntryForm() {
	this_app.Key = ""

	var form *tview.Form = nil

	quit_func := func() {
		this_app.App.Stop()
	}

	submit_func := func() {
		this_app.ModifyKey(form.GetFormItem(0).(*tview.InputField).GetText())
		this_app.App.SetInputCapture(nil)
		this_app.App.Stop()
		if this_app.CheckKey() {
			quit_func()
		} else {
			if !this_app.Authorized_key["has_been_false"] {
				this_app.Authorized_key["has_been_false"] = true
				this_app.RunEntryForm()
			}
			this_app.Authorized_key["has_been_false"] = true
		}
	}

	submit_on_enter := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			submit_func()
		}
		return event
	}

	form = tview.NewForm().
		AddPasswordField("Enter your key :", "", 35, '*', nil).
		SetFocus(0).
		AddButton("Submit", submit_func).
		AddButton("Quit", quit_func).
		AddTextView("NOTE :", "The key won't be saved anywhere, so you have to remember it.", 0, 0, false, false)

	form.SetBorder(true).SetTitle(" Welcome to PWMan ! ").SetTitleAlign(tview.AlignCenter)

	this_app.App.SetInputCapture(submit_on_enter)

	if err := this_app.App.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (this_app *PWMan_App) RunPswdList() {
	// TODO
}
