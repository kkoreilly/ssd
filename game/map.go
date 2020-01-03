// Copyright (c) 2020, The EFight Authors. All rights reserved.
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
	"House1": {"House", mat32.Vec3{0, 0, -10}, DefScale},
	"House2": {"House", mat32.Vec3{0, 15, -10}, DefScale},
}