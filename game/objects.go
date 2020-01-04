// Copyright (c) 2020, The EFight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/emer/eve/eve"
	"github.com/goki/gi/mat32"
	"github.com/goki/gi/gi3d"
)

// all units are in meters

// PhysMakeBrickHouse makes the EVE physics version of house
// in par = parent group.
func (gm *Game) PhysMakeBrickHouse(par *eve.Group, name string) *eve.Group {
	width := float32(15)
	depth := float32(15)
	height := float32(3.5)
	thick := float32(0.1) // wall, ceiling, floor
	roofThick := float32(3)

	house := eve.AddNewGroup(par, name)
	floor := eve.AddNewBox(house, "floor", mat32.Vec3{0, thick / 2, 0}, mat32.Vec3{width, thick, depth})
	floor.Color = "grey" // for debugging
	floor.Vis = "BrickHouse.Floor"
	ceiling := eve.AddNewBox(house, "ceiling", mat32.Vec3{0, float32(3.5) - thick / 2, 0}, mat32.Vec3{width, thick, depth})
	ceiling.Color = "grey" // for debugging
	ceiling.Vis = "BrickHouse.Ceiling"
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
	// Roof Top is Here. Currently uses box for physcis, need to make it into pyramid. Todo: Fix this.
	roof := eve.AddNewBox(house, "roof", mat32.Vec3{0, float32(5) - thick / 2, 0}, mat32.Vec3{width, roofThick, depth})
	roof.Color = "grey" // for debugging
	roof.Vis = "BrickHouse.Roof"
	return house
}

func (gm *Game) LibMakeBrickHouse() {
	sc := &gm.Scene.Scene
	_, err := sc.OpenToLibrary("objs/BrickHouse.DoorWall.obj", "BrickHouse.DoorWall")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("objs/BrickHouse.WinWall.obj", "BrickHouse.WinWall")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("objs/BrickHouse.Floor.obj", "BrickHouse.Floor")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("objs/BrickHouse.Ceiling.obj", "BrickHouse.Ceiling")
	if err != nil {
		log.Println(err)
	}

	rt, err := sc.OpenToLibrary("objs/BrickHouse.Roof.obj", "BrickHouse.Roof")
	if err != nil {
		log.Println(err)
		// } else {
		// 	rt.Pose.Pos.Set(-0.3725, 3.5, 0.3725)
		// 	rt.Pose.Scale.Set(1.05, 1, 1.05)
		// 	solidrt := rt.Child(0).Child(0).(*gi3da.Solid)
		// 	solidrt.Mat.CullBack = false
	}
	solidrt := rt.Child(0).Child(0).(*gi3d.Solid)
	solidrt.Mat.CullBack = false
	rt.Pose.Scale.Set(1.05, 1, 1.05)

	// bb, err := sc.OpenNewObj("bed1.obj", ogp)
	// if err != nil {
	// 	log.Println(err)
	//
	// } else {
	// 	bb.Pose.Pos.Set(0, 0, -13.5)
	// }

	// fi.Pose.SetAxisRotation(1, 0, 0, -45)
}
