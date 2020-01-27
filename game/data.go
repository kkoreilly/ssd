// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"database/sql"
	"fmt"
"io/ioutil"
"os"
	_ "github.com/lib/pq"
)

var db *sql.DB

func data() {
	var str string
	var b []byte
	home, err := os.UserHomeDir()
	fn := fmt.Sprintf("%v/dburl/url.txt", home)
	b, err = ioutil.ReadFile(fn)
	    if err != nil {
				// fmt.Printf("%v \n", err)
	        str = "example.com"
	    } else {
	    str = string(b) // convert content to a 'string'
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

func addUser(user string, password string) {
	tableCreateStatement := `CREATE TABLE IF NOT EXISTS users (
		username varchar,
		passwd varchar
		)`
	tableResult, err := db.Query(tableCreateStatement)
	fmt.Printf("Result: %v \n", tableResult)
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
