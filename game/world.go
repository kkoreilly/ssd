// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/goki/gi/svg"
)

type Territory struct {
	Name      string
	Owner     string
	Color     string
	SVGString string
}

type World map[string]*Territory

var FirstWorld = World{
	"Alaska":         {"Alaska", "team2", "red", ""},
	"USA":            {"USA", "team1", "blue", ""},
	"Canada":         {"Canada", "team3", "green", ""},
	"Brazil":         {"Brazil", "team4", "purple", ""},
	"SouthAmerica":   {"SouthAmerica", "team5", "orange", ""},
	"CentralAmerica": {"CentralAmerica", "team6", "yellow", ""},
	"WestAfrica":     {"WestAfrica", "team7", "pink", ""},
	"EastAfrica":     {"EastAfrica", "team8", "gold", ""},
	"Russia":         {"Russia", "team9", "maroon", ""},
	"ChinaRegion":    {"ChinaRegion", "team10", "violet", ""},
	"WestEurope":     {"WestEurope", "team1", "blue", ""},
	"NothernEurope":  {"NorthernEurope", "team2", "red", ""},
	"EastEurope":     {"EastEurope", "team3", "green", ""},
	"MiddleEast":     {"MiddleEast", "team4", "purple", ""},
	"Australia":      {"Australia", "team5", "orange", ""},
	"SouthEastAsia":  {"SouthEastAsia", "team6", "yellow", ""},
	"SouthWestAsia":  {"SouthWestAsia", "team7", "pink", ""},
	"Antarctica":     {"Antarctica", "none", "white", ""},
}

func (wr *World) RenderSVGs(sv *svg.SVG) {
	updt := sv.UpdateStart()
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
	antText.Width = 1000
	sv.UpdateEnd(updt)
}
