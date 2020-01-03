// Copyright (c) 2020, The EFight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/emer/eve/eve"
	"github.com/goki/gi/mat32"
)

// all units are in meters

// PhysMakeBrickHouse makes the EVE physics version of house
// in par = parent group.
func (gm *Game) PhysMakeBrickHouse(par *eve.Group, name string) *eve.Group {
	width := float32(15)
	depth := float32(15)
	height := float32(3.5)
	thick := float32(0.1) // wall, ceiling, floor

	house := eve.AddNewGroup(par, name)
	floor := eve.AddNewBox(house, "floor", mat32.Vec3{0, thick / 2, 0}, mat32.Vec3{width, thick, depth})
	floor.Color = "grey" // for debugging
	floor.Vis = "BrickHouse.Floor"
	bwall := eve.AddNewBox(house, "back-wall", mat32.Vec3{0, height / 2, -depth / 2}, mat32.Vec3{width, height, thick})
	bwall.Color = "blue"
	bwall.Vis = "BrickHouse.WinWall"
	lwall := eve.AddNewBox(house, "left-wall", mat32.Vec3{-width / 2, height / 2, 0}, mat32.Vec3{depth, height, thick})
	lwall.Initial.SetAxisRotation(0, 1, 0, -90)
	lwall.Color = "green"
	lwall.Vis = "BrickHouse.WinWall"
	rwall := eve.AddNewBox(house, "right-wall", mat32.Vec3{width / 2, height / 2, 0}, mat32.Vec3{depth, height, thick})
	rwall.Vis = "BrickHouse.WinWall"
	rwall.Initial.SetAxisRotation(0, 1, 0, -90)
	rwall.Color = "red"
	fwall := eve.AddNewBox(house, "front-wall", mat32.Vec3{0, height / 2, depth / 2}, mat32.Vec3{width, height, thick})
	fwall.Color = "yellow"
	fwall.Vis = "BrickHouse.DoorWall"
	return house
}

func (gm *Game) LibMakeBrickHouse() {
	sc := &gm.Scene.Scene
	_, err := sc.OpenToLibrary("doorWall1.obj", "BrickHouse.DoorWall")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("windowWall1.obj", "BrickHouse.WinWall")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("floor1.obj", "BrickHouse.Floor")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("roof1.obj", "BrickHouse.Ceiling")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("roofTop1.obj", "BrickHouse.Roof")
	if err != nil {
		log.Println(err)
		// } else {
		// 	rt.Pose.Pos.Set(-0.3725, 3.5, 0.3725)
		// 	rt.Pose.Scale.Set(1.05, 1, 1.05)
		// 	solidrt := rt.Child(0).Child(0).(*gi3d.Solid)
		// 	solidrt.Mat.CullBack = false
	}

	// bb, err := sc.OpenNewObj("bed1.obj", ogp)
	// if err != nil {
	// 	log.Println(err)
	//
	// } else {
	// 	bb.Pose.Pos.Set(0, 0, -13.5)
	// }

	// fi.Pose.SetAxisRotation(1, 0, 0, -45)
}
