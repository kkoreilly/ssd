package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func data() {
	var err error

	db, err = sql.Open("postgres", "postgres://mesnfvhztxvwes:abd64ff99d5342f8f88f93248bba2d1c34821cc228dd847ca119287417516604@ec2-107-20-185-27.compute-1.amazonaws.com:5432/da7ng88d7bth78")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected!  %T \n", db)

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

	checkUsernameStatement := fmt.Sprintf("SELECT * FROM users	WHERE username='%v'", user)

	checkResultRows, err := db.Query(checkUsernameStatement)

	if err != nil {
		panic(err)
	}
	gotResults := false

	for checkResultRows.Next() {
		gotResults = true
	}
	fmt.Printf("Results (Got): %v \n", gotResults)

	if gotResults == false {
		fmt.Printf("Username isn't in use, will create user. \n")

		// create user code
		createAccountStatement :=

			fmt.Sprintf("INSERT INTO users(username, passwd) VALUES ('%v', '%v')", user, password)

		fmt.Printf("STATEMENT: %v \n", createAccountStatement)

		_, err := db.Exec(createAccountStatement)
		if err != nil {
			panic(err)
		}

		signUpResult.SetText("<b>Account created</b>")

	} else {
		fmt.Printf("Failed, username exists. \n")
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
		err = rows.Scan(&username, &password)
		fmt.Printf("In rows \n")
		fmt.Printf("New username: %v New password %v \n", username, password)
		newText := fmt.Sprintf("Username: %v, Password: %v            ", username, password)
		inspectText.SetText(fmt.Sprintf("%v %v", inspectText.Text, newText))
	}

	// fixStatement :=
	// `ALTER TABLE Users RENAME TO users`
	// db.Query(fixStatement)

}
