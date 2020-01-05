// Copyright (c) 2020, The EFight Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/emer/eve/eve"
	"github.com/goki/gi/gi3d"
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
	roofThick := float32(3)
	bedHeight := float32(0.6)
	doorWalllWidth := float32(16)

	house := eve.AddNewGroup(par, name)
	floor := eve.AddNewBox(house, "floor", mat32.Vec3{0, thick / 2, 0}, mat32.Vec3{width, thick, depth})
	floor.Color = "grey" // for debugging
	floor.Vis = "BrickHouse.Floor"
	ceiling := eve.AddNewBox(house, "ceiling", mat32.Vec3{0, float32(3.5) - thick/2, 0}, mat32.Vec3{width, thick, depth})
	ceiling.Color = "grey" // for debugging
	ceiling.Vis = "BrickHouse.Ceiling"
	bwall := eve.AddNewBox(house, "back-wall", mat32.Vec3{0, height / 2, -depth / 2}, mat32.Vec3{3, height, thick})
	bwall.Color = "blue"
	bwall.Vis = "BrickHouse.WinWall"
	intwall := eve.AddNewBox(house, "int-wall", mat32.Vec3{-6, height / 2, 0}, mat32.Vec3{width, height, thick})
	intwall.Color = "blue"
	intwall.Vis = "BrickHouse.BlankWall"
	lwall := eve.AddNewBox(house, "left-wall", mat32.Vec3{-width / 2, height / 2, 0}, mat32.Vec3{depth, height, thick})
	lwall.Initial.SetAxisRotation(0, 1, 0, -90)
	lwall.Color = "green"
	lwall.Vis = "BrickHouse.WinWall"
	rwall := eve.AddNewBox(house, "right-wall", mat32.Vec3{width / 2, height / 2, 0}, mat32.Vec3{depth, height, thick})
	rwall.Vis = "BrickHouse.WinWall"
	rwall.Initial.SetAxisRotation(0, 1, 0, -90)
	rwall.Color = "red"
	lfwall := eve.AddNewBox(house, "front-wall-left", mat32.Vec3{-doorWalllWidth / 4, height / 2, depth / 2}, mat32.Vec3{7, height, thick})
	lfwall.Color = "yellow"
	lfwall.Vis = "BrickHouse.DoorWall.Left"
	rfwall := eve.AddNewBox(house, "front-wall-right", mat32.Vec3{doorWalllWidth / 4, height / 2, depth / 2}, mat32.Vec3{7, height, thick})
  rfwall.Color = "yellow"
	rfwall.Vis = "BrickHouse.DoorWall.Right"
	tfwall := eve.AddNewBox(house, "front-wall-top", mat32.Vec3{0, 3.25, depth / 2}, mat32.Vec3{1, 0.5, thick})
  tfwall.Color = "yellow"
	tfwall.Vis = "BrickHouse.DoorWall.Top"

	//Interior Wall 1:
	liwall := eve.AddNewBox(house, "int-wall-left", mat32.Vec3{-doorWalllWidth / 4, height / 2, ( 3 * depth) / 8}, mat32.Vec3{3.25, height, thick})
	liwall.Color = "yellow"
	liwall.Vis = "BrickHouse.IntWall.Left"
	liwall.Initial.SetAxisRotation(0, 1, 0, -90)
	riwall := eve.AddNewBox(house, "int-wall-right", mat32.Vec3{-doorWalllWidth / 4, height / 2, 1.625}, mat32.Vec3{3.25, height, thick})
  riwall.Color = "yellow"
	riwall.Vis = "BrickHouse.IntWall.Right"
	riwall.Initial.SetAxisRotation(0, 1, 0, -90)
	tiwall := eve.AddNewBox(house, "int-wall-top", mat32.Vec3{-doorWalllWidth / 4, 3.25, depth / 4}, mat32.Vec3{1, 0.5, thick})
  tiwall.Color = "yellow"
	tiwall.Vis = "BrickHouse.IntWall.Top"
	tiwall.Initial.SetAxisRotation(0, 1, 0, -90)
	//Interior Wall 2:
	liwall1 := eve.AddNewBox(house, "int-wall-left-1", mat32.Vec3{-doorWalllWidth / 4, height / 2, -( 3 * depth) / 8}, mat32.Vec3{3.25, height, thick})
	liwall1.Color = "yellow"
	liwall1.Vis = "BrickHouse.IntWall.Left"
	liwall1.Initial.SetAxisRotation(0, 1, 0, -90)
	riwall1 := eve.AddNewBox(house, "int-wall-right-1", mat32.Vec3{-doorWalllWidth / 4, height / 2, -1.625}, mat32.Vec3{3.25, height, thick})
	riwall1.Color = "yellow"
	riwall1.Vis = "BrickHouse.IntWall.Right"
	riwall1.Initial.SetAxisRotation(0, 1, 0, -90)
	tiwall1 := eve.AddNewBox(house, "int-wall-top-1", mat32.Vec3{-doorWalllWidth / 4, 3.25, -depth / 4}, mat32.Vec3{1, 0.5, thick})
	tiwall1.Color = "yellow"
	tiwall1.Vis = "BrickHouse.IntWall.Top"
	tiwall1.Initial.SetAxisRotation(0, 1, 0, -90)
	// Roof Top is Here. Currently uses box for physcis, need to make it into pyramid. Todo: Fix this.
	roof := eve.AddNewBox(house, "roof", mat32.Vec3{0, float32(5) - thick/2, 0}, mat32.Vec3{width, roofThick, depth})
	roof.Color = "grey" // for debugging
	roof.Vis = "BrickHouse.Roof"
	bed1 := eve.AddNewBox(house, "bed1", mat32.Vec3{-6.5, bedHeight / 2, -6.75}, mat32.Vec3{2, bedHeight, 1.5})
	bed1.Color = "yellow"
	bed1.Vis = "BrickHouse.Bed"
	// bed2 := eve.AddNewBox(house, "bed2", mat32.Vec3{-6.5, bedHeight / 2, -2.75}, mat32.Vec3{2, bedHeight, 1.5})
	// bed2.Color = "yellow"
	// bed2.Vis = "BrickHouse.Bed"
	bed3 := eve.AddNewBox(house, "bed3", mat32.Vec3{-6.5, bedHeight / 2, 2.75}, mat32.Vec3{2, bedHeight, 1.5})
	bed3.Color = "yellow"
	bed3.Vis = "BrickHouse.Bed"
	// bed4 := eve.AddNewBox(house, "bed4", mat32.Vec3{-6.5, bedHeight / 2, 6.75}, mat32.Vec3{2, bedHeight, 1.5})
	// bed4.Color = "yellow"
	// bed4.Vis = "BrickHouse.Bed"
	return house
}

func (gm *Game) LibMakeBrickDoorWall() {
	height := float32(3.5)
	thick := float32(.1)
	sc := &gm.Scene.Scene
	nm := "BrickHouse.DoorWall"
	// left part
	lnm := nm + ".Left"
	lwg := sc.NewInLibrary(lnm)
	lwm := gi3d.AddNewBox(sc, lnm, 7, height, thick)
	lws := gi3d.AddNewSolid(sc, lwg, lnm, lwm.Name())
	lws.Mat.Texture = gi3d.TexName("brick.jpg")
	lws.Mat.Tiling.Repeat.Set(4, 2)
// right part
	rnm := nm + ".Right"
	rwg := sc.NewInLibrary(rnm)
	rwm := gi3d.AddNewBox(sc, rnm, 7, height, thick)
	rws := gi3d.AddNewSolid(sc, rwg, rnm, rwm.Name())
	rws.Mat.Texture = gi3d.TexName("brick.jpg")
	rws.Mat.Tiling.Repeat.Set(4, 2)
// top part
	tnm := nm + ".Top"
	twg := sc.NewInLibrary(tnm)
	twm := gi3d.AddNewBox(sc, tnm, 1, 0.5, thick)
	tws := gi3d.AddNewSolid(sc, twg, tnm, twm.Name())
	tws.Mat.Texture = gi3d.TexName("brick.jpg")
	tws.Mat.Tiling.Repeat.Set(4.0/7.0, 2.0/7.0)
}

func (gm *Game) LibMakeIntDoorWall() {
	height := float32(3.5)
	thick := float32(.1)
	sc := &gm.Scene.Scene
	nm := "BrickHouse.IntWall"
	// left part
	lnm := nm + ".Left"
	lwg := sc.NewInLibrary(lnm)
	lwm := gi3d.AddNewBox(sc, lnm, 3.25, height, thick)
	lws := gi3d.AddNewSolid(sc, lwg, lnm, lwm.Name())
	lws.Mat.Color.SetName("white")
// right part
	rnm := nm + ".Right"
	rwg := sc.NewInLibrary(rnm)
	rwm := gi3d.AddNewBox(sc, rnm, 3.25, height, thick)
	rws := gi3d.AddNewSolid(sc, rwg, rnm, rwm.Name())
	rws.Mat.Color.SetName("white")
// top part
	tnm := nm + ".Top"
	twg := sc.NewInLibrary(tnm)
	twm := gi3d.AddNewBox(sc, tnm, 1, 0.5, thick)
	tws := gi3d.AddNewSolid(sc, twg, tnm, twm.Name())
	tws.Mat.Color.SetName("white")
}


func (gm *Game) LibMakeBrickHouse() {
	sc := &gm.Scene.Scene
	// _, err := sc.OpenToLibrary("objs/BrickHouse.DoorWall.sobj", "BrickHouse.DoorWall")
	// if err != nil {
	// 	log.Println(err)
	// }

	gm.LibMakeBrickDoorWall()
	gm.LibMakeIntDoorWall()

	_, err := sc.OpenToLibrary("objs/BrickHouse.WinWall.obj", "BrickHouse.WinWall")
	if err != nil {
		log.Println(err)
	}

	bw, err := sc.OpenToLibrary("objs/BrickHouse.BlankWall.obj", "BrickHouse.BlankWall")
	if err != nil {
		log.Println(err)
	}
	bw.Pose.Scale.Set(0.2, 1, 1)

	_, err = sc.OpenToLibrary("objs/BrickHouse.Floor.obj", "BrickHouse.Floor")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("objs/BrickHouse.Ceiling.obj", "BrickHouse.Ceiling")
	if err != nil {
		log.Println(err)
	}

	_, err = sc.OpenToLibrary("objs/BrickHouse.Bed.obj", "BrickHouse.Bed")
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
