// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	// "fmt"
	// "time"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
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
var tutorialTab *gi.Frame
var playTab *gi.Frame
var resourcesTab *gi.Frame
var simulationTab *gi.Frame
var simulationControlsTab *gi.Frame
var map2dTab *gi.Frame

// var map3dTab *gi.Frame // to be added later
var teamTab *gi.Frame
var goldResourcesText *gi.Label
var livesResourcesText *gi.Label
var tbrowH *gi.Layout
var tbrowR *gi.Layout
var keyRow *gi.Frame
var keyRow1 *gi.Frame
var keyRowM *gi.Frame
var win *gi.Window
var currentTrainingMap string
var currentMap Map
var currentMapString string
var simulateText *gi.Label
var simMapSVG *svg.SVG
var mapSVG *svg.SVG
var comebacks = false

func mainrun() {
	data()        // Connect to data base
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
	mainHeaderText := `<b>Welcome to <span style="color:grey">Singularity</span> <span style="color:red">Showdown</span> version 0.0.0 Alpha</b>`
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
	// initBorders()
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
			// t0 := time.Now()
			updateResource("gold", GOLD+1000)
			goldResourcesText.SetText("                                            ")
			readResources()
			// t1 := time.Now()
			// d := t1.Sub(t0)
			// fmt.Printf("Time to update: %v \n", d.Milliseconds())
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

	map2dTab = tv.AddNewTab(gi.KiT_Frame, "<b>Live Map of the World</b>").(*gi.Frame)

	map2dTab.Lay = gi.LayoutVert
	map2dTab.SetStretchMaxWidth()
	map2dTab.SetStretchMaxHeight()
	map2dTab.SetProp("background-color", "lightblue")

	map2dTitle := map2dTab.AddNewChild(gi.KiT_Label, "map2dTitle").(*gi.Label)
	map2dTitle.SetProp("font-size", "60px")
	map2dTitle.SetProp("font-family", "Times New Roman, serif")
	map2dTitle.SetProp("text-align", "center")
	map2dTitle.Text = "Live Map of the World"

	keyRowM = gi.AddNewFrame(map2dTab, "keyRowM", gi.LayoutHoriz)
	keyRowM.SetProp("spacing", units.NewEx(2))
	keyRowM.SetProp("horizontal-align", gi.AlignLeft)
	keyRowM.SetProp("background-color", "white")
	keyRowM.SetStretchMaxWidth()

	mapSVG = svg.AddNewSVG(map2dTab, "mapSVG")
	mapSVG.Fill = true
	mapSVG.SetProp("background-color", "white")
	// simMapSVG.SetProp("width", units.NewPx(float32(width-20)))
	// simMapSVG.SetProp("height", units.NewPx(float32(height-100)))
	mapSVG.SetStretchMaxWidth()
	mapSVG.SetStretchMaxHeight()

	FirstWorldLive.RenderSVGs(mapSVG)

	keyMainTextM := gi.AddNewLabel(keyRowM, "keyMainTextM", "<b>Team Key:</b>")
	keyMainTextM.SetProp("font-size", "30px")

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
	aboutText.Text = "Singularity Showdown is an open-source strategic 3D battle game. In Singularity Showdown, AI have become super-intelligent and attacked their human creators. The war has been going on for several months now, and both sides have split up into different groups fighting for control and resources. In Singularity Showdown, you get to choose a side and group to fight for. You fight other teams in 3D battles, the results of which effect the live map of the world. This map is shared by everyone, and the team that takes over the world wins. <br><br>This is Singularity Showdown version 0.0.0 Alpha, and more features will be added with new releases."
	aboutText.SetProp("white-space", gi.WhiteSpaceNormal)
	aboutText.SetProp("max-width", -1)
	aboutText.SetProp("width", "20em")

	tutorialTab = tv.AddNewTab(gi.KiT_Frame, "<b>How to Play</b>").(*gi.Frame)

	tutorialTab.Lay = gi.LayoutVert
	tutorialTab.SetStretchMaxWidth()
	tutorialTab.SetStretchMaxHeight()
	tutorialTab.SetProp("background-color", "lightblue")

	tutorialTitle := tutorialTab.AddNewChild(gi.KiT_Label, "tutorialTitle").(*gi.Label)
	tutorialTitle.SetProp("font-size", "60px")
	tutorialTitle.SetProp("font-family", "Times New Roman, serif")
	tutorialTitle.SetProp("text-align", "center")
	tutorialTitle.Text = "How to Play Singularity Showdown"

	tutorialText := tutorialTab.AddNewChild(gi.KiT_Label, "tutorialText").(*gi.Label)
	tutorialText.SetProp("font-size", "30px")
	tutorialText.SetProp("font-family", "Times New Roman, serif")
	tutorialText.SetProp("text-align", "left")
	tutorialText.Text = "<b>Keyboard Controls during Battles:</b> <br> <br><b>W, A, S, D:</b> Move (Forward, Left, Back, Right) <br><b>Space:</b> Jump <br><b>Move mouse:</b> Rotate screen<br><b>Escape:</b> Toggle Rotate<br><b>Click:</b> Shoot <br> <br><b>Game structure:</b> <br>From the home tab, you can choose to join a battle on a border between territories. Whoever gets to 10 kills first in the 3D Battle wins the battle, which gets their team one point on the border. Once a team gets to 10 points on the border, they take the opponent's territory on the border, which updates the live map for everyone. <br> <br>To get started playing Singularity Showdown, join a team and start battling!"
	tutorialText.SetProp("white-space", gi.WhiteSpaceNormal)
	tutorialText.SetProp("max-width", -1)
	tutorialText.SetProp("width", "20em")

	simulationControlsTab = tv.AddNewTab(gi.KiT_Frame, "<b>Simulation Settings</b>").(*gi.Frame)
	simulationControlsTab.Lay = gi.LayoutVert
	simulationControlsTab.SetStretchMaxWidth()
	simulationControlsTab.SetStretchMaxHeight()
	simulationControlsTab.SetProp("background-color", "lightblue")

	simulationControlsTitle := simulationControlsTab.AddNewChild(gi.KiT_Label, "simulationControlsTitle").(*gi.Label)
	simulationControlsTitle.SetProp("font-size", "60px")
	simulationControlsTitle.SetProp("font-family", "Times New Roman, serif")
	simulationControlsTitle.SetProp("text-align", "center")
	simulationControlsTitle.Text = "Simulation Settings"

	keyRow1 = gi.AddNewFrame(simulationControlsTab, "keyRow1", gi.LayoutHoriz)
	keyRow1.SetProp("spacing", units.NewEx(2))
	keyRow1.SetProp("horizontal-align", gi.AlignLeft)
	keyRow1.SetProp("background-color", "white")
	keyRow1.SetStretchMaxWidth()

	keyMainText1 := gi.AddNewLabel(keyRow1, "keyMainText1", "<b>Team Key:</b>")
	keyMainText1.SetProp("font-size", "30px")
	simulationRandomButton := gi.AddNewButton(simulationControlsTab, "simulationRandomButton")
	simulationRandomButton.Text = "Randomly choose strength"
	simulationRandomButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			InitStrength()
		}
	})
	gi.AddNewSpace(simulationControlsTab, "space2")
	comebackCheckbox := gi.AddNewCheckBox(simulationControlsTab, "comebackCheckbox")
	comebackCheckbox.Text = "Have more comebacks"
	comebackCheckbox.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			if comebacks {
				comebacks = false
			} else {
				comebacks = true
			}
		}
	})
	gi.AddNewSpace(simulationControlsTab, "space3")

	mv := giv.AddNewMapView(simulationControlsTab, "mv")
	mv.SetMap(&TeamStrength)
	mv.SetStretchMaxWidth()
	mv.SetStretchMaxHeight()

	simulationTab = tv.AddNewTab(gi.KiT_Frame, "<b>Simulation</b>").(*gi.Frame)

	simulationTab.Lay = gi.LayoutVert
	simulationTab.SetStretchMaxWidth()
	simulationTab.SetStretchMaxHeight()
	simulationTab.SetProp("background-color", "lightblue")

	simulationTitle := simulationTab.AddNewChild(gi.KiT_Label, "simulationTitle").(*gi.Label)
	simulationTitle.SetProp("font-size", "60px")
	simulationTitle.SetProp("font-family", "Times New Roman, serif")
	simulationTitle.SetProp("text-align", "center")
	simulationTitle.Text = "Simulation of the World"

	keyRow = gi.AddNewFrame(simulationTab, "keyRow", gi.LayoutHoriz)
	keyRow.SetProp("spacing", units.NewEx(2))
	keyRow.SetProp("horizontal-align", gi.AlignLeft)
	keyRow.SetProp("background-color", "white")
	keyRow.SetStretchMaxWidth()

	keyMainText := gi.AddNewLabel(keyRow, "keyMainText", "<b>Team Key:</b>")
	keyMainText.SetProp("font-size", "30px")

	addKeyItems()
	// InitStrength()
	simulationBrow := gi.AddNewFrame(simulationTab, "simulationBrow", gi.LayoutHoriz)
	simulationBrow.SetProp("spacing", units.NewEx(2))
	simulationBrow.SetProp("horizontal-align", gi.AlignLeft)
	simulationBrow.SetProp("background-color", "white")
	simulationBrow.SetStretchMaxWidth()

	simulateButton := gi.AddNewButton(simulationBrow, "simulateButton")
	simulateButton.Text = "Simulate (Full)"
	simulateButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			stopSimulation = false
			go FirstWorldBorders.simulateMap(true)
		}
	})
	simulateButton1 := gi.AddNewButton(simulationBrow, "simulateButton1")
	simulateButton1.Text = "Simulate (Step)"
	simulateButton1.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			FirstWorldBorders.simulateMap(false)
			FirstWorld.RenderSVGs(simMapSVG)
		}
	})
	// renderButton := gi.AddNewButton(simulationTab, "renderButton")
	// renderButton.Text = "Render SVGS"
	// renderButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
	// 	if sig == int64(gi.ButtonClicked) {
	// 		FirstWorld.RenderSVGs(simMapSVG)
	// 	}
	// })
	resetButton := gi.AddNewButton(simulationBrow, "resetButton")
	resetButton.Text = "Reset"
	resetButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			// InitStrength()
			curCountSimulation = 0
			FirstWorld = World{
				"Alaska":         {"Alaska", "human2", "green", ""},
				"Canada":         {"Canada", "human2", "green", ""},
				"USA":            {"USA", "human1", "blue", ""},
				"CentralAmerica": {"CentralAmerica", "human1", "blue", ""},
				"Brazil":         {"Brazil", "human3", "purple", ""},
				"SouthAmerica":   {"SouthAmerica", "human3", "purple", ""},
				"WestAfrica":     {"WestAfrica", "human4", "pink", ""},
				"EastAfrica":     {"EastAfrica", "human4", "pink", ""},
				"Russia":         {"Russia", "human5", "lightgreen", ""},
				"NorthAsia":      {"NorthAsia", "human5", "lightgreen", ""},
				"WestEurope":     {"WestEurope", "robot1", "red", ""},
				"NorthernEurope": {"NorthernEurope", "robot1", "red", ""},
				"EastEurope":     {"EastEurope", "robot2", "orange", ""},
				"MiddleEast":     {"MiddleEast", "robot2", "orange", ""},
				"Australia":      {"Australia", "robot3", "yellow", ""},
				"SouthEastAsia":  {"SouthEastAsia", "robot4", "violet", ""},
				"SouthWestAsia":  {"SouthWestAsia", "robot5", "maroon", ""},
				"Antarctica":     {"Antarctica", "none", "white", ""},
			}
			FirstWorldBorders = Borders{
				"AlaskaRussia":               {"Alaska", "Russia", "battle"},
				"AlaskaCanada":               {"Alaska", "Canada", "human2"},
				"CanadaUSA":                  {"Canada", "USA", "battle"},
				"USACentralAmerica":          {"USA", "CentralAmerica", "human1"},
				"CentralAmericaSouthAmerica": {"CentralAmerica", "SouthAmerica", "battle"},
				"SouthAmericaBrazil":         {"SouthAmerica", "Brazil", "human3"},
				"WestAfricaEastAfrica":       {"WestAfrica", "EastAfrica", "human4"},
				"WestAfricaWestEurope":       {"WestAfrica", "WestEurope", "battle"},
				"WestEuropeNorthernEurope":   {"WestEurope", "NorthernEurope", "robot1"},
				"WestEuropeEastEurope":       {"WestEurope", "EastEurope", "battle"},
				"EastEuropeMiddleEast":       {"EastEurope", "MiddleEast", "robot2"},
				"EastAfricaMiddleEast":       {"EastAfrica", "MiddleEast", "battle"},
				"NorthernEuropeRussia":       {"NorthernEurope", "Russia", "battle"},
				"EastEuropeRussia":           {"EastEurope", "Russia", "battle"},
				"MiddleEastRussia":           {"MiddleEast", "Russia", "battle"},
				"MiddleEastSouthWestAsia":    {"MiddleEast", "SouthWestAsia", "battle"},
				"SouthWestAsiaRussia":        {"SouthWestAsia", "Russia", "battle"},
				"SouthWestAsiaNorthAsia":     {"SouthWestAsia", "NorthAsia", "battle"},
				"NorthAsiaRussia":            {"NorthAsia", "Russia", "human5"},
				"SouthWestAsiaSouthEastAsia": {"SouthWestAsia", "SouthEastAsia", "battle"},
				"NorthAsiaSouthEastAsia":     {"NorthAsia", "SouthEastAsia", "battle"},
				"SouthEastAsiaAustralia":     {"SouthEastAsia", "Australia", "battle"},
				"CanadaWestEurope":           {"Canada", "WestEurope", "battle"},
			}
			// fmt.Printf("First World: %v Origin: %v \n", FirstWorld["USA"].Color, OriginFirstWorld["USA"].Color)
			simulateText.SetText("")
			FirstWorld.RenderSVGs(simMapSVG)
		}
	})

	resetButton1 := gi.AddNewButton(simulationBrow, "resetButton1")
	resetButton1.Text = "Reset Strength"
	resetButton1.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			for k := range TeamStrength {
				TeamStrength[k] = 1
			}
		}
	})

	resetButton2 := gi.AddNewButton(simulationBrow, "resetButton2")
	resetButton2.Text = "Reset All"
	resetButton2.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			for k := range TeamStrength {
				TeamStrength[k] = 1
			}
			curCountSimulation = 0
			FirstWorld = World{
				"Alaska":         {"Alaska", "human2", "green", ""},
				"Canada":         {"Canada", "human2", "green", ""},
				"USA":            {"USA", "human1", "blue", ""},
				"CentralAmerica": {"CentralAmerica", "human1", "blue", ""},
				"Brazil":         {"Brazil", "human3", "purple", ""},
				"SouthAmerica":   {"SouthAmerica", "human3", "purple", ""},
				"WestAfrica":     {"WestAfrica", "human4", "pink", ""},
				"EastAfrica":     {"EastAfrica", "human4", "pink", ""},
				"Russia":         {"Russia", "human5", "lightgreen", ""},
				"NorthAsia":      {"NorthAsia", "human5", "lightgreen", ""},
				"WestEurope":     {"WestEurope", "robot1", "red", ""},
				"NorthernEurope": {"NorthernEurope", "robot1", "red", ""},
				"EastEurope":     {"EastEurope", "robot2", "orange", ""},
				"MiddleEast":     {"MiddleEast", "robot2", "orange", ""},
				"Australia":      {"Australia", "robot3", "yellow", ""},
				"SouthEastAsia":  {"SouthEastAsia", "robot4", "violet", ""},
				"SouthWestAsia":  {"SouthWestAsia", "robot5", "maroon", ""},
				"Antarctica":     {"Antarctica", "none", "white", ""},
			}
			FirstWorldBorders = Borders{
				"AlaskaRussia":               {"Alaska", "Russia", "battle"},
				"AlaskaCanada":               {"Alaska", "Canada", "human2"},
				"CanadaUSA":                  {"Canada", "USA", "battle"},
				"USACentralAmerica":          {"USA", "CentralAmerica", "human1"},
				"CentralAmericaSouthAmerica": {"CentralAmerica", "SouthAmerica", "battle"},
				"SouthAmericaBrazil":         {"SouthAmerica", "Brazil", "human3"},
				"WestAfricaEastAfrica":       {"WestAfrica", "EastAfrica", "human4"},
				"WestAfricaWestEurope":       {"WestAfrica", "WestEurope", "battle"},
				"WestEuropeNorthernEurope":   {"WestEurope", "NorthernEurope", "robot1"},
				"WestEuropeEastEurope":       {"WestEurope", "EastEurope", "battle"},
				"EastEuropeMiddleEast":       {"EastEurope", "MiddleEast", "robot2"},
				"EastAfricaMiddleEast":       {"EastAfrica", "MiddleEast", "battle"},
				"NorthernEuropeRussia":       {"NorthernEurope", "Russia", "battle"},
				"EastEuropeRussia":           {"EastEurope", "Russia", "battle"},
				"MiddleEastRussia":           {"MiddleEast", "Russia", "battle"},
				"MiddleEastSouthWestAsia":    {"MiddleEast", "SouthWestAsia", "battle"},
				"SouthWestAsiaRussia":        {"SouthWestAsia", "Russia", "battle"},
				"SouthWestAsiaNorthAsia":     {"SouthWestAsia", "NorthAsia", "battle"},
				"NorthAsiaRussia":            {"NorthAsia", "Russia", "human5"},
				"SouthWestAsiaSouthEastAsia": {"SouthWestAsia", "SouthEastAsia", "battle"},
				"NorthAsiaSouthEastAsia":     {"NorthAsia", "SouthEastAsia", "battle"},
				"SouthEastAsiaAustralia":     {"SouthEastAsia", "Australia", "battle"},
				"CanadaWestEurope":           {"Canada", "WestEurope", "battle"},
			}
			// fmt.Printf("First World: %v Origin: %v \n", FirstWorld["USA"].Color, OriginFirstWorld["USA"].Color)
			simulateText.SetText("")
			FirstWorld.RenderSVGs(simMapSVG)
		}
	})
	stopSimulationButton := gi.AddNewButton(simulationBrow, "stopSimulationButton")
	stopSimulationButton.Text = "Stop Simulation"
	stopSimulationButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			stopSimulation = true
		}
	})
	simulateText = gi.AddNewLabel(simulationTab, "simulateText", "                                                                            ")
	simulateText.SetProp("font-size", "20px")
	simulateText.Redrawable = true

	// width := 1024 // pixel sizes of screen
	// height := 768 // pixel sizes of screen

	simMapSVG = svg.AddNewSVG(simulationTab, "simMapSVG")
	simMapSVG.Fill = true
	simMapSVG.SetProp("background-color", "white")
	// simMapSVG.SetProp("width", units.NewPx(float32(width-20)))
	// simMapSVG.SetProp("height", units.NewPx(float32(height-100)))
	simMapSVG.SetStretchMaxWidth()
	simMapSVG.SetStretchMaxHeight()

	FirstWorld.RenderSVGs(simMapSVG)

	createBattleJoinLayouts()
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

	usernameKeyTitle := gi.AddNewLabel(playTab, "usernameKeyTitle", "<b>Battle first to 10 kills:</b>")
	usernameKeyTitle.SetProp("text-align", "center")
	usernameKeyTitle.SetProp("font-size", "40px")

	usernameKey := gi.AddNewFrame(playTab, "usernameKey", gi.LayoutVert)
	usernameKey.SetStretchMaxWidth()

	TheGame = &Game{} // Set up game
	TheGame.Config()  // Set up game

	tv.SelectTabByName("<b>Game</b>")
	tv.UpdateEnd(updt)
}
