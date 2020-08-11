// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"

	"github.com/goki/gi/svg"
)

type Territory struct {
	Name      string
	Owner     string
	Color     string
	SVGString string
}
type Border struct {
	Territory1 string
	Territory2 string
	Owner      string // if a team owns both territories, then it owns it. If it is a battle zone, value is "battle"
}

type World map[string]*Territory
type Borders map[string]*Border

var curCountSimulation int
var stopSimulation = false

var TeamStrength = map[string]float32{
	"human1": 1,
	"human2": 1,
	"human3": 1,
	"human4": 1,
	"human5": 1,
	"robot1": 1,
	"robot2": 1,
	"robot3": 1,
	"robot4": 1,
	"robot5": 1,
}

var FirstWorld = World{
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

var TheWorldMap = World{
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

var FirstWorldBorders = Borders{
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

func InitStrength() {
	// fmt.Printf("init strength \n")
	for k := range TeamStrength {
		randNum := rand.Float32()
		TeamStrength[k] = randNum
		simulationControlsTab.SetFullReRender()
		// fmt.Printf("Team Strength for team %v: %v \n", k, TeamStrength[k])
	}
}

func (wr *World) RenderSVGs(sv *svg.SVG) {
	updt := sv.UpdateStart()
	// fmt.Printf("Update Value: %v \n", updt)
	sv.DeleteChildren(true)
	sv.Norm = true
	sv.ViewBox.Size.Set(2754, 1398)
	sv.ViewBox.Min.X = -30
	sv.ViewBox.Min.Y = 61
	readWorld()
	// Create ocean:
	op := svg.AddNewPath(sv, "Ocean", FirstSVG["Ocean"].Data)
	op.SetProp("fill", "lightblue")

	for _, t := range *wr {
		p := svg.AddNewPath(sv, t.Name, FirstSVG[t.Name].Data)
		p.SetProp("fill", t.Color)
	}
	antText := svg.AddNewText(sv, "antText", 1377, 1390, "Neutral")
	antText.SetProp("font-size", "30px")
	sv.UpdateEnd(updt)
}

func (bd *Borders) simulateMap(fullSim bool) {
	// fmt.Printf("Human1 value: %v \n", TeamStrength["human1"])
	// updt := simulationTab.UpdateStart()
	// defer simulationTab.UpdateEnd(updt)
	for i := 0; 1 < 2; i++ {
		for _, b := range *bd { // do the battles
			if stopSimulation {
				return
			}
			if b.Owner == "battle" { // if there is a battle to be had, randomly decide the battle
				randNum := rand.Float32()
				// fmt.Printf("Team 1 strength: %v Total Strength: %v \n", TeamStrength[FirstWorld[b.Territory1].Owner], (TeamStrength[FirstWorld[b.Territory1].Owner] + TeamStrength[FirstWorld[b.Territory2].Owner]))
				probNum := TeamStrength[FirstWorld[b.Territory1].Owner] / (TeamStrength[FirstWorld[b.Territory1].Owner] + TeamStrength[FirstWorld[b.Territory2].Owner])
				// fmt.Printf("Rand Num: %v Prob Num: %v \n", randNum, probNum)
				// fmt.Printf("Random Number: %v \n", randNum)
				if randNum <= probNum { // team1 wins the battle
					winTeam := FirstWorld[b.Territory1].Owner // get the winning team
					FirstWorld[b.Territory2].Owner = winTeam  // set the losing team's territory to be owned by the winning team
					FirstWorld[b.Territory2].Color = FirstWorld[b.Territory1].Color
					FirstWorldBorders[b.Territory1+b.Territory2].Owner = winTeam
					TeamStrength[FirstWorld[b.Territory1].Owner] += 0.5
					TeamStrength[FirstWorld[b.Territory2].Owner] -= 0.5
					// fmt.Printf("(team1 type) Team %v wins and takes the territory %v \n", winTeam, FirstWorld[b.Territory2].Name)
				} else { // team2 wins the battle
					winTeam := FirstWorld[b.Territory2].Owner // get the winning team
					FirstWorld[b.Territory1].Owner = winTeam  // set the losing team's territory to be owned by the winning team
					FirstWorld[b.Territory1].Color = FirstWorld[b.Territory2].Color
					FirstWorldBorders[b.Territory1+b.Territory2].Owner = winTeam
					TeamStrength[FirstWorld[b.Territory2].Owner] += 0.5
					TeamStrength[FirstWorld[b.Territory1].Owner] -= 0.5
					// fmt.Printf("(team2 type) Team %v wins and takes the territory %v \n", winTeam, FirstWorld[b.Territory2].Name)
				}
				FirstWorld.RenderSVGs(simMapSVG)
			}

		}
		//After doing the battles, update the borders:
		for _, b := range *bd {
			if FirstWorld[b.Territory1].Owner != FirstWorld[b.Territory2].Owner {
				b.Owner = "battle"
			}
		}

		if comebacks {
			// fmt.Printf("Comebacks! \n")
			teams := make(map[string]int)
			for _, d := range FirstWorld {
				if d.Owner == "none" {
					continue
				}
				teams[d.Owner] = teams[d.Owner] + 1
			}
			for k := range teams {
				if teams[k] == 1 {
					TeamStrength[k] = TeamStrength[k] + 20
				}
			}
		}
		// Now we did the battles, check if one team has won
		var x = ""
		var y = false
		for _, t := range FirstWorld {
			if t.Name == "Antarctica" {
				// fmt.Printf("Antarctica \n")
				continue
			} else if x == "" {
				x = t.Owner
			} else if x == t.Owner {
				// then we continue
			} else if x != t.Owner {
				y = true
			}
		}

		if y {
			if fullSim {
				curCountSimulation += 1
				continue
			} else {
				curCountSimulation += 1
				return
			}
		} else {
			curCountSimulation += 1
			updt := simulationTab.UpdateStart()
			simulateText.SetText(fmt.Sprintf("Amount of weeks taken: %v", curCountSimulation))
			simulationTab.SetFullReRender()
			simulationTab.UpdateEnd(updt)
			return
		}
	}
}
