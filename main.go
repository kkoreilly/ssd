// Copyright (c) 2020, The EFight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// var err error
	// tStr := os.Getenv("REPEAT")
	// repeat, err = strconv.Atoi(tStr)
	// if err != nil {
	// 	log.Print("Error converting $REPEAT to an int: %q - Using default", err)
	// 	repeat = 5
	// }
	//
	// db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatalf("Error opening database: %q", err)
	// }

	router := gin.New()
	router.Use(gin.Logger())

router.LoadHTMLGlob("webFiles/*.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
router.GET("/foodwaste", func(c *gin.Context) {
		c.HTML(http.StatusOK, "foodwaste.html", nil)
	})

	// router.GET("/mark", func(c *gin.Context) {
	// 	c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	// })
	//
	// router.GET("/repeat", repeatFunc)
	// router.GET("/db", dbFunc)

	router.Run(":" + port)
}
