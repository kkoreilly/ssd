// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	// "fmt"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/svg"
	"github.com/goki/gi/units"
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
var teamMainText *gi.Label
var tv *gi.TabView

// var SUPERMODE = false
var signUpTab *gi.Frame
var homeTab *gi.Frame
var aboutTab *gi.Frame
var playTab *gi.Frame
var resourcesTab *gi.Frame
var map2dTab *gi.Frame

// var map3dTab *gi.Frame // to be added later
var teamTab *gi.Frame
var goldResourcesText *gi.Label
var livesResourcesText *gi.Label
var tbrowH *gi.Layout
var tbrowR *gi.Layout
var keyRow *gi.Frame
var win *gi.Window
var currentTrainingMap string
var currentMap Map
var currentMapString string
var simulateText *gi.Label
var mapSVG *svg.SVG

func mainrun() {
	data() // Connect to data base

	width := 1024 // pixel sizes of screen
	height := 768 // pixel sizes of screen

	win = gi.NewMainWindow("singularity-showdown-main", "Singularity Showdown Home Screen", width, height)

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()

	mfr := win.SetMainFrame()
	rec := ki.Node{}
	rec.InitName(&rec, "rec")

	toprow := gi.AddNewFrame(mfr, "toprow", gi.LayoutVert)
	toprow.SetStretchMaxWidth()

	toprow.SetProp("background-color", "black")
	mainHeaderText := `<b>Welcome to <span style="color:grey">Singularity</span> <span style="color:red">Showdown</span> version 0.0.0 pre-alpha</b>`
	mainHeader := gi.AddNewLabel(toprow, "mainHeader", mainHeaderText)
	mainHeader.SetProp("font-size", "90px")
	mainHeader.SetProp("text-align", "center")
	mainHeader.SetProp("font-family", "Times New Roman, serif")
	mainHeader.SetProp("color", "white")

	tv = mfr.AddNewChild(gi.KiT_TabView, "tv").(*gi.TabView) // Create main tab view
	tv.NewTabButton = false
	tv.NoDeleteTabs = true
	tv.SetStretchMaxWidth()

	signUpTab = tv.AddNewTab(gi.KiT_Frame, "Sign Up").(*gi.Frame)

	signUpTab.Lay = gi.LayoutVert
	signUpTab.SetStretchMaxWidth()
	signUpTab.SetStretchMaxHeight()

	signUpTitle := signUpTab.AddNewChild(gi.KiT_Label, "signUpTitle").(*gi.Label)
	signUpTitle.SetProp("font-size", "x-large")
	signUpTitle.SetProp("text-align", "center")
	signUpTitle.Text = "<b>Enter your information to sign up for Singularity Showdown:</b>"
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

	byPassButton := signUpTab.AddNewChild(gi.KiT_Button, "byPassButton").(*gi.Button)
	byPassButton.Text = "<b>Log in with tester account</b>"
	byPassButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			USER = "tester"
			PASSWORD = "1234"

			tv.DeleteTabIndex(0, true)
			tv.DeleteTabIndex(0, true)
			initMainTabs()

			tv.SelectTabIndex(0)
		}
	})

	signUpResult = signUpTab.AddNewChild(gi.KiT_Label, "signUpResult").(*gi.Label)
	signUpResult.Text = "                                   "
	signUpResult.Redrawable = true

	playButton := signUpTab.AddNewChild(gi.KiT_Button, "playButton").(*gi.Button)
	playButton.Text = "<b>Play (Tester)</b>"

	playButton.SetProp("horizontal-align", gi.AlignCenter)

	playButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			initPlayTab()
		}
	})

	logInTab := tv.AddNewTab(gi.KiT_Frame, "Log In").(*gi.Frame)

	logInTab.Lay = gi.LayoutVert
	logInTab.SetStretchMaxWidth()
	logInTab.SetStretchMaxHeight()
	logInTitle := logInTab.AddNewChild(gi.KiT_Label, "logInTitle").(*gi.Label)
	logInTitle.SetProp("font-size", "x-large")
	logInTitle.SetProp("text-align", "center")
	logInTitle.Text = "<b>Enter your information to log into Singularity Showdown:</b>"

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

	// if SUPERMODE == true {
	//
	// 	inspectTab := tv.AddNewTab(gi.KiT_Frame, "Inspect Tab").(*gi.Frame)
	//
	// 	inspectTab.Lay = gi.LayoutVert
	//
	// 	inspectText = inspectTab.AddNewChild(gi.KiT_Label, "inspectText").(*gi.Label)
	// 	inspectText.Redrawable = true
	// 	inspectText.SetStretchMaxWidth()
	// 	initInspect()
	// }

	tv.SelectTabIndex(0)
	tv.ChildByName("tabs", 0).SetProp("background-color", "darkgrey")
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

func initMainTabs() {
	updt := tv.UpdateStart()
	tv.SetFullReRender()

	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	homeTab = tv.AddNewTab(gi.KiT_Frame, "<b>Home</b>").(*gi.Frame)

	homeTab.Lay = gi.LayoutVert
	homeTab.SetStretchMaxWidth()
	homeTab.SetStretchMaxHeight()

	mainTitle := homeTab.AddNewChild(gi.KiT_Label, "mainTitle").(*gi.Label)
	mainTitle.SetProp("font-size", "60px")
	mainTitle.SetProp("font-family", "Times New Roman, serif")
	mainTitle.SetProp("text-align", "center")
	mainTitle.Text = "Welcome to Singularity Showdown, a strategic 3D Battle Game"

	// playButton := homeTab.AddNewChild(gi.KiT_Button, "playButton").(*gi.Button)
	// playButton.Text = "<b>Play (Tester Mode)</b>"
	//
	// playButton.SetProp("horizontal-align", gi.AlignCenter)
	//
	// playButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
	// 	if sig == int64(gi.ButtonClicked) {
	// 		initPlayTab()
	// 	}
	// })
	trow := gi.AddNewLayout(homeTab, "trainingRow", gi.LayoutHoriz)
	trow.SetProp("spacing", units.NewEx(2))
	trow.SetProp("horizontal-align", gi.AlignLeft)
	trow.SetStretchMaxWidth()

	trainingText := gi.AddNewLabel(trow, "trainingRowText", "Practice and level up in Training Mode:")
	trainingText.SetProp("font-size", "30px")

	trainingDropdown := gi.AddNewMenuButton(trow, "trainingDropdown")
	trainingDropdown.SetText("Choose a map to train in")

	for _, value := range AllMaps {
		var value1 = value.Name
		var value2 = value.MapData
		// fmt.Printf("Value (0): %v \n", value.Name)
		trainingDropdown.Menu.AddAction(gi.ActOpts{Label: value1},
			win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				// fmt.Printf("Value (1): %v \n", value.Name)
				// fmt.Printf("Value (2): %v \n", value1)
				currentMap = value2 // Set the correct map for this dropdown
				currentMapString = value1
				trainingDropdown.SetText(value1)
				// fmt.Printf("Value: %v \n", value1)
			})
	}

	trainingPlayButton := trow.AddNewChild(gi.KiT_Button, "trainingPlayButton").(*gi.Button)
	trainingPlayButton.Text = "<b>Play in Training Mode</b>"
	trainingPlayButton.SetProp("background-color", "orange")

	trainingPlayButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			initPlayTab()
		}
	})

	homeTab.SetProp("background-color", "lightblue")

	resourcesTab = tv.AddNewTab(gi.KiT_Frame, "<b>Resources</b>").(*gi.Frame)

	resourcesTab.Lay = gi.LayoutVert
	resourcesTab.SetStretchMaxWidth()
	resourcesTab.SetStretchMaxHeight()
	resourcesTab.SetProp("background-color", "lightblue")

	resourcesTitle := resourcesTab.AddNewChild(gi.KiT_Label, "resourcesTitle").(*gi.Label)
	resourcesTitle.SetProp("font-size", "60px")
	resourcesTitle.SetProp("font-family", "Times New Roman, serif")
	resourcesTitle.SetProp("text-align", "center")
	resourcesTitle.Text = "Your Resources:"

	goldResourcesText = resourcesTab.AddNewChild(gi.KiT_Label, "goldResourcesText").(*gi.Label)
	goldResourcesText.SetProp("font-size", "30px")
	goldResourcesText.SetProp("font-family", "Times New Roman, serif")
	goldResourcesText.SetProp("text-align", "left")
	goldResourcesText.Text = "                                                                                                                                      "
	goldResourcesText.Redrawable = true

	brow := gi.AddNewLayout(resourcesTab, "gbrow", gi.LayoutHoriz)
	brow.SetProp("spacing", units.NewEx(2))
	brow.SetProp("horizontal-align", gi.AlignLeft)
	brow.SetStretchMaxWidth()

	goldButton := gi.AddNewButton(brow, "goldButton")
	goldButton.Text = "Purchase 100 gold for just 99 cents"
	goldButton.SetProp("background-color", "#D4AF37")
	goldButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			updateResource("gold", GOLD+100)
			goldResourcesText.SetText("                                            ")
			readResources()
		}
	})

	goldButton1 := gi.AddNewButton(brow, "goldButton1")
	goldButton1.Text = "BEST DEAL: Purchase 1000 gold for $8.99"
	goldButton1.SetProp("background-color", "#D4AF37")
	goldButton1.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			updateResource("gold", GOLD+1000)
			goldResourcesText.SetText("                                            ")
			readResources()
		}
	})
	livesResourcesText = resourcesTab.AddNewChild(gi.KiT_Label, "livesResourcesText").(*gi.Label)
	livesResourcesText.SetProp("font-size", "30px")
	livesResourcesText.SetProp("font-family", "Times New Roman, serif")
	livesResourcesText.SetProp("text-align", "left")
	livesResourcesText.Text = "                                                                                                                                      "
	livesResourcesText.Redrawable = true

	livesButton := gi.AddNewButton(resourcesTab, "livesButton")
	livesButton.Text = "Purchase 10 lives for 10 gold"
	livesButton.SetProp("background-color", "pink")
	livesButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			updateResource("gold", GOLD-10)
			updateResource("lives", LIVES+10)
			goldResourcesText.SetText("                                            ")
			livesResourcesText.SetText("                                            ")
			readResources()
		}
	})

	// updateResource("gold", 70)
	readResources()

	map2dTab = tv.AddNewTab(gi.KiT_Frame, "<b>Map (2D)</b>").(*gi.Frame)

	map2dTab.Lay = gi.LayoutVert
	map2dTab.SetStretchMaxWidth()
	map2dTab.SetStretchMaxHeight()
	map2dTab.SetProp("background-color", "lightblue")

	map2dTitle := map2dTab.AddNewChild(gi.KiT_Label, "map2dTitle").(*gi.Label)
	map2dTitle.SetProp("font-size", "60px")
	map2dTitle.SetProp("font-family", "Times New Roman, serif")
	map2dTitle.SetProp("text-align", "center")
	map2dTitle.Text = "Live Map of the World (2D):"

	keyRow = gi.AddNewFrame(map2dTab, "keyRow", gi.LayoutHoriz)
	keyRow.SetProp("spacing", units.NewEx(2))
	keyRow.SetProp("horizontal-align", gi.AlignLeft)
	keyRow.SetProp("background-color", "white")
	keyRow.SetStretchMaxWidth()

	keyMainText := gi.AddNewLabel(keyRow, "keyMainText", "<b>Team Key:</b>")
	keyMainText.SetProp("font-size", "30px")

	addKeyItems()

	simulateButton := gi.AddNewButton(map2dTab, "simulateButton")
	simulateButton.Text = "Simulate (Full)"
	simulateButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			FirstWorldBorders.simulateMap(true)
			FirstWorld.RenderSVGs(mapSVG)
		}
	})
	simulateButton1 := gi.AddNewButton(map2dTab, "simulateButton1")
	simulateButton1.Text = "Simulate (Step)"
	simulateButton1.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			FirstWorldBorders.simulateMap(false)
			FirstWorld.RenderSVGs(mapSVG)
		}
	})
	// renderButton := gi.AddNewButton(map2dTab, "renderButton")
	// renderButton.Text = "Render SVGS"
	// renderButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
	// 	if sig == int64(gi.ButtonClicked) {
	// 		FirstWorld.RenderSVGs(mapSVG)
	// 	}
	// })
	resetButton := gi.AddNewButton(map2dTab, "resetButton")
	resetButton.Text = "Reset"
	resetButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			FirstWorld = OriginFirstWorld
			FirstWorldBorders = OriginFirstWorldBorders
			simulateText.SetText("")
			FirstWorld.RenderSVGs(mapSVG)
		}
	})
	simulateText = gi.AddNewLabel(map2dTab, "simulateText", "                                                                            ")
	simulateText.SetProp("font-size", "20px")
	simulateText.Redrawable = true

	// width := 1024 // pixel sizes of screen
	// height := 768 // pixel sizes of screen

	mapSVG = svg.AddNewSVG(map2dTab, "mapSVG")
	mapSVG.Fill = true
	mapSVG.SetProp("background-color", "white")
	// mapSVG.SetProp("width", units.NewPx(float32(width-20)))
	// mapSVG.SetProp("height", units.NewPx(float32(height-100)))
	mapSVG.SetStretchMaxWidth()
	mapSVG.SetStretchMaxHeight()

	FirstWorld.RenderSVGs(mapSVG)

	// map3dTab = tv.AddNewTab(gi.KiT_Frame, "<b>Map (3D)</b>").(*gi.Frame)
	//
	// map3dTab.Lay = gi.LayoutVert
	// map3dTab.SetStretchMaxWidth()
	// map3dTab.SetStretchMaxHeight()
	// map3dTab.SetProp("background-color", "lightblue")
	//
	// map3dTitle := map3dTab.AddNewChild(gi.KiT_Label, "map3dTitle").(*gi.Label)
	// map3dTitle.SetProp("font-size", "60px")
	// map3dTitle.SetProp("font-family", "Times New Roman, serif")
	// map3dTitle.SetProp("text-align", "center")
	// map3dTitle.Text = "Live Map of the World (3D):"

	teamTab = tv.AddNewTab(gi.KiT_Frame, "<b>Your Team</b>").(*gi.Frame)

	teamTab.Lay = gi.LayoutVert
	teamTab.SetStretchMaxWidth()
	teamTab.SetStretchMaxHeight()
	teamTab.SetProp("background-color", "lightblue")

	teamTitle := teamTab.AddNewChild(gi.KiT_Label, "teamTitle").(*gi.Label)
	teamTitle.SetProp("font-size", "60px")
	teamTitle.SetProp("font-family", "Times New Roman, serif")
	teamTitle.SetProp("text-align", "center")
	teamTitle.Text = "<b>Your Team</b>"

	teamMainText = teamTab.AddNewChild(gi.KiT_Label, "teamMainText").(*gi.Label)
	teamMainText.SetProp("font-size", "30px")
	teamMainText.SetProp("font-family", "Times New Roman, serif")
	// teamMainText.SetProp("text-align", "center")
	teamMainText.Text = ""
	teamMainText.Redrawable = true
	readTeam()

	//if TEAM == "" { // when uncommented -- you can not switch teams. When commented, you can switch teams
	if TEAM == "" {
		teamMainText.SetText(teamMainText.Text + "\n\n<b>Since you have no team right now, you must join a team. Click one of the buttons below to join a team</b>.")
	} else {
		teamMainText.SetText(teamMainText.Text + "\n\n<b>Click one of the buttons below to switch your team<b>.")
	}
	gi.AddNewSpace(teamTab, "space1")
	tbrowH = gi.AddNewLayout(teamTab, "tbrowH", gi.LayoutHoriz)
	tbrowH.SetProp("spacing", units.NewEx(2))
	tbrowH.SetProp("horizontal-align", gi.AlignLeft)
	tbrowH.SetStretchMaxWidth()
	tbrowHText := gi.AddNewLabel(tbrowH, "tBrowHText", "<b>Join a human team:</b>")
	tbrowHText.SetProp("font-size", "30px")

	gi.AddNewSpace(teamTab, "space2")

	tbrowR = gi.AddNewLayout(teamTab, "tbrowR", gi.LayoutHoriz)
	tbrowR.SetProp("spacing", units.NewEx(2))
	tbrowR.SetProp("horizontal-align", gi.AlignLeft)
	tbrowR.SetStretchMaxWidth()

	tbrowRText := gi.AddNewLabel(tbrowR, "tBrowRText", "<b>Join a robot team:</b>")
	tbrowRText.SetProp("font-size", "30px")
	addTeamUpdateButtons()
	//}

	aboutTab = tv.AddNewTab(gi.KiT_Frame, "<b>About</b>").(*gi.Frame)

	aboutTab.Lay = gi.LayoutVert
	aboutTab.SetStretchMaxWidth()
	aboutTab.SetStretchMaxHeight()
	aboutTab.SetProp("background-color", "lightblue")

	aboutTitle := aboutTab.AddNewChild(gi.KiT_Label, "aboutTitle").(*gi.Label)
	aboutTitle.SetProp("font-size", "60px")
	aboutTitle.SetProp("font-family", "Times New Roman, serif")
	aboutTitle.SetProp("text-align", "center")
	aboutTitle.Text = "About Singularity Showdown"

	aboutText := aboutTab.AddNewChild(gi.KiT_Label, "aboutText").(*gi.Label)
	aboutText.SetProp("font-size", "30px")
	aboutText.SetProp("font-family", "Times New Roman, serif")
	aboutText.SetProp("text-align", "left")
	aboutText.Text = "Singularity Showdown is an open source, Strategic 3D Battle Game."

	tv.UpdateEnd(updt)
}

func initPlayTab() {
	updt := tv.UpdateStart()
	tv.SetFullReRender()

	rec := ki.Node{}
	rec.InitName(&rec, "rec")

	if currentMapString == "" { // if no map selected to join
		tv.UpdateEnd(updt)
		return // then don't create the game
	}
	_, err := tv.TabByNameTry("<b>Game</b>") // check if the game tab already exists -- there will not be an error if it already exists

	if err == nil { // if the tab Game already exists
		tv.SelectTabByName("<b>Game</b>")
		tv.UpdateEnd(updt)
		return // and don't create a new tab
	}

	playTab = tv.AddNewTab(gi.KiT_Frame, "<b>Game</b>").(*gi.Frame)

	playTab.Lay = gi.LayoutVert
	playTab.SetStretchMaxWidth()
	playTab.SetStretchMaxHeight()

	playTitleText := gi.AddNewLabel(playTab, "playTitleText", "Welcome to")
	playTitleText.SetText("Welcome to " + currentMapString)
	playTitleText.SetProp("text-align", "center")
	playTitleText.SetProp("font-size", "40px")

	TheGame = &Game{} // Set up game
	TheGame.Config()  // Set up game

	tv.SelectTabByName("<b>Game</b>")
	tv.UpdateEnd(updt)
}
