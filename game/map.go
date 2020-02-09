// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/goki/gi/mat32"

type MapObj struct {
	ObjType string
	Pos     mat32.Vec3
	Scale   mat32.Vec3
	// Color   string
}

type Map map[string]*MapObj

var DefScale = mat32.Vec3{1, 1, 1}
var FirstMap = Map{
	// "BigComplex1": {"BigComplex", mat32.Vec3{0, 0, -30}, DefScale},
	// "House1":    {"House", mat32.Vec3{10, 0, -40}, DefScale},
	// "House1": {"House", mat32.Vec3{0, 0, -10}, DefScale},
	// "House2": {"House", mat32.Vec3{20, 0, -10}, DefScale},
	"Block1":   {"Block", mat32.Vec3{-80, 0, -10}, DefScale},
	"TheWall1": {"TheWall", mat32.Vec3{0, 0, 0}, DefScale},
}

type MapInfo struct {
	Name    string
	MapData Map
}

type Maps map[string]*MapInfo

var AllMaps = Maps{
	"FirstMap": {"Training Map 1", FirstMap},
}
