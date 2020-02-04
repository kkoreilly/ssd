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
	"Alaska": {"Alaska", "team2", "blue", ""},
	"USA":    {"USA", "team1", "red", ""},
	"Canada": {"Canada", "team3", "green", ""},
}

func (wr *World) RenderSVGs(sv *svg.SVG) {
	updt := sv.UpdateStart()
	sv.DeleteChildren(true)
	sv.Norm = true
	sv.ViewBox.Size.Set(2754, 1398)
	readWorld()
	for _, t := range *wr {
		p := svg.AddNewPath(sv, t.Name, FirstSVG[t.Name].Data)
		p.SetProp("fill", t.Color)
	}
	sv.UpdateEnd(updt)
}
