package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func data() {

	db, err := sql.Open("postgres", "postgres://mesnfvhztxvwes:abd64ff99d5342f8f88f93248bba2d1c34821cc228dd847ca119287417516604@ec2-107-20-185-27.compute-1.amazonaws.com:5432/da7ng88d7bth78")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected! \n")

}
