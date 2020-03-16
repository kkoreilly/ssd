// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	// "math/rand"
	"os"
	"path/filepath"
	"strings"
	// "time"

	"github.com/goki/gi/gi"
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
	_ "github.com/lib/pq"
)

var db *sql.DB
var USER string     // Global variable for your username
var PASSWORD string // Global variable for your password
var GOLD int        // Global variable for the amount of gold you have in game
var LIVES int       // Global variable for the amount of lives you have in game
var TEAM string     // Global variable for what team you're on
var POINTS int      // Global variable for the currrent amount of points you have in a battle
var goldNum int
var livesNum int
var gameOpen = true
var curBattleTerritory1, curBattleTerritory2 string

func data() {
	var str string
	home, _ := os.UserHomeDir()
	fn := filepath.Join(filepath.Join(home, "dburl"), "url.txt")
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		// fmt.Printf("%v \n", err)
		str = "example.com"
	} else {
		str = strings.TrimSpace(string(b)) // convert content to a 'string'
	}
	// fmt.Printf("Test String: %v \n", str)
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
	statement := fmt.Sprintf("INSERT INTO players(username, posX, posY, posZ, battleName, points) VALUES ('%v', '%v', '%v', '%v', '%v', 0)", USER, 0, 1, 0, battleName)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
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
			joinBattleButton.Text = "Create a 1v1 Battle in this Battlefield"
			joinBattleButton.SetProp("font-size", "30px")
			joinBattleButton.SetProp("horizontal-align", gi.AlignCenter)
			rec := ki.Node{}
			rec.InitName(&rec, "rec")
			joinBattleButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.ButtonClicked) {
					currentMapString = "Training Map 1"
					currentMap = FirstMap
					initPlayTab()
					joinPlayersTable(territory1 + territory2)
					curBattleTerritory1 = territory1
					curBattleTerritory2 = territory2
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
		gameResultText.SetText(fmt.Sprintf("<b>Congratulations on winning the battle with %v points. \nYour team (%v) wins one point in the battle %v vs. %v</b>", POINTS, TEAM, curBattleTerritory1, curBattleTerritory2))
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
			go removePlayer()
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
	statement := fmt.Sprintf("UPDATE players SET points = '%v' WHERE username = '%v'", value, username)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}
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
	fmt.Printf("%v \n", joinTeamStatement)

	result, err := db.Exec(joinTeamStatement)
	fmt.Printf("%v \n", result)

	if err != nil {
		fmt.Printf("Error")
		panic(err)
	}
	TEAM = name
	readTeam()
	teamMainText.SetText(teamMainText.Text + "\n\n<b>Click one of the buttons below to switch your team<b>.")

}

func (gm *Game) GetPosFromServer() { // GetPosFromServer loops through the players database and updates gm.OtherPos with the new data
	for {
		// fmt.Printf("Working 1 \n")
		getStatement := "SELECT * FROM players"
		rows, err := db.Query(getStatement)
		if err != nil {
			panic(err)
		}
		gm.PosMu.Lock()
		for rows.Next() {
			var username, battleName string
			var posX, posY, posZ float32
			var points int
			rows.Scan(&username, &posX, &posY, &posZ, &battleName, &points)
			// fmt.Printf("Username: %v \n", username)
			// fmt.Printf("User: %v \n", USER)
			if username != USER {
				gm.OtherPos[username] = &CurPosition{username, mat32.Vec3{posX, posY, posZ}, points}
				// fmt.Printf("Other Pos: %v \n", gm.OtherPos[username])
			}
		}
		// time.Sleep(100 * time.Millisecond)
		// gm.OtherPos["testyother"] = &CurPosition{"testyother", mat32.Vec3{rand.Float32()*5 - 2.5, 1, rand.Float32()*5 - 2.5}, 50}

		if !gm.GameOn {
			close(gm.PosUpdtChan)
			gm.battleOver(gm.Winner)
			gm.PosMu.Unlock()
			return
		}
		gm.PosMu.Unlock()
		gm.PosUpdtChan <- true // tells UpdatePeopleWorldPos to update to new positions!
		// todo: don't know from sender perspective if channel is still open!
		// if !ok {
		// 	return // game over
		// }
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
		goldResourcesText.SetText(fmt.Sprintf("%v \n \n You have %v gold", goldResourcesText.Text, goldNum))
		GOLD = goldNum
		livesResourcesText.SetText(fmt.Sprintf("%v \n \n You have %v lives", livesResourcesText.Text, livesNum))
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
}
func updatePosition(t string, value mat32.Vec3) {
	statement := fmt.Sprintf("UPDATE players SET posX = '%v' WHERE username='%v'", value.X, USER)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}

	statement2 := fmt.Sprintf("UPDATE players SET posZ = '%v' WHERE username='%v'", value.Z, USER)
	_, err = db.Exec(statement2)
	if err != nil {
		panic(err)
	}
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
