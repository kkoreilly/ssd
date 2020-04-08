// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/goki/gi/gi"
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
	_ "github.com/lib/pq"
	"net/http"
	// "strconv"
	"strings"
	"time"
	// "io/ioutil"
)

var db *sql.DB
var USER string     // Global variable for your username
var PASSWORD string // Global variable for your password
var GOLD int        // Global variable for the amount of gold you have in game
var LIVES int       // Global variable for the amount of lives you have in game
var TEAM string     // Global variable for what team you're on
var POINTS int      // Global variable for the currrent amount of points you have in a battle
var WEAPON = "Basic"
var goldNum int
var livesNum int
var gameOpen = true
var curBattleTerritory1, curBattleTerritory2 string
var CURBATTLE string

func data() {
	var str string
	if URL_GLOBAL != "" {
		str = URL_GLOBAL
	} else {
		fmt.Printf("Unable to connect to database, URL missing \n")
		return
	}
	var err error
	db, err = sql.Open("postgres", str)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Connected!  %T \n", db)

}

func serverGetPlayerPos() {
	resp, err := http.Get("http://ssdserver.herokuapp.com/playerPosGet")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	// for i := 0; scanner.Scan(); i++ {
	fmt.Printf("Got data: %v \n", scanner.Text())
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func writePlayerPosToServer(pos mat32.Vec3, battleName string) {
	// fmt.Printf("Battle Name: %v \n", battleName)
	info := &CurPosition{USER, battleName, POINTS, pos, TheGame.KilledBy, TheGame.SpawnCount}
	// fmt.Printf("Info: %v \n", info)
	b, _ := json.Marshal(info)
	// b := []byte(fmt.Sprintf("username: %v, battleName: %v, posX: %v, posY: %v, posZ: %v, points: %v", USER, battleName, pos.X, pos.Y, pos.Z, POINTS))
	buff := bytes.NewBuffer(b)
	resp, err := http.Post("http://ssdserver.herokuapp.com/playerPosPost", "application/json", buff)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
func writeFireEventToServer(origin mat32.Vec3, dir mat32.Vec3, dmg int, battleName string) {
	// fmt.Printf("Battle Name: %v \n", battleName)
	info := &FireEventInfo{USER, origin, dir, dmg, battleName, time.Now()}
	b, _ := json.Marshal(info)
	// b := []byte(fmt.Sprintf("username: %v, battleName: %v, posX: %v, posY: %v, posZ: %v, points: %v", USER, battleName, pos.X, pos.Y, pos.Z, POINTS))
	buff := bytes.NewBuffer(b)
	resp, err := http.Post("http://ssdserver.herokuapp.com/fireEventsPost", "application/json", buff)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

//
// func writeEnemyPlayerPosToServer(username string, pos mat32.Vec3, battleName string, points int) {
// 	// fmt.Printf("Battle Name: %v \n", battleName)
// 	info := &CurPosition{username, battleName, points, pos, "", 0}
// 	fmt.Printf("Info: %v \n", info)
// 	b, _ := json.Marshal(info)
// 	// b := []byte(fmt.Sprintf("username: %v, battleName: %v, posX: %v, posY: %v, posZ: %v, points: %v", USER, battleName, pos.X, pos.Y, pos.Z, POINTS))
// 	buff := bytes.NewBuffer(b)
// 	resp, err := http.Post("http://ssdserver.herokuapp.com/playerPosPost", "application/json", buff)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// }

// func removeFireEventFromServer(key int, battleName string) {
// 	keyS := strconv.Itoa(key)
// 	buff := bytes.NewBuffer([]byte(""))
// 	resp, err := http.Post("http://ssdserver.herokuapp.com/fireEventsDelete?battleName="+battleName+"&key="+keyS, "application/json", buff)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// }
func readTeam() {
	findUserStatement := fmt.Sprintf("SELECT * FROM users WHERE username='%v'", USER)

	findUserResult, err := db.Query(findUserStatement)
	if err != nil {
		panic(err)
	}
	findUserResult.Scan(&USER, &PASSWORD, &goldNum, &livesNum, &TEAM)
	teamMainText.SetText(fmt.Sprintf("<b>Your team is<b>: <i><u>%v</u></i>", TEAM))

}
func addTeamUpdateButtons() {
	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	findTeamsStatement := "SELECT * FROM teams"
	findTeamsResult, err := db.Query(findTeamsStatement)
	if err != nil {
		panic(err)
	}
	for findTeamsResult.Next() {
		var teamName, color string
		var numOfPeople int
		findTeamsResult.Scan(&teamName, &numOfPeople, &color)
		var button *gi.Button
		if strings.Contains(teamName, "robot") {
			button = gi.AddNewButton(tbrowR, fmt.Sprintf("teamButton%v", teamName))
		} else if strings.Contains(teamName, "human") {
			button = gi.AddNewButton(tbrowH, fmt.Sprintf("teamButton%v", teamName))
		}
		button.Text = fmt.Sprintf("Join the team %v", teamName)
		button.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.ButtonClicked) {
				joinTeam(teamName)
			}
		})
	}
}
func (gm *Game) GetFireEvents() {
	for {
		// startTime := time.Now()
		if !gm.GameOn {
			return
		}

		resp, err := http.Get("http://ssdserver.herokuapp.com/fireEventsGet/?battleName=" + CURBATTLE + "&username=" + USER)
		if err != nil {
			panic(err)
		}
		if resp.Status == "422" {
			fmt.Printf("422, Battle maps nil or battle name nil")
			gm.PosMu.Unlock()
			gm.PosUpdtChan <- true // tells UpdatePeopleWorldPos to update to new positions
			// gm.FireUpdtChan <- true
			continue
		}
		defer resp.Body.Close()
		gm.FireEventMu.Lock()
		decoder := json.NewDecoder(resp.Body)
		newInfo := make([]*FireEventInfo, 0)
		decoder.Decode(&newInfo)
		gm.FireEvents = append(gm.FireEvents, newInfo...)
		gm.FireEventMu.Unlock()
		// fmt.Printf("Total time for GetFireEvents: %v \n", time.Since(startTime).Milliseconds())
	}
}
func addFireEventToDB(creator string, damage int, origin mat32.Vec3, dir mat32.Vec3) {
	statement := fmt.Sprintf("INSERT INTO fireEvents(creator, damage, originX, originY, originZ, dirX, dirY, dirZ) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')", creator, damage, origin.X, origin.Y, origin.Z, dir.X, dir.Y, dir.Z)
	_, err := db.Exec(statement)
	if err != nil {
		fmt.Printf("Err: %v \n", err)
	}
}
func initBorders() {
	for _, d := range FirstWorldBorders {
		territory1 := d.Territory1
		territory2 := d.Territory2
		team1 := FirstWorld[territory1].Owner
		team2 := FirstWorld[territory2].Owner
		// activeString := "false"
		// if d.Owner == "battle" {
		// 	activeString = "true"
		// }
		statement := fmt.Sprintf("INSERT INTO borders(territory1, territory2, team1, team2) VALUES ('%v', '%v', '%v', '%v')", territory1, territory2, team1, team2)
		_, err := db.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
}
func joinPlayersTable(battleName string) {
	writePlayerPosToServer(mat32.Vec3{0, 1, 0}, battleName)
	// fmt.Printf("Battle Name 1: %v \n", battleName)
	CURBATTLE = battleName
	// fmt.Printf("Curbattle 1: %v \n", CURBATTLE)
	// statement := fmt.Sprintf("INSERT INTO players(username, posX, posY, posZ, battleName, points) VALUES ('%v', '%v', '%v', '%v', '%v', 0)", USER, 0, 1, 0, battleName)
	POINTS = 0
	// // fmt.Printf("Points Data: %v", POINTS)
	// _, err := db.Exec(statement)
	// if err != nil {
	// 	panic(err)
	// }
}
func createBattleJoinLayouts() {

	// updt := homeTab.UpdateStart()
	// defer homeTab.UpdateEnd(updt)
	homeTab.SetFullReRender()
	statement := "SELECT * FROM borders"
	rows, err := db.Query(statement)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var territory1, territory2, team1, team2 string
		var team1points, team2points int
		rows.Scan(&territory1, &territory2, &team1, &team2, &team1points, &team2points)
		if FirstWorldLive[territory1].Owner != team1 {
			fixStatement := fmt.Sprintf("UPDATE borders SET team1 = '%v' WHERE territory1 = '%v' AND territory2 = '%v'", FirstWorldLive[territory1].Owner, territory1, territory2)
			_, err := db.Exec(fixStatement)
			if err != nil {
				panic(err)
			}
		}

		if FirstWorldLive[territory2].Owner != team2 {
			fixStatement := fmt.Sprintf("UPDATE borders SET team2 = '%v' WHERE territory1 = '%v' AND territory2 = '%v'", FirstWorldLive[territory2].Owner, territory1, territory2)
			_, err := db.Exec(fixStatement)
			if err != nil {
				panic(err)
			}
		}
	}
	teamJoinTitle := gi.AddNewLabel(homeTab, "teamJoinTitle", "<b>Battles that you can join:</b>")
	teamJoinTitle.SetProp("text-align", "center")
	teamJoinTitle.SetProp("font-size", "40px")
	joinLayoutG := gi.AddNewFrame(homeTab, "joinLayoutG", gi.LayoutVert)
	joinLayoutG.SetStretchMaxWidth()
	rows, err = db.Query(statement)
	for rows.Next() {
		var territory1, territory2, team1, team2 string
		var team1points, team2points int
		rows.Scan(&territory1, &territory2, &team1, &team2, &team1points, &team2points)
		// fmt.Printf("Team 1 points: %v   Team 2 points: %v \n", team1points, team2points)
		// fmt.Printf("TEAM Global var: %v \n", TEAM)
		if (FirstWorldLive[territory1].Owner != FirstWorldLive[territory2].Owner) && (team1 == TEAM || team2 == TEAM) {
			joinLayout := gi.AddNewFrame(joinLayoutG, "joinLayout", gi.LayoutVert)
			joinLayout.SetStretchMaxWidth()
			scoreText := gi.AddNewLabel(joinLayout, "scoreText", fmt.Sprintf("<b>%v             -                %v</b>", team1points, team2points))
			scoreText.SetProp("font-size", "35px")
			scoreText.SetProp("text-align", "center")
			teamsText := gi.AddNewLabel(joinLayout, "teamsText", "Team "+team1+"           vs.             Team "+team2)
			teamsText.SetProp("font-size", "30px")
			teamsText.SetProp("text-align", "center")
			territoriesText := gi.AddNewLabel(joinLayout, "territoriesText", territory1+"   vs.  "+territory2)
			territoriesText.SetProp("font-size", "25px")
			territoriesText.SetProp("text-align", "center")
			joinBattleButton := gi.AddNewButton(joinLayout, "joinBattleButton")
			joinBattleButton.Text = "Join Battle"
			joinBattleButton.SetProp("font-size", "30px")
			joinBattleButton.SetProp("horizontal-align", gi.AlignCenter)
			rec := ki.Node{}
			rec.InitName(&rec, "rec")
			joinBattleButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.ButtonClicked) {
					currentMapString = "Training Map 1"
					currentMap = FirstMap
					curBattleTerritory1 = territory1
					curBattleTerritory2 = territory2
					initPlayTab()
				}
			})
		}
	}
	teamNoJoinTitle := gi.AddNewLabel(homeTab, "teamNoJoinTitle", "<b>Other Battles:</b>")
	teamNoJoinTitle.SetProp("text-align", "center")
	teamNoJoinTitle.SetProp("font-size", "40px")
	rows, err = db.Query(statement)
	if err != nil {
		panic(err)
	}
	joinLayoutG1 := gi.AddNewFrame(homeTab, "joinLayoutG1", gi.LayoutVert)
	joinLayoutG1.SetStretchMaxWidth()
	for rows.Next() {
		var territory1, territory2, team1, team2 string
		var team1points, team2points int
		rows.Scan(&territory1, &territory2, &team1, &team2, &team1points, &team2points)
		if FirstWorldLive[territory1].Owner != FirstWorldLive[territory2].Owner && (team1 != TEAM && team2 != TEAM) {
			joinLayout := gi.AddNewFrame(joinLayoutG1, "joinLayout1", gi.LayoutVert)
			joinLayout.SetStretchMaxWidth()
			scoreText := gi.AddNewLabel(joinLayout, "scoreText", fmt.Sprintf("<b>%v             -                %v</b>", team1points, team2points))
			scoreText.SetProp("font-size", "35px")
			scoreText.SetProp("text-align", "center")
			teamsText := gi.AddNewLabel(joinLayout, "teamsText", "Team "+team1+"           vs.             Team "+team2)
			teamsText.SetProp("font-size", "30px")
			teamsText.SetProp("text-align", "center")
			territoriesText := gi.AddNewLabel(joinLayout, "territoriesText", territory1+"   vs.  "+territory2)
			territoriesText.SetProp("font-size", "25px")
			territoriesText.SetProp("text-align", "center")
		}
	}
}
func (gm *Game) setGameOver(winner string) {
	gm.WorldMu.Lock()
	gm.PosMu.Lock()
	gm.Scene.Win.OSWin.SetCursorEnabled(true, false)
	gm.Scene.TrackMouse = false
	gm.GameOn = false
	gm.Winner = winner
	gm.WorldMu.Unlock()
	gm.PosMu.Unlock()
}
func (gm *Game) battleOver(winner string) {
	gm.WorldMu.Lock()
	tabIndex, _ := tv.TabIndexByName("<b>Game</b>")
	tv.DeleteTabIndex(tabIndex, true)
	gameResultTab := tv.AddNewTab(gi.KiT_Frame, "<b>Game Result</b>").(*gi.Frame)

	gameResultTab.Lay = gi.LayoutVert
	gameResultTab.SetStretchMaxWidth()
	gameResultTab.SetStretchMaxHeight()

	gameResultText := gi.AddNewLabel(gameResultTab, "gameResultText", "")
	if winner == USER {
		gameResultText.SetText(fmt.Sprintf("<b>Congratulations on winning the battle with %v points. \nYour team (%v) wins one point in the battle %v vs. %v. \nYou win 10 gold.</b>", POINTS, TEAM, curBattleTerritory1, curBattleTerritory2))
		updateResource("gold", GOLD+10)
		readResources()
	} else {
		oppTeam := getEnemyTeamFromName(winner)
		gameResultText.SetText(fmt.Sprintf("<b>User %v won the battle with %v points. \nTheir team (%v) wins one point in the battle %v vs. %v</b>", winner, gm.OtherPos[winner].Points, oppTeam, curBattleTerritory1, curBattleTerritory2))
	}
	tabIndexResult, _ := tv.TabIndexByName("<b>Game Result</b>")
	gameResultText.SetProp("text-align", "center")
	gameResultText.SetProp("font-size", "60px")

	returnToHomeTab := gi.AddNewButton(gameResultTab, "returnToHomeTab")
	returnToHomeTab.Text = "Return to home"
	returnToHomeTab.SetProp("font-size", "40px")
	returnToHomeTab.SetProp("horizontal-align", "center")
	gameResultTab.SetFullReRender()
	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	returnToHomeTab.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			tv.DeleteTabIndex(tabIndexResult, true)
			tv.SelectTabIndex(0)
			go resetServer()
		}
	})

	updateBorderPoints(getEnemyTeamFromName(winner), 1, curBattleTerritory1, curBattleTerritory2)
	joinLayout := homeTab.ChildByName("joinLayoutG", 0)
	joinLayout1 := homeTab.ChildByName("joinLayoutG1", 0)
	joinLayout.Delete(true)
	joinLayout1.Delete(true)
	joinLayoutTitle := homeTab.ChildByName("teamJoinTitle", 0)
	joinLayoutNoTitle := homeTab.ChildByName("teamNoJoinTitle", 0)
	joinLayoutTitle.Delete(true)
	joinLayoutNoTitle.Delete(true)
	readWorld()
	go createBattleJoinLayouts()
	tv.SelectTabIndex(tabIndexResult)
	gm.WorldMu.Unlock()
}
func resetServer() {
	buff := bytes.NewBuffer([]byte(""))
	resp, err := http.Post("http://ssdserver.herokuapp.com/cleanUpBattle?battleName="+CURBATTLE, "text", buff)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
func updateBorderPoints(team string, changeNum int, territory1, territory2 string) {
	rowsB, err := db.Query(fmt.Sprintf("SELECT * FROM borders WHERE territory1 = '%v' AND territory2 = '%v'", territory1, territory2))
	if err != nil {
		panic(err)
	}
	var teamType string
	var curPoints int
	for rowsB.Next() {
		var team1, team2 string
		var team1points, team2points int
		rowsB.Scan(&territory1, &territory2, &team1, &team2, &team1points, &team2points)
		if team1 == team {
			teamType = "team1"
			curPoints = team1points
		} else {
			teamType = "team2"
			curPoints = team2points
		}
		// fmt.Printf("Team: %v  Team 1: %v   Team 2: %v    Team Type: %v \n", team, team1, team2, teamType)

	}
	statement := fmt.Sprintf("UPDATE borders SET %v = '%v' WHERE territory1 = '%v' AND territory2='%v'", teamType+"points", changeNum+curPoints, territory1, territory2)
	// fmt.Printf("Statement: %v \n", statement)
	_, err = db.Exec(statement)
	if err != nil {
		panic(err)
	}
	if changeNum+curPoints >= 10 { // then ten points have been reached and the border battle has been won
		var losingTerritory string // the territory that has been taken over
		if teamType == "team1" {
			losingTerritory = territory2
		} else {
			losingTerritory = territory1
		}
		// fmt.Printf("Losing territory: %v \n", losingTerritory)
		updateTStatement := fmt.Sprintf("UPDATE world SET owner = '%v' WHERE name = '%v'", team, losingTerritory)
		_, err = db.Exec(updateTStatement)
		if err != nil {
			panic(err)
		}
		rowsT, err := db.Query(fmt.Sprintf("SELECT * FROM teams WHERE name = '%v'", team))
		if err != nil {
			panic(err)
		}
		var color string
		for rowsT.Next() {
			var name string
			var numOfPeople int
			rowsT.Scan(&name, &numOfPeople, &color)
		}
		updateCStatement := fmt.Sprintf("UPDATE world SET color = '%v' WHERE name = '%v'", color, losingTerritory)
		_, err = db.Exec(updateCStatement)
		if err != nil {
			panic(err)
		}

		updatePStatement := fmt.Sprintf("UPDATE borders SET team1points = 0 WHERE territory1 = '%v' AND territory2='%v'", territory1, territory2)
		_, err = db.Exec(updatePStatement)
		if err != nil {
			panic(err)
		}
		updatePStatement1 := fmt.Sprintf("UPDATE borders SET team2points = 0 WHERE territory1 = '%v' AND territory2='%v'", territory1, territory2)
		_, err = db.Exec(updatePStatement1)
		if err != nil {
			panic(err)
		}

		FirstWorldLive.RenderSVGs(mapSVG)
	}
}
func getEnemyTeamFromName(username string) (team string) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE username = '%v'", username))
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var username string
		var password string
		var gold int
		var lives int
		rows.Scan(&username, &password, &gold, &lives, &team)
	}
	return team
}
func updateBattlePoints(username string, value int) {
	// statement := fmt.Sprintf("UPDATE players SET points = '%v' WHERE username = '%v'", value, username)
	// _, err := db.Exec(statement)
	// if err != nil {
	// 	panic(err)
	// }
}
func setActive() {
	for _, d := range FirstWorldBorders {
		activeString := "f"
		if d.Owner == "battle" {
			activeString = "t"
		}
		fmt.Printf("Active string: %v \n", activeString)

		statement := fmt.Sprintf("UPDATE borders SET active='t'")
		_, err := db.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
	rowsB, err := db.Query("SELECT * FROM borders")
	if err != nil {
		panic(err)
	}
	for rowsB.Next() {
		var t string
		var active string
		rowsB.Scan(t, t, t, t, t, t, &active)
		// fmt.Printf("Active: %v \n", active)
		// fmt.Printf("In rows \n")
		// fmt.Printf("\n \n USER: %v TEAM: %v \n \n", username, team)
		// fmt.Printf("<b>Username:</b> %v        <b>Password:</b> %v        <b>Gold:</b> %v        <b>Lives:</b> %v        <b>Team:</b> %v\n \n", username, password, gold, lives, team)
		fmt.Printf("Active result: %v \n", active)
	}
}
func addKeyItems() {
	// the ordering of doing this twice and the if statements will make the key be in the correct order
	// todo: make this code more efficient
	findTeamsStatement := "SELECT * FROM teams"
	findTeamsResult, err := db.Query(findTeamsStatement)
	if err != nil {
		panic(err)
	}
	for findTeamsResult.Next() {
		var teamName, color string
		var numOfPeople int
		findTeamsResult.Scan(&teamName, &numOfPeople, &color)
		var keyItemText *gi.Label
		var keyItemText1 *gi.Label
		var keyItemTextM *gi.Label
		if strings.Contains(teamName, "human") {
			keyItemText = gi.AddNewLabel(keyRow, "keyItemText", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
			keyItemText1 = gi.AddNewLabel(keyRow1, "keyItemText1", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
			keyItemTextM = gi.AddNewLabel(keyRowM, "keyItemTextM", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
		} else if strings.Contains(teamName, "robot") {
			continue
		}
		keyItemText.SetProp("font-size", "20px")
		keyItemText.SetProp("background-color", color)
		keyItemText1.SetProp("font-size", "20px")
		keyItemText1.SetProp("background-color", color)
		keyItemTextM.SetProp("font-size", "20px")
		keyItemTextM.SetProp("background-color", color)
		clr := gi.Color{}
		clr.SetString(color, nil)
		if clr.IsDark() || color == "red" || color == "blue" { // if dark, text is white
			keyItemText.SetProp("color", "white")
			keyItemText1.SetProp("color", "white")
			keyItemTextM.SetProp("color", "white")
		} else { // else, text is black
			keyItemText.SetProp("color", "black")
			keyItemText1.SetProp("color", "black")
			keyItemTextM.SetProp("color", "black")
		}
		keyItemText.Redrawable = true
		keyItemText1.Redrawable = true
		keyItemTextM.Redrawable = true
	}
	findTeamsStatement = "SELECT * FROM teams"
	findTeamsResult, err = db.Query(findTeamsStatement)
	if err != nil {
		panic(err)
	}
	for findTeamsResult.Next() {
		var teamName, color string
		var numOfPeople int
		findTeamsResult.Scan(&teamName, &numOfPeople, &color)
		var keyItemText *gi.Label
		var keyItemText1 *gi.Label
		var keyItemTextM *gi.Label
		if strings.Contains(teamName, "robot") {
			keyItemText = gi.AddNewLabel(keyRow, "keyItemText", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
			keyItemText1 = gi.AddNewLabel(keyRow1, "keyItemText1", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
			keyItemTextM = gi.AddNewLabel(keyRowM, "keyItemText1", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
		} else if strings.Contains(teamName, "human") {
			continue
		}
		keyItemText.SetProp("font-size", "20px")
		keyItemText.SetProp("background-color", color)
		keyItemText1.SetProp("font-size", "20px")
		keyItemText1.SetProp("background-color", color)
		keyItemTextM.SetProp("font-size", "20px")
		keyItemTextM.SetProp("background-color", color)
		clr := gi.Color{}
		clr.SetString(color, nil)
		if clr.IsDark() || color == "red" || color == "blue" { // if dark, text is white
			keyItemText.SetProp("color", "white")
			keyItemText1.SetProp("color", "white")
			keyItemTextM.SetProp("color", "white")
		}
		if !clr.IsDark() || color == "yellow" || color == "orange" && color != "red" && color != "blue" { // else, text is black
			keyItemText.SetProp("color", "black")
			keyItemText1.SetProp("color", "black")
			keyItemTextM.SetProp("color", "black")
		}
		keyItemText.Redrawable = true
		keyItemText1.Redrawable = true
		keyItemTextM.Redrawable = true
	}
}
func joinTeam(name string) {
	joinTeamStatement := fmt.Sprintf("UPDATE users SET %v = '%v' WHERE username='%v'", "team", name, USER)
	// fmt.Printf("%v \n", joinTeamStatement)

	_, err := db.Exec(joinTeamStatement)
	// fmt.Printf("%v \n", result)

	if err != nil {
		fmt.Printf("Error")
		panic(err)
	}
	TEAM = name
	readTeam()
	teamMainText.SetText(teamMainText.Text + "\n\n<b>Click one of the buttons below to switch your team<b>.")
	joinLayout := homeTab.ChildByName("joinLayoutG", 0)
	joinLayout1 := homeTab.ChildByName("joinLayoutG1", 0)
	joinLayout.Delete(true)
	joinLayout1.Delete(true)
	joinLayoutTitle := homeTab.ChildByName("teamJoinTitle", 0)
	joinLayoutNoTitle := homeTab.ChildByName("teamNoJoinTitle", 0)
	joinLayoutTitle.Delete(true)
	joinLayoutNoTitle.Delete(true)
	readWorld()
	createBattleJoinLayouts()

}
func removeBulletFromDB(originP, dirP mat32.Vec3) {
	readStatement := "SELECT * FROM fireEvents"
	rows, err := db.Query(readStatement)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var creator string
		var damage int
		var origin, dir mat32.Vec3
		rows.Scan(&creator, &damage, &origin.X, &origin.Y, &origin.Z, &dir.X, &dir.Y, &dir.Z)
		if origin.Round() == originP.Round() && dir.Round() == dirP.Round() {
			statement := fmt.Sprintf("DELETE FROM fireEvents WHERE originX='%v' AND originY='%v' AND originZ='%v' AND dirX = '%v' AND dirY = '%v' AND dirZ = '%v'", origin.X, origin.Y, origin.Z, dir.X, dir.Y, dir.Z)
			_, err := db.Exec(statement)
			if err != nil {
				panic(err)
			}
		}

	}

}
func (gm *Game) clearAllBullets() {
	// for k, _ := range gm.FireEvents {
	// 	// removeFireEventFromServer(k, CURBATTLE)
	// }
}
func (gm *Game) GetPosFromServer() { // GetPosFromServer loops through the players database and updates gm.OtherPos with the new data
	for {
		// startTime := time.Now()
		// fmt.Printf("GetPosFromServer Lock: %v Milliseconds\n", time.Since(startTime).Milliseconds())
		// startServerTime := time.Now()
		resp, err := http.Get("http://ssdserver.herokuapp.com/playerPosGet/?battleName=" + CURBATTLE)
		if err != nil {
			panic(err)
		}
		if resp.Status == "422" {
			fmt.Printf("422, Battle maps nil or battle name nil")
			gm.PosMu.Unlock()
			gm.PosUpdtChan <- true // tells UpdatePeopleWorldPos to update to new positions
			// gm.FireUpdtChan <- true
			continue
		}
		defer resp.Body.Close()
		// fmt.Printf("Time for GetPosFromServer server stuff: %v Milliseconds \n", time.Since(startServerTime).Milliseconds())
		// startDecodingTime := time.Now()
		tempOtherPos := make(map[string]*CurPosition)
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&tempOtherPos)
		// fmt.Printf("Time for GetPosFromServer Decoding: %v Milliseconds \n", time.Since(startDecodingTime).Milliseconds())
		// startTempTime := time.Now()
		gm.PosMu.Lock()
		for _, d := range tempOtherPos {
			if gm.OtherPos[d.Username] == nil {
				continue
			}
			if (d.KilledBy == USER) && ((d.SpawnCount - 1) == gm.OtherPos[d.Username].SpawnCount) {
				POINTS += 1
				resultText.SetText("<b>You killed " + d.Username + "! You get one point.</b>")
				resultText.SetFullReRender()
			}
		}
		gm.OtherPos = tempOtherPos
		// fmt.Printf("Time for GetPosFromServer check for kills: %v Milliseconds \n", time.Since(startTempTime).Milliseconds())
		// otherTime := time.Now()
		if !gm.GameOn {
			close(gm.PosUpdtChan)
			// close(gm.FireUpdtChan)
			gm.battleOver(gm.Winner)
			gm.PosMu.Unlock()
			return
		}
		// fmt.Printf("Time for GetPosFromServer check game on: %v Milliseconds \n", time.Since(otherTime).Milliseconds())
		gm.PosMu.Unlock()
		gm.PosUpdtChan <- true // tells UpdatePeopleWorldPos to update to new positions
		// fmt.Printf("Time for PosUpdtChan: %v Milliseconds \n", time.Since(otherTime).Milliseconds())
		// gm.FireUpdtChan <- true
		// fmt.Printf("Time for FireUpdtChan: %v Milliseconds \n", time.Since(otherTime).Milliseconds())
		// fmt.Printf("Time for GetPosFromServer other stuff: %v Milliseconds \n", time.Since(otherTime).Milliseconds())
		// since := time.Since(startTime)
		// fmt.Printf("Total time for GetPosFromServer: %v Milliseconds\n", since.Milliseconds())
	}
}

func readResources() {
	findUserStatement := fmt.Sprintf("SELECT * FROM users WHERE username='%v'", USER)
	findUserResult, err := db.Query(findUserStatement)

	if err != nil {
		panic(err)
	}

	for findUserResult.Next() {
		findUserResult.Scan(&USER, &PASSWORD, &goldNum, &livesNum, &TEAM)
		// fmt.Printf("Gold: %v \n", goldNum)
		// fmt.Printf("Lives: %v \n", livesNum)
		goldResourcesText.SetText(fmt.Sprintf("You have %v gold", goldNum))
		GOLD = goldNum
		// livesResourcesText.SetText(fmt.Sprintf("%v \n \n You have %v lives", livesResourcesText.Text, livesNum))
		LIVES = livesNum
	}
}
func updateResource(name string, value int) {
	updateResourceStatement := fmt.Sprintf("UPDATE users SET %v = '%v' WHERE username='%v'", name, value, USER)
	_, err := db.Exec(updateResourceStatement)
	if err != nil {
		panic(err)
	} else {
		// fmt.Printf("Updated resource")
	}

}
func removePlayer() {
	statement := fmt.Sprintf("DELETE FROM players WHERE username='%v'", USER)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
}

// fmt.Printf("Find User Result: %v \n", findUserResult)
func readWorld() {
	// fmt.Printf("In function \n")
	readStatement := `SELECT * FROM world`
	readResult, err := db.Query(readStatement)
	if err != nil {
		panic(err)
	}
	var name, owner, color string
	for readResult.Next() {
		readResult.Scan(&name, &owner, &color)
		// fmt.Printf("In loop. Name = %v \n", name)
		tr, has := FirstWorldLive[name]
		if !has {
			// fmt.Printf("Leaving loop")
			continue
		}
		tr.Owner = owner
		tr.Color = color
	}
	var previousMapObjOwner string
	mapDone := true
	var i = 0
	for _, d := range FirstWorldLive {
		// fmt.Printf("Map Done During: %v \n", mapDone)
		if d.Owner == previousMapObjOwner || i == 0 || d.Owner == "none" || previousMapObjOwner == "none" {
			previousMapObjOwner = d.Owner
		} else {
			// fmt.Printf("False, Owner: %v Prev owner: %v \n", d.Owner, previousMapObjOwner)
			mapDone = false
			break
		}
		i = i + 1
	}
	// Code currently doesn't actually reset map, fix later
	// fmt.Printf("Map Done: %v \n", mapDone)
	if mapDone == true { // one team has taken over the world
		winTeam := previousMapObjOwner // the team that has taken over the world
		resetWorld()
		resetBorders()
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var username string
			var password string
			var gold int
			var lives int
			var team string
			rows.Scan(&username, &password, &gold, &lives, &team)
			if team == winTeam { // this person won!
				updateGoldStatement := fmt.Sprintf("UPDATE users SET gold = '%v' WHERE username = '%v'", gold+1000, username) // give them a thousand gold
				_, err := db.Exec(updateGoldStatement)
				if err != nil {
					panic(err)
				}
				sendMessage("important", fmt.Sprintf("Your team (%v) has taken over the world! You have been awarded 1000 gold for being on the winning team! The map has been reset and a new game started.", winTeam), username)
			} else {
				sendMessage("important", fmt.Sprintf("Team %v has taken over the world. The map has been reset and a new game started.", winTeam), username)
			}

			readMessages()
			readWorld()

		}

	}

}
func resetBorders() {
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
	statement := "DELETE FROM borders"
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
	initBorders()
}
func resetWorld() {
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
	for _, d := range FirstWorld {
		statement := fmt.Sprintf("UPDATE world SET owner = '%v' WHERE name = '%v'", d.Owner, d.Name)
		_, err := db.Exec(statement)
		if err != nil {
			panic(err)
		}
		statement1 := fmt.Sprintf("UPDATE world SET color = '%v' WHERE name = '%v'", d.Color, d.Name)
		_, err = db.Exec(statement1)
		if err != nil {
			panic(err)
		}
	}
}
func sendMessage(messageType string, message string, username string) {
	statement := fmt.Sprintf("INSERT INTO messages(messageType, message, username) VALUES ('%v', '%v', '%v')", messageType, message, username)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
}
func readMessages() {
	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	statement := "SELECT * FROM messages"
	rows, err := db.Query(statement)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("messageFrame%v", i)
		if homeTab.ChildByName(name, 0) != nil {
			homeTab.ChildByName(name, 0).Delete(true)
		}
	}
	for rows.Next() {
		var messageType, message, username string
		rows.Scan(&messageType, &message, &username)
		if username == USER {
			if messageType == "important" {
				var name string
				for i := 0; 1 < 2; i++ {
					name = fmt.Sprintf("messageFrame%v", i)
					if homeTab.ChildByName(name, 0) == nil {
						break
					}
				}
				messageFrame := gi.AddNewFrame(homeTab, name, gi.LayoutVert)
				messageFrame.SetStretchMaxWidth()
				// messageFrame.Lay = gi.LayoutVert
				messageFrame.SetProp("background-color", "lightgreen")
				messageText := gi.AddNewLabel(messageFrame, "importantMessageText", "")
				messageText.Text = "<b>IMPORTANT MESSAGE:</b> \n<b>" + message + "</b>"
				messageText.SetProp("font-size", "50px")
				messageText.SetProp("text-align", "center")
				messageText.SetProp("white-space", gi.WhiteSpaceNormal)
				messageText.SetProp("max-width", -1)
				messageText.SetProp("width", "20em")
				messageButton := gi.AddNewButton(messageFrame, "messageButton")
				messageButton.Text = "OK"
				messageButton.SetProp("font-size", "30px")
				messageButton.SetProp("horizontal-align", "center")
				messageButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
					if sig == int64(gi.ButtonClicked) {
						removeMessage(message, username)
						messageFrame.Delete(true)
					}
				})
			}
		}
	}

}
func removeMessage(message string, username string) {
	statement := fmt.Sprintf("DELETE FROM messages WHERE message = '%v' AND username='%v'", message, username)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
}
func updatePosition(t string, value mat32.Vec3) {
	writePlayerPosToServer(mat32.Vec3{value.X, value.Y, value.Z}, CURBATTLE)
	// statement := fmt.Sprintf("UPDATE players SET posX = '%v' WHERE username='%v'", value.X, USER)
	// _, err := db.Exec(statement)
	// if err != nil {
	// 	fmt.Printf("DB err: %v \n", err)
	// }
	//
	// statement2 := fmt.Sprintf("UPDATE players SET posZ = '%v' WHERE username='%v'", value.Z, USER)
	// _, err = db.Exec(statement2)
	// if err != nil {
	// 	fmt.Printf("DB err: %v \n", err)
	// }
}
func addUser(user string, password string) {
	tableCreateStatement := `CREATE TABLE IF NOT EXISTS users (
		username varchar,
		passwd varchar
		)`
	_, err := db.Query(tableCreateStatement)
	// fmt.Printf("Result: %v \n", tableResult)
	if err != nil {
		panic(err)
	}

	checkUsernameStatement := fmt.Sprintf("SELECT * FROM users WHERE username='%v'", user)

	checkResultRows, err := db.Query(checkUsernameStatement)

	if err != nil {
		panic(err)
	}
	gotResults := false

	for checkResultRows.Next() {
		gotResults = true
	}
	// fmt.Printf("Results (Got): %v \n", gotResults)

	if gotResults == false {
		// fmt.Printf("Username isn't in use, will create user. \n")

		// create user code
		createAccountStatement :=

			fmt.Sprintf("INSERT INTO users(username, passwd) VALUES ('%v', '%v')", user, password)

		// fmt.Printf("STATEMENT: %v \n", createAccountStatement)

		_, err := db.Exec(createAccountStatement)
		if err != nil {
			panic(err)
		}

		signUpResult.SetText("<b>Account created</b>")

	} else {
		// fmt.Printf("Failed, username exists. \n")
		signUpResult.SetText("<b>Username exists, failed</b>")
	}

}
func initInspect() {

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var username string
		var password string
		rows.Scan(&username, &password)
		// fmt.Printf("In rows \n")
		// fmt.Printf("New username: %v New password %v \n", username, password)
		newText := fmt.Sprintf("Username: %v, Password: %v            ", username, password)
		inspectText.SetText(fmt.Sprintf("%v %v", inspectText.Text, newText))
	}

	// fixStatement :=
	// `ALTER TABLE Users RENAME TO users`
	// db.Query(fixStatement)

}

func logIn(user string, password string) {
	loginCheckStatement := fmt.Sprintf("SELECT * FROM users WHERE username='%v' AND passwd='%v'", user, password)
	results, err := db.Query(loginCheckStatement)
	var in = false
	if err != nil {
		panic(err)
	}
	for results.Next() {
		in = true
	}
	if in == true {
		// fmt.Printf("Found pair, logging in \n")
		updt := tv.UpdateStart()

		USER = user
		PASSWORD = password
		tv.Viewport.SetFullReRender()
		tv.DeleteTabIndex(0, true)
		tv.DeleteTabIndex(0, true)
		initMainTabs()

		tv.SelectTabIndex(0)
		tv.UpdateEnd(updt)

	} else {
		// fmt.Printf("Username and password do not match \n")
		logInResult.SetText("<b>Username and password do not match</b>")
	}

}
