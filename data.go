// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	// "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/goki/gi/gi"
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
	_ "github.com/lib/pq"
	// "time"
)

var db *sql.DB
var WEAPON = "Basic"
var gameOpen = true
var curBattleTerritory1, curBattleTerritory2 string
var CURBATTLE string
var curTeam1, curTeam2 string
var POINTS int

func InitDatabase() {
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
}

func InitDataMaps() {
	AllUserInfo = make(map[string]*UserInfo)
	AllBorders = make(map[string]*BorderInfo)
	ThisUserInfo = &UserInfo{"", "", "", 1}
	userRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	for userRows.Next() {
		var UserInfoTemp = &UserInfo{"", "", "", 1}
		var UsernameTemp string
		userRows.Scan(&UsernameTemp, &UserInfoTemp.Password, &UserInfoTemp.Team, &UserInfoTemp.Gold)
		UserInfoTemp.Username = UsernameTemp
		AllUserInfo[UsernameTemp] = UserInfoTemp
	}

	borderRows, err := db.Query("SELECT * FROM borders")
	if err != nil {
		panic(err)
	}
	for borderRows.Next() {
		BorderInfoTemp := &BorderInfo{"", "", "", "", 1, 1}
		borderRows.Scan(&BorderInfoTemp.Territory1, &BorderInfoTemp.Territory2, &BorderInfoTemp.Team1, &BorderInfoTemp.Team2, &BorderInfoTemp.Team1Points, &BorderInfoTemp.Team2Points)
		AllBorders[BorderInfoTemp.Territory1+BorderInfoTemp.Territory2] = BorderInfoTemp
	}
}

// The two functions above are well written, used and good

// Function below is good and used.
// func writePlayerPosToServer(pos mat32.Vec3, battleName string) {
// 	info := &Player{ThisUserInfo.Username, POINTS, pos, TheGame.KilledBy, TheGame.SpawnCount}
// 	b, _ := json.Marshal(info)
// 	buff := bytes.NewBuffer(b)
// 	resp, err := http.Post("http://ssdserver.herokuapp.com/playerPosPost", "application/json", buff)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// }
//
// // Function below is used and good
// func writeFireEventToServer(origin mat32.Vec3, dir mat32.Vec3, dmg int, battleName string) {
// 	info := &FireEventInfo{ThisUserInfo.Username, origin, dir, dmg, battleName, time.Now()}
// 	b, _ := json.Marshal(info)
// 	buff := bytes.NewBuffer(b)
// 	resp, err := http.Post("http://ssdserver.herokuapp.com/fireEventsPost", "application/json", buff)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// }

// Below is used and good
func readTeam() {
	findUserStatement := fmt.Sprintf("SELECT * FROM users WHERE username='%v'", ThisUserInfo.Username)
	findUserResult, err := db.Query(findUserStatement)
	if err != nil {
		panic(err)
	}
	findUserResult.Scan(&ThisUserInfo.Username, &ThisUserInfo.Password, &ThisUserInfo.Gold, &ThisUserInfo.Team)
	teamMainText.SetText(fmt.Sprintf("<b>Your team is<b>: <i><u>%v</u></i>", ThisUserInfo.Team))
}

//Little sketch but othwerside good. Make better later
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

// LAG CAUSER, sketch focus later.
func (gm *Game) GetFireEvents() {
	for {
		if !gm.GameOn {
			return
		}

		resp, err := http.Get("http://ssdserver.herokuapp.com/fireEventsGet/?battleName=" + CURBATTLE + "&username=" + ThisUserInfo.Username)
		if err != nil {
			panic(err)
		}
		if resp.Status == "422" {
			fmt.Printf("422, Battle maps nil or battle name nil")
			gm.PosUpdtChan <- true // tells UpdatePeopleWorldPos to update to new positions
			continue
		}
		defer resp.Body.Close()
		gm.FireEventMu.Lock()
		decoder := json.NewDecoder(resp.Body)
		newInfo := make([]*FireEventInfo, 0)
		decoder.Decode(&newInfo)
		gm.FireEvents = append(gm.FireEvents, newInfo...)
		gm.FireEventMu.Unlock()
	}
}

//Backup function, only used in the case of someone winning. Leave for now...
func initBorders() {
	for _, d := range FirstWorldBorders {
		territory1 := d.Territory1
		territory2 := d.Territory2
		team1 := FirstWorld[territory1].Owner
		team2 := FirstWorld[territory2].Owner
		statement := fmt.Sprintf("INSERT INTO borders(territory1, territory2, team1, team2) VALUES ('%v', '%v', '%v', '%v')", territory1, territory2, team1, team2)
		_, err := db.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
}

// Simple, good, used
func joinPlayersTable(battleName string) {
	// writePlayerPosToServer(mat32.Vec3{0, 1, 0}, battleName)
	CURBATTLE = battleName
	POINTS = 0
}

// The function below is being retired and will be completely re-written and updated with new updates:

// func createBattleJoinLayouts() {
// 	homeTab.SetFullReRender()
// 	statement := "SELECT * FROM borders"
// 	rows, err := db.Query(statement)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for rows.Next() {
// 		var territory1, territory2, team1, team2 string
// 		var team1points, team2points int
// 		rows.Scan(&territory1, &territory2, &team1, &team2, &team1points, &team2points)
// 		if TheWorldMap[territory1].Owner != team1 {
// 			fixStatement := fmt.Sprintf("UPDATE borders SET team1 = '%v' WHERE territory1 = '%v' AND territory2 = '%v'", TheWorldMap[territory1].Owner, territory1, territory2)
// 			_, err := db.Exec(fixStatement)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
//
// 		if TheWorldMap[territory2].Owner != team2 {
// 			fixStatement := fmt.Sprintf("UPDATE borders SET team2 = '%v' WHERE territory1 = '%v' AND territory2 = '%v'", TheWorldMap[territory2].Owner, territory1, territory2)
// 			_, err := db.Exec(fixStatement)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// 	teamJoinTitle := gi.AddNewLabel(homeTab, "teamJoinTitle", "<b>Battles that you can join:</b>")
// 	teamJoinTitle.SetProp("text-align", "center")
// 	teamJoinTitle.SetProp("font-size", "40px")
// 	joinLayoutG := gi.AddNewFrame(homeTab, "joinLayoutG", gi.LayoutVert)
// 	joinLayoutG.SetStretchMaxWidth()
// 	rows, err = db.Query(statement)
// 	for rows.Next() {
// 		var territory1, territory2, team1, team2 string
// 		var team1points, team2points int
// 		rows.Scan(&territory1, &territory2, &team1, &team2, &team1points, &team2points)
// 		if (TheWorldMap[territory1].Owner != TheWorldMap[territory2].Owner) && (team1 == ThisUserInfo.Team || team2 == ThisUserInfo.Team) {
// 			joinLayout := gi.AddNewFrame(joinLayoutG, "joinLayout", gi.LayoutVert)
// 			joinLayout.SetStretchMaxWidth()
// 			scoreText := gi.AddNewLabel(joinLayout, "scoreText", fmt.Sprintf("<b>%v             -                %v</b>", team1points, team2points))
// 			scoreText.SetProp("font-size", "35px")
// 			scoreText.SetProp("text-align", "center")
// 			teamsText := gi.AddNewLabel(joinLayout, "teamsText", "Team "+team1+"           vs.             Team "+team2)
// 			teamsText.SetProp("font-size", "30px")
// 			teamsText.SetProp("text-align", "center")
// 			territoriesText := gi.AddNewLabel(joinLayout, "territoriesText", territory1+"   vs.  "+territory2)
// 			territoriesText.SetProp("font-size", "25px")
// 			territoriesText.SetProp("text-align", "center")
// 			joinBattleButton := gi.AddNewButton(joinLayout, "joinBattleButton")
// 			joinBattleButton.Text = "Join Battle"
// 			joinBattleButton.SetProp("font-size", "30px")
// 			joinBattleButton.SetProp("horizontal-align", gi.AlignCenter)
// 			rec := ki.Node{}
// 			rec.InitName(&rec, "rec")
// 			joinBattleButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
// 				if sig == int64(gi.ButtonClicked) {
// 					currentMapString = "The Arena"
// 					currentMap = TheArenaMap
// 					curBattleTerritory1 = territory1
// 					curBattleTerritory2 = territory2
// 					curTeam1 = team1
// 					curTeam2 = team2
// 					initPlayTab()
// 				}
// 			})
// 		}
// 	}
// 	teamNoJoinTitle := gi.AddNewLabel(homeTab, "teamNoJoinTitle", "<b>Other Battles:</b>")
// 	teamNoJoinTitle.SetProp("text-align", "center")
// 	teamNoJoinTitle.SetProp("font-size", "40px")
// 	rows, err = db.Query(statement)
// 	if err != nil {
// 		panic(err)
// 	}
// 	joinLayoutG1 := gi.AddNewFrame(homeTab, "joinLayoutG1", gi.LayoutVert)
// 	joinLayoutG1.SetStretchMaxWidth()
// 	for rows.Next() {
// 		var territory1, territory2, team1, team2 string
// 		var team1points, team2points int
// 		rows.Scan(&territory1, &territory2, &team1, &team2, &team1points, &team2points)
// 		if TheWorldMap[territory1].Owner != TheWorldMap[territory2].Owner && (team1 != ThisUserInfo.Team && team2 != ThisUserInfo.Team) {
// 			joinLayout := gi.AddNewFrame(joinLayoutG1, "joinLayout1", gi.LayoutVert)
// 			joinLayout.SetStretchMaxWidth()
// 			scoreText := gi.AddNewLabel(joinLayout, "scoreText", fmt.Sprintf("<b>%v             -                %v</b>", team1points, team2points))
// 			scoreText.SetProp("font-size", "35px")
// 			scoreText.SetProp("text-align", "center")
// 			teamsText := gi.AddNewLabel(joinLayout, "teamsText", "Team "+team1+"           vs.             Team "+team2)
// 			teamsText.SetProp("font-size", "30px")
// 			teamsText.SetProp("text-align", "center")
// 			territoriesText := gi.AddNewLabel(joinLayout, "territoriesText", territory1+"   vs.  "+territory2)
// 			territoriesText.SetProp("font-size", "25px")
// 			territoriesText.SetProp("text-align", "center")
// 		}
// 	}
// }

// 2 functions below a little sketch
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
	if winner == ThisUserInfo.Username {
		gameResultText.SetText(fmt.Sprintf("<b>Congratulations on winning the battle with %v points. \nYour team (%v) wins one point in the battle %v vs. %v. \nYou win 10 gold.</b>", POINTS, ThisUserInfo.Team, curBattleTerritory1, curBattleTerritory2))
		updateResource("gold", ThisUserInfo.Gold+10)
		readResources()
	} else {
		oppTeam := getEnemyTeamFromName(winner)
		gameResultText.SetText(fmt.Sprintf("<b>User %v won the battle with %v points. \nTheir team (%v) wins one point in the battle %v vs. %v</b>", winner, gm.Players[winner].Points, oppTeam, curBattleTerritory1, curBattleTerritory2))
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
	// go createBattleJoinLayouts()
	tv.SelectTabIndex(tabIndexResult)
	gm.WorldMu.Unlock()
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

		TheWorldMap.RenderSVGs(mapSVG)
	}
}

// Function above is trash TOP PRI FIX
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

// Function above needs to be fixed with new maps

// Way too complex, fix
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

// Good and used
func joinTeam(name string) {
	joinTeamStatement := fmt.Sprintf("UPDATE users SET %v = '%v' WHERE username='%v'", "team", name, ThisUserInfo.Username)

	_, err := db.Exec(joinTeamStatement)

	if err != nil {
		fmt.Printf("Error")
		panic(err)
	}
	ThisUserInfo.Team = name
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
	// createBattleJoinLayouts()

}

// Needs improvement
func (gm *Game) GetPosFromServer() { // GetPosFromServer loops through the players database and updates gm.Players with the new data
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
			// gm.PosMu.Unlock()
			gm.PosUpdtChan <- true // tells UpdatePeopleWorldPos to update to new positions
			// gm.FireUpdtChan <- true
			continue
		}
		defer resp.Body.Close()
		// fmt.Printf("Time for GetPosFromServer server stuff: %v Milliseconds \n", time.Since(startServerTime).Milliseconds())
		// startDecodingTime := time.Now()
		tempPlayers := make(map[string]*Player)
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&tempPlayers)
		// fmt.Printf("Time for GetPosFromServer Decoding: %v Milliseconds \n", time.Since(startDecodingTime).Milliseconds())
		// startTempTime := time.Now()
		gm.PosMu.Lock()
		for _, d := range tempPlayers {
			if gm.Players[d.Username] == nil {
				continue
			}
			if (d.KilledBy == ThisUserInfo.Username) && ((d.SpawnCount - 1) == gm.Players[d.Username].SpawnCount) {
				POINTS += 1
				resultText.SetText("<b>You killed " + d.Username + "! You get one point.</b>")
				resultText.SetFullReRender()
			}
		}
		gm.Players = tempPlayers
		// fmt.Printf("Time for GetPosFromServer check for kills: %v Milliseconds \n", time.Since(startTempTime).Milliseconds())
		// otherTime := time.Now()
		if !gm.GameOn {
			close(gm.PosUpdtChan)
			// close(gm.FireUpdtChan)
			gm.PosMu.Unlock()
			gm.battleOver(gm.Winner)
			return
		}
		gm.PosMu.Unlock()
		gm.PosUpdtChan <- true // tells UpdatePeopleWorldPos to update to new positions
	}
}

// Shouldn't be nessecary, review
func readResources() {
	findUserStatement := fmt.Sprintf("SELECT * FROM users WHERE username='%v'", ThisUserInfo.Username)
	findUserResult, err := db.Query(findUserStatement)

	if err != nil {
		panic(err)
	}

	for findUserResult.Next() {
		findUserResult.Scan(&ThisUserInfo.Username, &ThisUserInfo.Password, &ThisUserInfo.Gold, &ThisUserInfo.Team)
		// fmt.Printf("Gold: %v \n", ThisUserInfo.Gold)
		// fmt.Printf("Lives: %v \n", livesNum)
		goldResourcesText.SetText(fmt.Sprintf("You have %v gold", ThisUserInfo.Gold))
		// livesResourcesText.SetText(fmt.Sprintf("%v \n \n You have %v lives", livesResourcesText.Text, livesNum))
	}
}

// Good, TODO: Update map while doing this
func updateResource(name string, value int) {
	updateResourceStatement := fmt.Sprintf("UPDATE users SET %v = '%v' WHERE username='%v'", name, value, ThisUserInfo.Username)
	_, err := db.Exec(updateResourceStatement)
	if err != nil {
		panic(err)
	} else {
	}

}

// Needs to be fixed
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
		tr, has := TheWorldMap[name]
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
	for _, d := range TheWorldMap {
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

//Too complex of functions happening below
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

// Good
func sendMessage(messageType string, message string, username string) {
	statement := fmt.Sprintf("INSERT INTO messages(messageType, message, username) VALUES ('%v', '%v', '%v')", messageType, message, username)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
}

// Make less complex
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
		if username == ThisUserInfo.Username {
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

// Good
func removeMessage(message string, username string) {
	statement := fmt.Sprintf("DELETE FROM messages WHERE message = '%v' AND username='%v'", message, username)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
}

// Still used?
func updatePosition(t string, value mat32.Vec3) {
	// writePlayerPosToServer(mat32.Vec3{value.X, value.Y, value.Z}, CURBATTLE)
}

func addUser(user string, password string) {
	gotResults := false
	for k := range AllUserInfo {
		if user == k {
			gotResults = true
		}
	}
	if gotResults == false { // No user with our name exists, so create

		// create user code
		createAccountStatement := fmt.Sprintf("INSERT INTO users(username, passwd) VALUES ('%v', '%v')", user, password)
		AllUserInfo[user] = &UserInfo{user, password, "", 500}
		ThisUserInfo = &UserInfo{user, password, "", 500}
		_, err := db.Exec(createAccountStatement)
		if err != nil {
			panic(err)
		}

		signUpResult.SetText("<b>Account created</b>")

	} else {
		signUpResult.SetText("<b>Username exists, failed</b>")
	}

}

func logIn(user string, password string) {
	var in = false
	for k, d := range AllUserInfo {
		if k == user && d.Password == password {
			in = true

			ThisUserInfo.Username = user
			ThisUserInfo.Password = password
			ThisUserInfo.Gold = d.Gold
			ThisUserInfo.Team = d.Team
		}
	}
	if in == true {
		updt := tv.UpdateStart()
		tv.Viewport.SetFullReRender()
		tv.DeleteTabIndex(0, true)
		tv.DeleteTabIndex(0, true)
		initMainTabs()

		tv.SelectTabIndex(0)
		tv.UpdateEnd(updt)

	} else {
		logInResult.SetText("<b>Username and password do not match</b>")
	}

}
