// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/goki/gi/gi"
	"github.com/goki/ki/ki"
	_ "github.com/lib/pq"
)

var db *sql.DB
var USER string
var PASSWORD string
var GOLD int
var LIVES int
var TEAM string
var goldNum int
var livesNum int

func data() {
	var str string
	var b []byte
	home, err := os.UserHomeDir()
	fn := filepath.Join(filepath.Join(home, "dburl"), "url.txt")
	b, err = ioutil.ReadFile(fn)
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
		if strings.Contains(teamName, "human") {
			keyItemText = gi.AddNewLabel(keyRow, "keyItemText", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
		} else if strings.Contains(teamName, "robot") {
			continue
		}
		keyItemText.SetProp("font-size", "20px")
		keyItemText.SetProp("background-color", color)
		clr := gi.Color{}
		clr.SetString(color, nil)
		if clr.IsDark() {
			keyItemText.SetProp("color", "white")
		} else {
			keyItemText.SetProp("color", "black")
		}
		keyItemText.Redrawable = true
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
		if strings.Contains(teamName, "robot") {
			keyItemText = gi.AddNewLabel(keyRow, "keyItemText", fmt.Sprintf("<b>%v:</b> %v", teamName, color))
		} else if strings.Contains(teamName, "human") {
			continue
		}
		keyItemText.SetProp("font-size", "20px")
		keyItemText.SetProp("background-color", color)
		clr := gi.Color{}
		clr.SetString(color, nil)
		if clr.IsDark() {
			keyItemText.SetProp("color", "white")
		} else {
			keyItemText.SetProp("color", "black")
		}
		keyItemText.Redrawable = true
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
		tr, has := FirstWorld[name]
		if !has {
			// fmt.Printf("Leaving loop")
			continue
		}
		tr.Owner = owner
		tr.Color = color
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
