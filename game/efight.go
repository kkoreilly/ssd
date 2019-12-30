// Copyright (c) 2019, The EFight Authors. All rights reserved.

package main

import (
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/ki/ki"
)

func main() {
	gimain.Main(func() {
		mainrun()
	})
}

var signUpResult *gi.Label
var logInResult *gi.Label
var inspectText *gi.Label
var tv *gi.TabView
var SUPERMODE = false
var signUpTab *gi.Frame
var homeTab *gi.Frame

func mainrun() {
	data()

	width := 1024
	height := 768

	win := gi.NewMainWindow("efight-main", "EFight Home Screen", width, height) // pixel sizes

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()

	mfr := win.SetMainFrame()
	rec := ki.Node{}
	rec.InitName(&rec, "rec")

	toprow := gi.AddNewFrame(mfr, "toprow", gi.LayoutVert)
	toprow.SetStretchMaxWidth()

	toprow.SetProp("background-color", "lightgreen")
	mainHeader := gi.AddNewLabel(toprow, "mainHeader", "Welcome to EFight version 0.0.0 alpha")
	mainHeader.SetProp("font-size", "100px")
	mainHeader.SetProp("text-align", "center")
	mainHeader.SetProp("font-family", "Times New Roman, serif")

	tv = mfr.AddNewChild(gi.KiT_TabView, "tv").(*gi.TabView)
	tv.NewTabButton = false
	tv.SetStretchMaxWidth()
	tv.SetProp("background-color", "lightgreen")

	signUpTab = tv.AddNewTab(gi.KiT_Frame, "Sign Up").(*gi.Frame)
	startGame()

	signUpTab.Lay = gi.LayoutVert
	signUpTab.SetStretchMaxWidth()
	signUpTab.SetStretchMaxHeight()
	signUpTitle := signUpTab.AddNewChild(gi.KiT_Label, "signUpTitle").(*gi.Label)
	signUpTitle.SetProp("font-size", "x-large")
	signUpTitle.SetProp("text-align", "center")
	signUpTitle.Text = "<b>Enter your information to sign up for EFight:</b>"
	signUpText := signUpTab.AddNewChild(gi.KiT_TextField, "signUpText").(*gi.TextField)
	signUpText.SetProp("horizontal-align", gi.AlignCenter)
	signUpText.Placeholder = "Enter what you want your username to be"
	signUpText.SetStretchMaxWidth()
	signUpText2 := signUpTab.AddNewChild(gi.KiT_TextField, "signUpText2").(*gi.TextField)
	signUpText2.SetProp("horizontal-align", gi.AlignCenter)
	signUpText2.Placeholder = "Enter what you want your password to be"
	signUpText2.SetStretchMaxWidth()

	signUpButton := signUpTab.AddNewChild(gi.KiT_Button, "signUpButton").(*gi.Button)
	signUpButton.Text = "<b>Sign Up!</b>"
	signUpButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			username := signUpText.Text()
			password := signUpText2.Text()
			// fmt.Printf("User: %v Password: %v \n", username, password)
			addUser(username, password)
		}
	})

	signUpResult = signUpTab.AddNewChild(gi.KiT_Label, "signUpResult").(*gi.Label)
	signUpResult.Text = "                                   "
	signUpResult.Redrawable = true

	logInTab := tv.AddNewTab(gi.KiT_Frame, "Log In").(*gi.Frame)

	logInTab.Lay = gi.LayoutVert
	logInTab.SetStretchMaxWidth()
	logInTab.SetStretchMaxHeight()
	logInTitle := logInTab.AddNewChild(gi.KiT_Label, "logInTitle").(*gi.Label)
	logInTitle.SetProp("font-size", "x-large")
	logInTitle.SetProp("text-align", "center")
	logInTitle.Text = "<b>Enter your information to log into EFight:</b>"

	logInText := logInTab.AddNewChild(gi.KiT_TextField, "logInText").(*gi.TextField)
	logInText.SetProp("horizontal-align", gi.AlignCenter)
	logInText.Placeholder = "Username"
	logInText.SetStretchMaxWidth()
	logInText2 := logInTab.AddNewChild(gi.KiT_TextField, "logInText2").(*gi.TextField)
	logInText2.SetProp("horizontal-align", gi.AlignCenter)
	logInText2.Placeholder = "Password"
	logInText2.SetStretchMaxWidth()

	logInButton := logInTab.AddNewChild(gi.KiT_Button, "logInButton").(*gi.Button)
	logInButton.Text = "<b>Log In</b>"
	logInButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			username := logInText.Text()
			password := logInText2.Text()
			// fmt.Printf("User: %v Password: %v \n", username, password)
			logIn(username, password)
		}
	})

	logInResult = logInTab.AddNewChild(gi.KiT_Label, "logInResult").(*gi.Label)
	logInResult.Text = "                                                                                                                                                                                  "
	logInResult.Redrawable = true

	if SUPERMODE == true {

		inspectTab := tv.AddNewTab(gi.KiT_Frame, "Inspect Tab").(*gi.Frame)

		inspectTab.Lay = gi.LayoutVert

		inspectText = inspectTab.AddNewChild(gi.KiT_Label, "inspectText").(*gi.Label)
		inspectText.Redrawable = true
		inspectText.SetStretchMaxWidth()
		initInspect()
	}

	tv.SelectTabIndex(0)
	//
	// 	// main menu
	// 	appnm := oswin.TheApp.Name()
	// 	mmen := win.MainMenu
	// 	mmen.ConfigMenus([]string{appnm, "Edit", "Window"})
	//
	// 	amen := win.MainMenu.ChildByName(appnm, 0).(*gi.Action)
	// 	amen.Menu = make(gi.Menu, 0, 10)
	// 	amen.Menu.AddAppMenu(win)
	//
	// 	emen := win.MainMenu.ChildByName("Edit", 1).(*gi.Action)
	// 	emen.Menu = make(gi.Menu, 0, 10)
	// 	emen.Menu.AddCopyCutPaste(win)
	//
	// 	win.OSWin.SetCloseCleanFunc(func(w oswin.Window) {
	// 		go oswin.TheApp.Quit() // once main window is closed, quit
	// 	})
	//
	// win.MainMenuUpdated()

	vp.UpdateEndNoSig(updt)

	win.StartEventLoop()
}

func initMainTab() {
	updt := tv.UpdateStart()
	tv.SetFullReRender()
	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	homeTab = tv.AddNewTab(gi.KiT_Frame, "Home Tab").(*gi.Frame)

	homeTab.Lay = gi.LayoutVert
	homeTab.SetStretchMaxWidth()
	homeTab.SetStretchMaxHeight()

	mainTitle := homeTab.AddNewChild(gi.KiT_Label, "mainTitle").(*gi.Label)
	mainTitle.SetProp("font-size", "x-large")
	mainTitle.SetProp("font-family", "Times New Roman, serif")
	mainTitle.SetProp("text-align", "center")
	mainTitle.Text = "<b>Welcome to EFight, an Energy Based 3D battle game!</b>"
	playButton := homeTab.AddNewChild(gi.KiT_Button, "playButton").(*gi.Button)
	playButton.Text = "<b>Play!</b>"

	playButton.SetProp("horizontal-align", gi.AlignCenter)

	playButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
		//	startGame()
		}
	})
	homeTab.SetProp("background-color", "lightgreen")
	// tv.SetStretchMaxWidth()
	tv.UpdateEnd(updt)
}