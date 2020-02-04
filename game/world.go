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
	//"EastUSA": {"EastUSA", "team1", "red", "m 75.43682,31.6938 c 6.626012,0 11.960311,0.17705 11.960311,0.39697 v 35.95563 c 0,0.21992 -5.334299,0.39696 -11.960311,0.39696 -6.626012,0 -11.960311,-0.17704 -11.960311,-0.39696 v -35.95563 c 0,-0.21992 5.334299,-0.39697 11.960311,-0.39697 z"},
	//"WestUSA": {"WestUSA", "team2", "blue", "m 51.516196,132.22835 c 6.626013,0 11.960311,0.17704 11.960311,0.39696 v 35.95563 c 0,0.21992 -5.334298,0.39697 -11.960311,0.39697 -6.626012,0 -11.960311,-0.17705 -11.960311,-0.39697 v -35.95563 c 0,-0.21992 5.334299,-0.39696 11.960311,-0.39696 z"},
	"Alaska": {"Alaska", "team2", "blue", ""},
	"USA":    {"USA", "team1", "red", ""},
	"Canada": {"Canada", "team3", "green", ""},
}

func (wr *World) RenderSVGs(sv *svg.SVG) {
	updt := sv.UpdateStart()
	sv.DeleteChildren(true)
	sv.Norm = true
	sv.ViewBox.Size.Set(2754, 1398)
	for _, t := range *wr {
		p := svg.AddNewPath(sv, t.Name, FirstSVG[t.Name].Data)
		p.SetProp("fill", t.Color)
	}
	sv.UpdateEnd(updt)
}
