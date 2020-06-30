// Copyright (c) 2020, The Singularity Showdown Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sort"
	"sync"

	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/emer/eve/eve"
	"github.com/emer/eve/evev"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gi3d"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/oswin"
	"github.com/goki/gi/oswin/key"
	"github.com/goki/gi/oswin/mouse"
	"github.com/goki/gi/units"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

const (
	BulletTimeMsec = 300
)

type CurPosition struct {
	Username   string
	BattleName string
	Points     int
	Pos        mat32.Vec3
	KilledBy   string
	SpawnCount int
}

type Weapon struct {
	Name     string
	MinD     float32 `desc:"Minimum amount of damage that this weapon will do"`
	MaxD     float32 `desc:"Maximum amount of damage that this weapon will do"`
	FireRate float32 `desc:"How many times this weapon can be fired in a second"`
}

type FireEventInfo struct {
	Creator    string
	Origin     mat32.Vec3
	Dir        mat32.Vec3
	Damage     int
	BattleName string
	StartTime  time.Time
}

type Game struct {
	World          *eve.Group
	View           *evev.View
	Scene          *Scene
	Map            Map
	MapObjs        map[string]bool
	OtherPos       map[string]*CurPosition
	PosUpdtChan    chan bool  `desc:"channel connecting server pos updates with world state update"`
	PosMu          sync.Mutex `desc:"protects updates to OtherPos map"`
	WorldMu        sync.Mutex `desc:"protects updates to World physics and view"`
	GameOn         bool       // starts on when game turned out, turn off when close game
	Winner         string
	PersHitWall    bool
	Gravity        float32
	AbleToFire     bool
	FireEvents     []*FireEventInfo
	FireEventMu    sync.Mutex
	FireUpdtChan   chan bool
	KilledBy       string
	SpawnCount     int
	SpawnPositions []mat32.Vec3
	IsDead         bool
}

// TheGame is the game instance for the current game
var TheGame *Game
var HEALTH float32 = 100 // how much health you have
type Weapons map[string]*Weapon

var TheWeapons = Weapons{
	"Basic":  {"Basic", 10, 30, 2},
	"Sniper": {"Sniper", 60, 100, 0.25},
}

type Scene struct {
	gi3d.Scene
	TrackMouse bool
	CamRotUD   float32 // current camera rotation up / down
	CamRotLR   float32 // current camera rotation LR
}

var KiT_Scene = kit.Types.AddType(&Scene{}, nil)

func (gm *Game) BuildMap() {
	ogp := eve.AddNewGroup(gm.World, "FirstPerson")
	gm.PhysMakePerson(ogp, "FirstPerson", true)
	posX, posZ := gm.GetRandomSpawnPoint()
	ogp.Initial.Pos.Set(posX, 1, posZ)

	eve.AddNewGroup(gm.World, "PeopleGroup")
	gi3d.AddNewGroup(&gm.Scene.Scene, &gm.Scene.Scene, "RayGroup")
	gi3d.AddNewGroup(&gm.Scene.Scene, &gm.Scene.Scene, "PeopleTextGroup")
	for nm, obj := range gm.Map {
		// fmt.Printf("Object type: %v \n", obj.ObjType)
		gm.MakeObj(obj, nm)
	}
}

func (gm *Game) MakeObj(obj *MapObj, nm string) *eve.Group {
	var ogp *eve.Group
	switch obj.ObjType {
	case "TheWall":
		ogp = eve.AddNewGroup(gm.World, nm)
		gm.PhysMakeTheWall(ogp, nm)
	case "House":
		ogp = eve.AddNewGroup(gm.World, nm)
		gm.PhysMakeBrickHouse(ogp, nm)
	case "Block":
		ogp = eve.AddNewGroup(gm.World, nm)
		for i := 0; i < 3; i++ {
			house := gm.PhysMakeBrickHouse(ogp, fmt.Sprintf("%v-House%v", nm, i))
			house.Initial.Pos.Set(float32(20*i), 0, 0)
		}
	case "Road":
		ogp = eve.AddNewGroup(gm.World, nm)
		gm.PhysMakeRoad(ogp, nm)
	case "ArenaWalls":
		ogp = eve.AddNewGroup(gm.World, nm)
		gm.PhysMakeArenaWalls(ogp, nm)
	case "LavaBlockFour":
		ogp = eve.AddNewGroup(gm.World, nm)
		gm.PhysMakeLava(ogp, nm, 2, 2)
	case "LavaWallsSetup":
		ogp = eve.AddNewGroup(gm.World, nm)
		backRight := gm.PhysMakeLava(ogp, nm+"backRight", 8, 1)
		backRight.Initial.Pos.Set(10, 0, 85)
		backLeft := gm.PhysMakeLava(ogp, nm+"backLeft", 8, 1)
		backLeft.Initial.Pos.Set(-90, 0, 85)

		frontRight := gm.PhysMakeLava(ogp, nm+"frontRight", 8, 1)
		frontRight.Initial.Pos.Set(10, 0, -85)
		frontLeft := gm.PhysMakeLava(ogp, nm+"frontLeft", 8, 1)
		frontLeft.Initial.Pos.Set(-90, 0, -85)

		rightRight := gm.PhysMakeLava(ogp, nm+"rightRight", 1, 7)
		rightRight.Initial.Pos.Set(90, 0, 10)
		rightLeft := gm.PhysMakeLava(ogp, nm+"rightLeft", 1, 7)
		rightLeft.Initial.Pos.Set(90, 0, -10)

		leftRight := gm.PhysMakeLava(ogp, nm+"leftRight", 1, 7)
		leftRight.Initial.Pos.Set(-90, 0, 10)
		leftLeft := gm.PhysMakeLava(ogp, nm+"leftLeft", 1, 7)
		leftLeft.Initial.Pos.Set(-90, 0, -10)
	}
	/*
		case "Hill":
				ogp = gi3d.AddNewGroup(sc, sc, nm)
				o := gi3d.AddNewObject(sc, ogp, "hill", "Hill")
				o.Pose.Pos.Set(0, 0, 0)
				// o.Pose.Scale.Set(1, 10, 1)
				o.Mat.Color.SetString("green", nil)
			case "House":
				ogp = gi3d.AddNewGroup(sc, sc, nm)
				o := gi3d.AddNewObject(sc, ogp, "house_ground", "HouseFloor")
				o.Pose.Pos.Set(0, 0, 0)
				// o.Pose.Scale.Set(10, 0.01, 10)
				o.Mat.Color.SetString("brown", nil)

				o = gi3d.AddNewObject(sc, ogp, "house_roof", "HouseRoof")
				o.Pose.Pos.Set(0, 4.995, 0)
				// o.Pose.Scale.Set(10, 0.01, 10)
				o.Mat.Color.SetString("brown", nil)

				o = gi3d.AddNewObject(sc, ogp, "house_wall1", "HouseWallOne")
				o.Pose.Pos.Set(-5.25, 0, 0)
				// o.Pose.Scale.Set(0.5, 10, 10)
				o.Mat.Color.SetString("brown", nil)

				o = gi3d.AddNewObject(sc, ogp, "house_wall2", "HouseWallOne")
				o.Pose.Pos.Set(4.75, 0, 0)
				// o.Pose.Scale.Set(0.5, 10, 10)
				o.Mat.Color.SetString("brown", nil)

				o = gi3d.AddNewObject(sc, ogp, "house_wall3", "HouseWallTwo")
				o.Pose.Pos.Set(0, 0, -5)
				// o.Pose.Scale.Set(10, 10, 0.5)
				o.Mat.Color.SetString("brown", nil)

				o = gi3d.AddNewObject(sc, ogp, "house_bed1", "HouseBedOne")
				o.Pose.Pos.Set(-3.5, 0, -4)
				// o.Mat.Color.SetString("green", nil)
				o.Mat.SetTextureName(sc, "HouseBed")

				o = gi3d.AddNewObject(sc, ogp, "house_blanket1", "HouseBlanketOne")
				o.Pose.Pos.Set(-3.5, 1.05, -4)
				o.Mat.SetTextureName(sc, "HouseBlanket")

				o = gi3d.AddNewObject(sc, ogp, "house_pillow1", "HousePillowOne")
				o.Pose.Pos.Set(-4.75, 1.15, -4)
				o.Mat.SetTextureName(sc, "HousePillow")

				o = gi3d.AddNewObject(sc, ogp, "house_couch_base1", "HouseCouchBaseOne")
				o.Pose.Pos.Set(2, 0, -3.5)
				o.Mat.SetTextureName(sc, "HouseCouch")

				o = gi3d.AddNewObject(sc, ogp, "house_couch_top1", "HouseCouchTopOne")
				o.Pose.Pos.Set(2, 1.5, -4.5)
				o.Mat.SetTextureName(sc, "HouseCouch")

				o = gi3d.AddNewObject(sc, ogp, "house_window1", "HouseWindowOne")
				o.Pose.Pos.Set(2, 3.5, -4.8)
				o.Mat.SetTextureName(sc, "HouseWindow")

				o = gi3d.AddNewObject(sc, ogp, "house_window2", "HouseWindowOne")
				o.Pose.Pos.Set(-3, 3.5, -4.8)
				o.Mat.SetTextureName(sc, "HouseWindow")

			case "Center_Blue":
				ogp = gi3d.AddNewGroup(sc, sc, nm)
				o := gi3d.AddNewObject(sc, ogp, "center_blue", "Center_Blue")
				o.Pose.Pos.Set(0, 0, 0)
				o.Mat.Color.SetString("blue", nil)
			case "Table":
				ogp = gi3d.AddNewGroup(sc, sc, nm)
				o := gi3d.AddNewObject(sc, ogp, "table", "Table")
				o.Pose.Pos.Set(0, 0, 0)

				o.Mat.SetTextureName(sc, "Table")
			case "BigComplex":
				ogp = gi3d.AddNewGroup(sc, sc, nm)
				o := gi3d.AddNewObject(sc, ogp, "bigComplexPlaceholder", "BigComplexPlaceholder")
				o.Pose.Pos.Set(0, 20, 0)
				// o.Mat.Color.SetString("red", nil)
				o.Mat.SetTextureName(sc, "Metal1")
			case "TestBed":
				ogp = gi3d.AddNewGroup(sc, sc, nm)
				err := sc.OpenObj([]string{"bed1.obj"}, ogp)
				if err != nil {
					log.Println(err)

				}
	*/

	if ogp != nil {
		ogp.Initial.Pos = obj.Pos
		// ogp.Initial.Scale = obj.Scale
	}
	return ogp
}

// todo: could have a smarter way of figuring out all the stuff you need
// to make, or not..
func (gm *Game) MakeLibrary() {
	gm.LibMakeBrickHouse()
	gm.LibMakeTheWall()
	gm.LibMakePerson()
	gm.LibMakePerson1()
	gm.LibMakeRoad()
	gm.LibMakeGrass()
	gm.LibMakeArenaWalls()
	gm.LibMakeLava()
	gm.LibMakeLavaContainer()
}

func (gm *Game) MakeMeshes() {
	sc := &gm.Scene.Scene
	gi3d.AddNewBox(sc, "Gun", 0.1, 0.1, 1)
	gi3d.AddNewBox(sc, "Person", 0.5, 2, 0.5)
	/*
		gi3d.AddNewBox(sc, "Hill", 1, 10, 1)
		gi3d.AddNewBox(sc, "Table", 5, 2.5, 5)
		gi3d.AddNewBox(sc, "Center_Blue", 3, 2, 3)
		gi3d.AddNewBox(sc, "HouseFloor", 10, 0.01, 10)
		gi3d.AddNewBox(sc, "HouseRoof", 10, 0.01, 10)
		gi3d.AddNewBox(sc, "HouseWallOne", 0.5, 10, 10)
		gi3d.AddNewBox(sc, "HouseWallTwo", 10, 10, 0.5)
		gi3d.AddNewBox(sc, "HouseBedOne", 3, 2, 2)
		gi3d.AddNewBox(sc, "HouseBlanketOne", 3, 0.1, 2)
		gi3d.AddNewBox(sc, "HousePillowOne", 0.5, 0.25, 2)
		gi3d.AddNewBox(sc, "HouseCouchBaseOne", 5, 2, 3)
		gi3d.AddNewBox(sc, "HouseCouchTopOne", 5, 1, 1)
		gi3d.AddNewBox(sc, "HouseWindowOne", 1, 1, 0.1)
		gi3d.AddNewBox(sc, "BigComplexPlaceholder", 40, 40, 40)
	*/
}

// func MakeTextures(sc *gi3d.Scene) {
// 	gi3d.AddNewTextureFile(sc, "Table", "table.jpg")
// 	gi3d.AddNewTextureFile(sc, "HouseBed", "bed.png")
// 	gi3d.AddNewTextureFile(sc, "HouseBlanket", "blanket.png")
// 	gi3d.AddNewTextureFile(sc, "HousePillow", "pillow.png")
// 	gi3d.AddNewTextureFile(sc, "HouseCouch", "couch.jpg")
// 	gi3d.AddNewTextureFile(sc, "HouseWindow", "window.png")
// 	gi3d.AddNewTextureFile(sc, "Metal1", "metal1.jpg")
// 	gi3d.AddNewTextureFile(sc, "Brick1", "brick.jpg")
// }

func (gm *Game) GetRandomSpawnPoint() (posX, posZ float32) {
	rand.Seed(time.Now().UTC().UnixNano())
	randomNum := int(rand.Intn(len(gm.SpawnPositions)))
	// fmt.Printf("Random number: %v\n", randomNum)
	return gm.SpawnPositions[randomNum].X, gm.SpawnPositions[randomNum].Z

}

// MakeWorld constructs a new virtual physics world
func (gm *Game) MakeWorld() {
	gm.World = &eve.Group{}
	gm.World.InitName(gm.World, "World")
	gm.BuildMap()
	gm.World.WorldInit() // key to put things in their places!
}

// MakeView makes the view
func (gm *Game) MakeView() {
	sc := &gm.Scene.Scene
	gm.MakeMeshes()
	wgp := gi3d.AddNewGroup(sc, sc, "world")
	gm.View = evev.NewView(gm.World, sc, wgp)
	// gm.View.InitLibrary() // this makes a basic library based on body shapes, sizes
	gm.MakeLibrary()
	gm.View.Sync()
}

func (gm *Game) Config() {
	gm.SpawnPositions = []mat32.Vec3{mat32.Vec3{-50, 0, -50}, mat32.Vec3{-50, 0, 50}, mat32.Vec3{50, 0, 50}, mat32.Vec3{50, 0, -50}}
	gamerow := gi.AddNewLayout(mfr2, "gamerow", gi.LayoutVert)
	gamerow.SetStretchMaxWidth()
	gamerow.SetStretchMaxHeight()

	sc := AddNewScene(gamerow, "scene")
	gm.Scene = sc
	sc.SetStretchMaxWidth()
	sc.SetStretchMaxHeight()

	// first, add lights, set camera
	sc.BgColor.SetUInt8(230, 230, 255, 255) // sky blue-ish
	gi3d.AddNewAmbientLight(&sc.Scene, "ambient", 0.5, gi3d.DirectSun)

	dir := gi3d.AddNewDirLight(&sc.Scene, "dir", 1, gi3d.DirectSun)
	dir.Pos.Set(0, 1, 1) // default: 0,1,1 = above and behind us (we are at 0,0,X)

	// point := gi3d.AddNewPointLight(sc, "point", 1, gi3d.DirectSun)
	// point.Pos.Set(0, 5, 5)

	// spot := gi3d.AddNewSpotLight(sc, "spot", 1, gi3d.DirectSun)
	// spot.Pose.Pos.Set(0, 0, 2)
	sc.Camera.FOV = 50
	sc.Camera.Pose.Pos.Y = 2
	sc.Camera.Pose.Pos.Z = 45
	gm.Gravity = 0.5
	gm.Map = currentMap
	gm.MakeWorld()

	gm.MakeView()
	gm.AbleToFire = true
	crossGroup := gi3d.AddNewGroup(&sc.Scene, &sc.Scene, "CrossGroup")
	color := gi.Color{255, 255, 255, 200}
	right := gi3d.AddNewLine(&sc.Scene, crossGroup, "RightLine", mat32.Vec3{0.1, 0, 0}, mat32.Vec3{0.5, 0, 0}, 0.1, color)
	left := gi3d.AddNewLine(&sc.Scene, crossGroup, "LeftLine", mat32.Vec3{-0.1, 0, 0}, mat32.Vec3{-0.5, 0, 0}, 0.1, color)
	bottom := gi3d.AddNewLine(&sc.Scene, crossGroup, "BottomLine", mat32.Vec3{0, -0.1, 0}, mat32.Vec3{0, -0.5, 0}, 0.1, color)
	top := gi3d.AddNewLine(&sc.Scene, crossGroup, "TopLine", mat32.Vec3{0, 0.1, 0}, mat32.Vec3{0, 0.5, 0}, 0.1, color)
	right.Mat.Emissive = color
	left.Mat.Emissive = color
	bottom.Mat.Emissive = color
	top.Mat.Emissive = color
	// text := gi3d.AddNewText2D(&sc.Scene, &sc.Scene, "CrossText", "+")
	// text.SetProp("color", "white")
	// // text.SetProp("background-color", "black")
	// text.SetProp("text-align", "center")
	// text.Pose.Scale.SetScalar(0.1)
	// text.Pose.Pos = sc.Camera.Pose.Pos
	// text.Pose.Pos.Z -= 10
	// text.Pose.Pos.X += 2
	joinPlayersTable(curBattleTerritory1 + curBattleTerritory2)

	gi.AddNewSpace(gamerow, "space1")

	brow := gi.AddNewLayout(mfr2, "brow", gi.LayoutHoriz)
	brow.SetProp("spacing", units.NewEx(2))
	brow.SetProp("horizontal-align", gi.AlignLeft)
	brow.SetStretchMaxWidth()

	epbut := gi.AddNewButton(brow, "edit-phys")
	epbut.SetText("Edit Phys")
	epbut.ButtonSig.Connect(gm.World.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			giv.GoGiEditorDialog(gm.World)
		}
	})
	cgbut := gi.AddNewButton(brow, "close-game")
	cgbut.SetText("Close Game")
	cgbut.ButtonSig.Connect(gm.World.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			gm.GameOn = false
			go removePlayer()
			tabIndex, _ := tv.TabIndexByName("<b>Game</b>")
			tv.DeleteTabIndex(tabIndex, true)
			tv.SelectTabIndex(0)
		}
	})

	rec := ki.Node{}
	rec.InitName(&rec, "rec")

	takeDamage := gi.AddNewButton(brow, "takeDamage")
	takeDamage.Text = "Take Damage"
	takeDamage.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			go gm.fireWeapon()
		}
	})
	// center_bluem :=
	// cbm.Segs.Set(10, 10, 10) // not clear if any diff really..
	// fpobj = gm.MakeObj(&MapObj{"FirstPerson", mat32.Vec3{0,0,0}, mat32.Vec3{1,1,1}}, "FirstPerson")
	// rcb = gm.MakeObj(&MapObj{"FirstPerson", mat32.Vec3{0, 0, 10}, mat32.Vec3{1, 1, 1}}, "FirstPerson")
	// gm.World.WorldInit()
	// rcb = gi3d.AddNewSolid(&sc.Scene, fpobj, "red-cube", "Person")
	// rcb.Pose.Pos.Set(0, -1, -8)
	// // rcb.Pose.Scale.Set(0.1, 0.1, 1)
	// rcb.Mat.Color.SetString("red", nil)

	// center_blue := sc.AddNewObject("center_blue", center_bluem.Name())
	// center_blue.Pose.Pos.Set(0, 0, 0)
	// center_blue.Mat.Color.SetString("blue", nil)
	//
	// green_hill := sc.AddNewObject("green_hill", cbm.Name())
	// green_hill.Pose.Pos.Set(1, 0, -20)
	// green_hill.Pose.Scale.Set(1, 10, 1)
	// green_hill.Mat.Color.SetString("green", nil)
	//
	// tbtx := gi3d.AddNewTextureFile(&sc.Scene, "table", "table.jpg")
	// var posy float32 = 0
	// var posx float32 = -5
	// var posz float32 = -20
	// for r := 0; r < 4; r++ {
	// 	posx = -5
	// 	for c := 0; c < 4; c++ {
	// 		market := sc.AddNewObject(fmt.Sprintf("market%v", c*r), cbm.Name())
	// 		market.Pose.Pos.Set(posx, posy, posz)
	// 		market.Pose.Scale.Set(5, 2.5, 5)
	// 		// market1.Mat.Color.SetString("red", nil)
	// 		market.Mat.SetTexture(&sc.Scene, tbtx.Name())
	// 		posx = posx - 6.5
	// 	}
	// 	posz = posz - 6.5
	//
	// }
	//
	// // market1 := sc.AddNewObject("market1", cbm.Name())
	// // market1.Pose.Pos.Set(-5, 0, -20)
	// // market1.Pose.Scale.Set(5, 2.5, 5)
	// // // market1.Mat.Color.SetString("red", nil)
	// // market1.Mat.SetTexture(&sc.Scene, tbtx.Name())
	gm.clearAllBullets()
	ogp := eve.AddNewGroup(gm.World, "Grass")
	gm.PhysMakeGrass(ogp, "Grass")
	// floorp := gi3d.AddNewPlane(&sc.Scene, "floor-plane", 200, 200)
	// floor := gi3d.AddNewSolid(&sc.Scene, &sc.Scene, "floor", floorp.Name())
	// floor.Pose.Pos.Set(0, 0, 0)
	// // floor.Mat.Emissive.SetString("green", nil)
	// grtx := gi3d.AddNewTextureFile(&sc.Scene, "ground", "objs/grass.jpg")
	// floor.Mat.SetTexture(&sc.Scene, grtx)
	// floor.Mat.Tiling.Repeat.Set(50, 50)

	gi.FilterLaggyKeyEvents = true // fix key lag

	gm.PosUpdtChan = make(chan bool) // todo: close channel when ending game, will terminate goroutines
	gm.FireUpdtChan = make(chan bool)

	gm.OtherPos = make(map[string]*CurPosition)
	gm.FireEvents = make([]*FireEventInfo, 0)
	gm.GameOn = true
	RayGroup := gm.Scene.Scene.ChildByName("RayGroup", 0).(*gi3d.Group)
	for i := 0; i < 30; i++ {
		color := gi.Color{255, 0, 0, 200}
		line := gi3d.AddNewLine(&gm.Scene.Scene, RayGroup, fmt.Sprintf("bullet_arrow_enemy%v", i), mat32.Vec3{0, 0, 0}, mat32.Vec3{1, 1, 1}, .05, color)
		line.SetInvisible()
	}

	go gm.GetPosFromServer()     // this is loop getting positions from server
	go gm.UpdatePeopleWorldPos() // this is loop updating positions
	go gm.UpdatePersonYPos()     // deals with jumping and gravity
	go gm.GetFireEvents()
	go gm.RenderEnemyShots()
}

func (gm *Game) RenderEnemyShots() {
	RayGroup := gm.Scene.Scene.ChildByName("RayGroup", 0).(*gi3d.Group)
	for {
		// startTime := time.Now()
		if !gm.GameOn {
			return
		}
		// fmt.Printf("Before fire event mu \n")
		gm.FireEventMu.Lock()
		// fmt.Printf("After fire event mu \n")
		// fmt.Printf("Lock time: %v \n", time.Since(startTime).Milliseconds())
		for k, d := range gm.FireEvents {
			if d.Creator != USER {
				bi := k % 30
				rayObj := RayGroup.ChildByName(fmt.Sprintf("bullet_arrow_enemy%v", bi), 0).(*gi3d.Solid)
				ray := mat32.NewRay(d.Origin, d.Dir)
				endPos := mat32.Vec3{0, 0, 0}
				sepPos := d.Dir.Mul(mat32.Vec3{100, 100, 100})
				cts := gm.World.RayBodyIntersections(*ray)
				killed := false
				for _, d1 := range cts {
					if d1.Body.Name() == "FirstPerson" && !gm.IsDead {
						gm.FireEventMu.Unlock()
						gm.removeHealthPoints(d.Damage, d.Creator)
						gi3d.SetLineStartEnd(rayObj, mat32.Vec3{500, 500, 500}, mat32.Vec3{500, 500, 500})
						gm.FireEventMu.Lock()
						killed = true
					}

					endPos = d1.Point
					break
				}
				if cts == nil {

					endPos = sepPos
				}
				if !killed {
					gi3d.SetLineStartEnd(rayObj, d.Origin, endPos)
					rayObj.ClearInvisible()
				}
			}
			if time.Now().Sub(d.StartTime) >= time.Millisecond*BulletTimeMsec {
				bi := k % 30
				rayObjKi := RayGroup.ChildByName(fmt.Sprintf("bullet_arrow_enemy%v", bi), 0)
				gm.FireEvents[k].Origin = mat32.Vec3{500, 500, 500}
				gm.FireEvents[k].Dir = mat32.Vec3{-1, 0, 0}
				if rayObjKi == nil {
					continue
				}
				rayObj := rayObjKi.(*gi3d.Solid)
				rayObj.SetInvisible()
				gi3d.SetLineStartEnd(rayObj, mat32.Vec3{500, 500, 500}, mat32.Vec3{500, 500, 500})
			}
		}

		gm.FireEventMu.Unlock()
		// fmt.Printf("Total time for render enemy shots: %v \n", time.Since(startTime).Milliseconds())
	}
}
func (gm *Game) fireWeapon() { // standard event for what happens when you fire
	// Currently just deal damage to yourself, at interval
	// todo: actually fire and deal damage to others
	if !gm.AbleToFire {
		return
	}
	// what to do on fire in here:

	cursor := gm.Scene.Scene.ChildByName("CrossGroup", 0).(*gi3d.Group)
	endPos := cursor.Pose
	rayPos := cursor.Pose
	sPos := cursor.Pose
	// endPos.MoveOnAxis(0, 0, -1, 100)
	rayPos.Pos = mat32.Vec3{0, 0, 0}
	rayPos.MoveOnAxis(0, 0, -1, 1)
	ray := mat32.NewRay(cursor.Pose.Pos, rayPos.Pos)
	// fmt.Printf("Ray: %v \n", ray)
	cts := gm.World.RayBodyIntersections(*ray)
	var closest *eve.BodyPoint
	// fmt.Printf("Contacts: %v \n", cts)
	for _, d := range cts {
		// fmt.Printf("Key: %v Body: %v  Point: %v \n", k, d.Body, d.Point)
		if closest == nil {
			closest = d
		} else {
			if cursor.Pose.Pos.DistTo(closest.Point) > cursor.Pose.Pos.DistTo(d.Point) {
				closest = d
			}
		}
		if d.Body.Name() == "FirstPerson" {
			// gm.removeHealthPoints(WEAPON)
		}
	}
	if cts != nil {
		endPos.Pos = closest.Point
	} else {
		sPos.MoveOnAxis(0, 0, -1, 100)
		endPos.Pos = sPos.Pos
	}
	index := len(gm.FireEvents)
	RayGroup := gm.Scene.Scene.ChildByName("RayGroup", 0).(*gi3d.Group)
	index = index % 30
	// fmt.Printf("Name: bullet_arrow_enemy%v \n", index)
	rayObj := RayGroup.ChildByName(fmt.Sprintf("bullet_arrow_enemy%v", index), 0).(*gi3d.Solid)
	gi3d.SetLineStartEnd(rayObj, cursor.Pose.Pos, endPos.Pos)
	rayObj.ClearInvisible()
	// bullet = gi3d.AddNewLine(&gm.Scene.Scene, RayGroup, "bullet_arrow_you", cursor.Pose.Pos, endPos.Pos, .05, color)
	go gm.removeBulletLoop(rayObj, cursor.Pose.Pos, rayPos.Pos, index)
	// done with what to fire
	gm.AbleToFire = false
	writeFireEventToServer(cursor.Pose.Pos, rayPos.Pos, generateDamageAmount(WEAPON), CURBATTLE)
	addFireEventToDB(USER, generateDamageAmount(WEAPON), cursor.Pose.Pos, rayPos.Pos)
	numOfSeconds := TheWeapons[WEAPON].FireRate
	time.Sleep(time.Duration(1/numOfSeconds) * time.Second)
	gm.AbleToFire = true
}

func (gm *Game) removeBulletLoop(bullet *gi3d.Solid, origin mat32.Vec3, dir mat32.Vec3, index int) {
	time.Sleep(300 * time.Millisecond)
	gm.FireEventMu.Lock()
	// delete(gm.FireEvents, index)
	// removeFireEventFromServer(index, CURBATTLE)
	// removeBulletFromDB(origin, dir)
	bullet.SetInvisible()
	gi3d.SetLineStartEnd(bullet, mat32.Vec3{500, 500, 500}, mat32.Vec3{500, 500, 500})
	gm.FireEventMu.Unlock()
}

func generateDamageAmount(wp string) (damage int) {
	var minD, maxD, rangeDiff float32
	minD = TheWeapons[wp].MinD
	maxD = TheWeapons[wp].MaxD
	rangeDiff = maxD - minD
	randNum := rand.Float32()
	addNum := math.Round(float64(rangeDiff * randNum))
	damage = int(float32(addNum) + minD)
	return damage
}
func (gm *Game) removeHealthPoints(dmg int, from string) {
	HEALTH -= float32(dmg)
	healthBar.SetValue(HEALTH)
	healthText.SetText(fmt.Sprintf("You have %v health", HEALTH))
	if HEALTH >= 67 {
		healthBar.SetProp(":value", ki.Props{"background-color": "green"})
	} else if HEALTH <= 33 {
		healthBar.SetProp(":value", ki.Props{"background-color": "red"})
	} else {
		healthBar.SetProp(":value", ki.Props{"background-color": "yellow"})
	}
	healthBar.SetFullReRender()
	if HEALTH <= 0 {
		gm.SpawnCount += 1
		gm.IsDead = true
		pers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)
		camOff := gm.Scene.Camera.Pose.Pos.Sub(pers.Rel.Pos) // currrent offset of camera vs. person
		pers.Rel.Pos = mat32.Vec3{1000, 1, 1000}
		gm.Scene.Camera.Pose.Pos = pers.Rel.Pos.Add(camOff)
		gm.KilledBy = from
		writePlayerPosToServer(pers.Rel.Pos, CURBATTLE)
		gm.World.WorldRelToAbs()
		gm.Scene.UpdateSig()
		gm.Scene.Win.OSWin.SetCursorEnabled(true, false)
		gm.Scene.TrackMouse = false
		resultText.SetText("<b>You were killed by " + from + " - Respawning in 5</b>")
		resultText.SetFullReRender()
		// updateBattlePoints(from, gm.OtherPos[from].Points+1)
		go gm.timerForResult(from)
	}
}
func (gm *Game) timerForResult(from string) {
	time.Sleep(1 * time.Second)
	resultText.SetText("<b>You were killed by " + from + " - Respawning in 4</b>")
	resultText.SetFullReRender()
	time.Sleep(1 * time.Second)
	resultText.SetText("<b>You were killed by " + from + " - Respawning in 3</b>")
	resultText.SetFullReRender()
	time.Sleep(1 * time.Second)
	resultText.SetText("<b>You were killed by " + from + " - Respawning in 2</b>")
	resultText.SetFullReRender()
	time.Sleep(1 * time.Second)
	resultText.SetText("<b>You were killed by " + from + " - Respawning in 1</b>")
	resultText.SetFullReRender()
	time.Sleep(1 * time.Second)
	resultText.SetText("<b>You were killed by " + from + "</b>")
	updt := resultRow.UpdateStart()
	defer resultRow.UpdateEnd(updt)
	resultButton := gi.AddNewButton(resultRow, "resultButton")
	resultButton.Text = "<b>Respawn</b>"
	resultButton.SetProp("horizontal-align", "center")
	resultButton.SetProp("font-size", "40px")
	resultButton.SetFullReRender()
	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	resultButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			gm.IsDead = false
			gm.clearAllBullets()
			resultText.SetText("")
			resultText.SetFullReRender()
			resultButton.Delete(true)
			HEALTH = 100
			healthBar.SetValue(HEALTH)
			healthBar.SetProp(":value", ki.Props{"background-color": "green"})
			healthBar.SetFullReRender()
			healthText.SetText(fmt.Sprintf("You have %v health", HEALTH))
			pers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)
			camOff := gm.Scene.Camera.Pose.Pos.Sub(pers.Rel.Pos) // currrent offset of camera vs. person
			gm.Scene.TrackMouse = true
			gm.Scene.Win.OSWin.SetCursorEnabled(false, true)
			posX, posY := gm.GetRandomSpawnPoint()
			pers.Rel.Pos = mat32.Vec3{posX, 1, posY}
			gm.Scene.Camera.Pose.Pos = pers.Rel.Pos.Add(camOff)
			gm.World.WorldRelToAbs()
			gm.Scene.UpdateSig()

		}
	})
	resultButton.SetFullReRender()
	resultText.SetFullReRender()
	resultRow.SetFullReRender()
}

func (gm *Game) updateCursorPosition() {
	cursor := gm.Scene.Scene.ChildByName("CrossGroup", 0).(*gi3d.Group)
	// pers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)
	cursor.Pose = gm.Scene.Camera.Pose
	cursor.Pose.MoveOnAxis(0, 0, -1, 3)
	cursor.Pose.Scale.SetScalar(0.3)
}

func (gm *Game) UpdatePersonYPos() {
	for {
		if !gm.GameOn {
			return
		}
		pers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)
		camOff := gm.Scene.Camera.Pose.Pos.Sub(pers.Rel.Pos) // currrent offset of camera vs. person
		if pers.Rel.Pos.Y != 1 {
			pers.Rel.LinVel.Y = pers.Rel.LinVel.Y - gm.Gravity
			pers.Rel.Pos.Y += pers.Rel.LinVel.Y
			if pers.Rel.Pos.Y <= 1 {
				pers.Rel.Pos.Y = 1
				pers.Rel.LinVel.Y = 0
			}
			writePlayerPosToServer(pers.Rel.Pos, CURBATTLE)
		}
		gm.World.WorldRelToAbs()
		gm.Scene.Camera.Pose.Pos = pers.Rel.Pos.Add(camOff)
		if pers.Rel.Pos.Y != 1 {
			gm.Scene.UpdateSig()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func AddNewScene(parent ki.Ki, name string) *Scene {
	sc := parent.AddNewChild(KiT_Scene, name).(*Scene)
	sc.Defaults()
	return sc
}

func (gm *Game) UpdatePeopleWorldPos() {
	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	// Get all of the groups
	pGp := gm.World.ChildByName("PeopleGroup", 0).(*eve.Group)
	pgt := gm.Scene.Scene.ChildByName("PeopleTextGroup", 0)
	uk := mfr2.ChildByName("usernameKey", 0)
	for i := 0; true; i++ {
		_, ok := <-gm.PosUpdtChan // we wait here to receive channel message sent when positions have been updated
		if !ok {                  // this means channel was closed, we need to bail, game over!
			return
		}
		gm.PosMu.Lock()
		keys := make([]string, len(gm.OtherPos))
		ctr := 0
		for k := range gm.OtherPos {
			keys[ctr] = k
			ctr++
		}
		sort.Strings(keys) // it is "key" to have others in same order so if there are no changes, nothing happens
		// config := kit.TypeAndNameList{}
		// config1 := kit.TypeAndNameList{}
		// for _, k := range keys {
		// 	config.Add(eve.KiT_Group, k)
		// 	if uk.ChildByName("ukt_"+k, 0) != nil {
		// 		config1.Add(gi.KiT_Label, "ukt_"+k)
		// 	}
		// 	// if uk.ChildByName("ukt_"+k+"_button", 0) != nil {
		// 	// 	config1.Add(gi.KiT_Button, "ukt_"+k+"_button")
		// 	// }
		// 	if i >= 1 {
		// 		config1.Add(gi.KiT_Label, "ukt_"+USER)
		// 		// config1.Add(gi.KiT_Button, "ukt_"+USER+"_button")
		// 	}
		// }
		// mods, updt := pGp.ConfigChildren(config, ki.NonUniqueNames)
		// mods1, updt1 := pgt.ConfigChildren(config, ki.NonUniqueNames)
		// mods2, updt2 := uk.ConfigChildren(config1, ki.NonUniqueNames)
		// if !mods {
		// 	updt = pGp.UpdateStart() // updt is automatically set if mods = true, so we're just doing it here
		// }
		// if !mods1 {
		// 	updt1 = pgt.UpdateStart() // updt is automatically set if mods = true, so we're just doing it here
		// }
		// if !mods2 {
		// 	updt2 = uk.UpdateStart() // updt is automatically set if mods = true, so we're just doing it here
		// }
		updt := pGp.UpdateStart()
		updt1 := pgt.UpdateStart()
		updt2 := uk.UpdateStart()
		needToSync := false
		for i, k := range keys {
			if k == USER {
				continue
			}
			ppos := gm.OtherPos[k]
			var pers *eve.Group
			if i >= len(*pGp.Children()) {
				pers = eve.AddNewGroup(pGp, k)
			} else {
				pers = pGp.Child(i).(*eve.Group) // this is guaranteed to be for person "k"
			}
			if !pers.HasChildren() { // if has not already been made
				needToSync = true
				gm.PhysMakePerson(pers, k, false) // make
				text := gi3d.AddNewText2D(&gm.Scene.Scene, &gm.Scene.Scene, k+"Text", k)
				text.SetProp("color", "black")
				text.SetProp("background-color", "white")
				text.SetProp("text-align", gi.AlignLeft)
				text.Pose.Scale.SetScalar(0.3)
				text.Pose.Pos = ppos.Pos
				text.Pose.Pos.Y = text.Pose.Pos.Y + 1.3
				uktn := "ukt_" + k
				ukt := gi.AddNewLabel(uk, uktn, "")
				ukt.SetText(fmt.Sprintf("<b>%v:</b>         %v kills         ", k, gm.OtherPos[k].Points))
				ukt.SetProp("font-size", "30px")
				ukt.Redrawable = true
				ukt.SetProp("width", "20em")
				ukt.SetFullReRender()
			} else {
				text1 := gm.Scene.Scene.ChildByName(k+"Text", 0)
				if text1 == nil {
					fmt.Printf("Continue \n")
					continue
				}
				text := text1.(*gi3d.Text2D)
				text.Pose.Pos = ppos.Pos
				text.Pose.Pos.Y = text.Pose.Pos.Y + 1.3
				text.SetProp("text-align", gi.AlignLeft)
				uktt, err := uk.ChildByNameTry("ukt_"+k, 0)
				if err != nil {
					panic(err)
				}
				ukt := uktt.(*gi.Label)
				ukt.SetText(fmt.Sprintf("<b>%v:</b>         %v kills              ", k, gm.OtherPos[k].Points))
				ukt.SetProp("width", "20em")
				text.Pose.Pos.X = text.Pose.Pos.X - 0.2
				if gm.OtherPos[k].Points >= 100 {
					gm.PosMu.Unlock()
					gm.setGameOver(k)
					gm.PosMu.Lock()
				}
			}
			pers.Rel.Pos = ppos.Pos
			firstPers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)
			d := pers.Rel.Pos.Sub(firstPers.Rel.Pos)
			dn := d.Normal()
			text1 := gm.Scene.Scene.ChildByName(k+"Text", 0)
			if text1 == nil {
				continue
			}
			text := text1.(*gi3d.Text2D)
			text.Pose.Quat.SetFromUnitVectors(mat32.Vec3{1, 0, 0}, dn)
			text.Pose.RotateOnAxis(0, 1, 0, -90)
		}
		_, err := uk.ChildByNameTry("ukt_"+USER, 0)
		if err != nil {
			ukt := gi.AddNewLabel(uk, "ukt_"+USER, "")
			ukt.SetText(fmt.Sprintf("<b>%v:</b>         %v kills              ", USER, POINTS))
			ukt.SetProp("font-size", "30px")
			ukt.Redrawable = true
			ukt.SetProp("width", "20em")
			ukt.SetFullReRender()
		} else {
			ukt := uk.ChildByName("ukt_"+USER, 0).(*gi.Label)
			ukt.SetText(fmt.Sprintf("<b>%v:</b>         %v kills            ", USER, POINTS))
			ukt.SetProp("width", "20em")
		}
		if POINTS >= 100 {
			gm.PosMu.Unlock()
			gm.setGameOver(USER)
			gm.PosMu.Lock()
		}
		if needToSync {
			gm.View.Sync() // if something was created or destroyed, it must use Sync to update Scene
		} else {
			gm.View.Sync()
			gm.View.UpdatePose() // UpdatePose is much faster and assumes no changes in objects
		}
		gm.PosMu.Unlock()
		// so now everyone's updated
		gm.World.WorldRelToAbs()
		pGp.UpdateEnd(updt)
		pgt.UpdateEnd(updt1)
		uk.UpdateEnd(updt2)
		gm.Scene.UpdateSig()
	}
}

func (sc *Scene) Render2D() {
	if sc.PushBounds() {
		if !sc.NoNav {
			sc.NavEvents()
		}
		if gi.Render2DTrace {
			// fmt.Printf("3D Render2D: %v\n", sc.PathUnique())
		}
		sc.Render()
		sc.PopBounds()
	} else {
		sc.DisconnectAllEvents(gi.RegPri)
	}
}

func (sc *Scene) NavEvents() {
	sc.ConnectEvent(oswin.MouseMoveEvent, gi.RegPri, func(recv, send ki.Ki, sig int64, d interface{}) {
		if !sc.TrackMouse {
			return
		}
		me := d.(*mouse.MoveEvent)
		me.SetProcessed()
		ssc := recv.Embed(KiT_Scene).(*Scene)
		orbDel := float32(.2)
		orbDels := orbDel * 1
		panDel := float32(.05)
		del := me.Where.Sub(me.From)
		dx := float32(-del.X)
		dy := float32(-del.Y)
		// fmt.Printf("pos: %v  fm: %v del: %v\n", me.Where, me.From, del)
		switch {
		case key.HasAllModifierBits(me.Modifiers, key.Shift):
			ssc.Camera.Pan(dx*panDel, -dy*panDel)
		case key.HasAllModifierBits(me.Modifiers, key.Control):
			ssc.Camera.PanAxis(dx*panDel, -dy*panDel)
		case key.HasAllModifierBits(me.Modifiers, key.Alt):
			ssc.Camera.PanTarget(dx*panDel, -dy*panDel, 0)
		default:
			if mat32.Abs(dx) > mat32.Abs(dy) {
				dy = 0
			} else {
				dx = 0
			}
			sc.CamRotUD += dy * orbDels
			if sc.CamRotUD > 90 {
				sc.CamRotUD = 90
			}
			if sc.CamRotUD < -90 {
				sc.CamRotUD = -90
			}

			sc.CamRotLR += dx * orbDels * 2
			pers := TheGame.World.ChildByName("FirstPerson", 0).(*eve.Group)
			ssc.Camera.Pose.Pos = pers.Abs.Pos
			ssc.Camera.Pose.Quat = pers.Abs.Quat
			ssc.Camera.Pose.Pos.Y += 1
			ssc.Camera.Pose.SetAxisRotation(0, 1, 0, sc.CamRotLR)
			ssc.Camera.Pose.RotateOnAxis(1, 0, 0, sc.CamRotUD)
			ssc.Camera.Pose.MoveOnAxis(0, 0, 1, 3)
			ssc.Camera.Pose.MoveOnAxis(1, 0, 0, 1)
			pers.Rel.SetAxisRotation(0, 1, 0, sc.CamRotLR)

			TheGame.updateCursorPosition()
		}
		ssc.UpdateSig()
	})
	// sc.ConnectEvent(oswin.MouseScrollEvent, gi.RegPri, func(recv, send ki.Ki, sig int64, d interface{}) {
	// 	me := d.(*mouse.ScrollEvent)
	// 	me.SetProcessed()
	// 	ssc := recv.Embed(KiT_Scene).(*Scene)
	// 	if ssc.SetDragCursor {
	// 		oswin.TheApp.Cursor(ssc.Viewport.Win.OSWin).Pop()
	// 		ssc.SetDragCursor = false
	// 	}
	// 	zoom := float32(me.NonZeroDelta(false))
	// 	zoomPct := float32(.05)
	// 	zoomDel := float32(.05)
	// 	switch {
	// 	case key.HasAllModifierBits(me.Modifiers, key.Alt):
	// 		ssc.Camera.PanTarget(0, 0, zoom*zoomDel)
	// 	default:
	// 		ssc.Camera.Zoom(zoomPct * zoom)
	// 	}
	// 	ssc.UpdateSig()
	// })
	sc.ConnectEvent(oswin.MouseEvent, gi.RegPri, func(recv, send ki.Ki, sig int64, d interface{}) {
		me := d.(*mouse.Event)
		me.SetProcessed()
		ssc := recv.Embed(KiT_Scene).(*Scene)

		if ssc.SetDragCursor {
			oswin.TheApp.Cursor(ssc.Viewport.Win.OSWin).Pop()
			ssc.SetDragCursor = false
		}
		if !ssc.IsInactive() && !ssc.HasFocus() {
			ssc.GrabFocus()
			sc.TrackMouse = true
			sc.Win.OSWin.SetCursorEnabled(false, true)
		} else {
			go TheGame.fireWeapon()
		}
		// obj := ssc.FirstContainingPoint(me.Where, true)
		// if me.Action == mouse.Release && me.Button == mouse.Right {
		// 	me.SetProcessed()
		// 	if obj != nil {
		// 		giv.StructViewDialog(ssc.Viewport, obj, giv.DlgOpts{Title: "sc Element View"}, nil, nil)
		// 	}
		// }
	})
	sc.ConnectEvent(oswin.MouseHoverEvent, gi.RegPri, func(recv, send ki.Ki, sig int64, d interface{}) {
		me := d.(*mouse.HoverEvent)
		me.SetProcessed()
		// ssc := recv.Embed(KiT_Scene).(*Scene)
		// obj := ssc.FirstContainingPoint(me.Where, true)
		// if obj != nil {
		// 	pos := me.Where
		// 	ttxt := fmt.Sprintf("element name: %v -- use right mouse click to edit", obj.Name())
		// 	gi.PopupTooltip(obj.Name(), pos.X, pos.Y, sc.Viewport, ttxt)
		// }
	})
	sc.ConnectEvent(oswin.KeyChordEvent, gi.RegPri, func(recv, send ki.Ki, sig int64, d interface{}) {
		ssc := recv.Embed(KiT_Scene).(*Scene)
		kt := d.(*key.ChordEvent)
		ssc.NavKeyEvents(kt)
	})
}

func (sc *Scene) NavKeyEvents(kt *key.ChordEvent) {
	ch := string(kt.Chord())
	// fmt.Printf(ch)
	// orbDeg := float32(5)
	// panDel := float32(.1)
	// zoomPct := float32(.05)

	gm := TheGame
	gm.WorldMu.Lock()
	defer gm.WorldMu.Unlock()

	wupdt := gm.World.UpdateStart()

	pers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)
	camOff := sc.Camera.Pose.Pos.Sub(pers.Rel.Pos) // currrent offset of camera vs. person
	// todo: get current camera axis-angle

	switch ch {
	case "Escape":
		sc.TrackMouse = !sc.TrackMouse
		if sc.TrackMouse {
			sc.Win.OSWin.SetCursorEnabled(false, true) // turn off mouse cursor
		} else {
			sc.Win.OSWin.SetCursorEnabled(true, false)
		}
		kt.SetProcessed()

		// case "UpArrow":
		//
		// 	sc.Camera.Pose.SetEulerRotation(orbDeg, 0, 0)
		// kt.SetProcessed()

		// case "Shift+UpArrow":
		// 	sc.Camera.Pan(0, panDel)
		// 	kt.SetProcessed()
		// case "Control+UpArrow":
		// 	sc.Camera.PanAxis(0, panDel)
		// 	kt.SetProcessed()
		// case "Alt+UpArrow":
		// 	sc.Camera.PanTarget(0, panDel, 0)
		// 	kt.SetProcessed()
		// case "DownArrow":
		// sc.Camera.Orbit(0, -orbDeg)
		// kt.SetProcessed()
		// case "Shift+DownArrow":
		// 	sc.Camera.Pan(0, -panDel)
		// 	kt.SetProcessed()
		// case "Control+DownArrow":
		// 	sc.Camera.PanAxis(0, -panDel)
		// 	kt.SetProcessed()
		// case "Alt+DownArrow":
		// 	sc.Camera.PanTarget(0, -panDel, 0)
		// 	kt.SetProcessed()
		// case "LeftArrow":
		// sc.Camera.Orbit(orbDeg, 0)
		// kt.SetProcessed()
		// case "Shift+LeftArrow":
		// 	sc.Camera.Pan(-panDel, 0)
		// 	kt.SetProcessed()
		// case "Control+LeftArrow":
		// 	sc.Camera.PanAxis(-panDel, 0)
		// 	kt.SetProcessed()
		// case "Alt+LeftArrow":
		// 	sc.Camera.PanTarget(-panDel, 0, 0)
		// 	kt.SetProcessed()
		// case "RightArrow":
		// sc.Camera.Orbit(-orbDeg, 0)
		// kt.SetProcessed()
		// case "Shift+RightArrow":
		// 	sc.Camera.Pan(panDel, 0)
		// 	kt.SetProcessed()
		// case "Control+RightArrow":
		// 	sc.Camera.PanAxis(panDel, 0)
		// 	kt.SetProcessed()
		// case "Alt+RightArrow":
		// 	sc.Camera.PanTarget(panDel, 0, 0)
		// 	kt.SetProcessed()
		// case "Alt++", "Alt+=":
		// 	sc.Camera.PanTarget(0, 0, panDel)
		// 	kt.SetProcessed()
		// case "Alt+-", "Alt+_":
		// 	sc.Camera.PanTarget(0, 0, -panDel)
		// 	kt.SetProcessed()
		// case "+", "=":
		// 	sc.Camera.Zoom(-zoomPct)
		// 	kt.SetProcessed()
		// case "-", "_":
		// 	sc.Camera.Zoom(zoomPct)
		// 	kt.SetProcessed()
	case " ":
		if pers.Rel.Pos.Y == 1 {
			pers.Rel.LinVel.Y = 1
			pers.Rel.Pos.Y += pers.Rel.LinVel.Y
		}
	case "r":
		pers.Rel.Pos.Set(0, 1, 0)
		pers.Rel.Quat.SetFromAxisAngle(mat32.Vec3{0, 1, 0}, 0)
	case "w":
		kt.SetProcessed()
		y := pers.Rel.Pos.Y // keep height fixed -- no jumping right now.
		if !gm.PersHitWall {
			pers.Rel.MoveOnAxis(0, 0, -1, .5) // todo: use camera axis not fixed axis
		} else {
			prevPosX := pers.Rel.Pos.X
			prevPosZ := pers.Rel.Pos.Z
			pers.Rel.MoveOnAxis(0, 0, -1, .5)
			stillNessecary := gm.WorldStep(true)
			if stillNessecary {
				pers.Rel.Pos.X = prevPosX
				pers.Rel.Pos.Z = prevPosZ
			}
		}
		pers.Rel.Pos.Y = y
		gm.WorldStep(false)
	case "s":
		kt.SetProcessed()
		y := pers.Rel.Pos.Y // keep height fixed -- no jumping right now.
		if !gm.PersHitWall {
			pers.Rel.MoveOnAxis(0, 0, 1, .5)
		} else {
			prevPosX := pers.Rel.Pos.X
			prevPosZ := pers.Rel.Pos.Z
			pers.Rel.MoveOnAxis(0, 0, 1, .5)
			stillNessecary := gm.WorldStep(true)
			// fmt.Printf("Still for s: %v \n", stillNessecary)
			if stillNessecary {
				pers.Rel.Pos.X = prevPosX
				pers.Rel.Pos.Z = prevPosZ
			}
		}

		pers.Rel.Pos.Y = y
		gm.WorldStep(false)
	case "a":
		kt.SetProcessed()
		y := pers.Rel.Pos.Y // keep height fixed -- no jumping right now.
		if !gm.PersHitWall {
			pers.Rel.MoveOnAxis(-0.75, 0, 0, .5)
		} else {
			prevPosX := pers.Rel.Pos.X
			prevPosZ := pers.Rel.Pos.Z
			pers.Rel.MoveOnAxis(-0.75, 0, 0, .5)
			stillNessecary := gm.WorldStep(true)
			if stillNessecary {
				pers.Rel.Pos.X = prevPosX
				pers.Rel.Pos.Z = prevPosZ
			}
		}
		pers.Rel.Pos.Y = y
		gm.WorldStep(false)
		// sc.Camera.Pan(panDel, 0)
		// kt.SetProcessed()
		// go updatePosition("posX", rcb.Initial.Pos.X)
	case "d":
		kt.SetProcessed()
		y := pers.Rel.Pos.Y // keep height fixed -- no jumping right now.
		if !gm.PersHitWall {
			pers.Rel.MoveOnAxis(0.75, 0, 0, .5)
		} else {
			prevPosX := pers.Rel.Pos.X
			prevPosZ := pers.Rel.Pos.Z
			pers.Rel.MoveOnAxis(0.75, 0, 0, .5)
			stillNessecary := gm.WorldStep(true)
			if stillNessecary {
				pers.Rel.Pos.X = prevPosX
				pers.Rel.Pos.Z = prevPosZ
			}
		}
		pers.Rel.Pos.Y = y
		gm.WorldStep(false)
		// sc.Camera.Pan(-panDel, 0)
		// kt.SetProcessed()
		// go updatePosition("posX", rcb.Initial.Pos.X)
	}

	go updatePosition("pos", pers.Abs.Pos) // this was updated from UpdateWorld
	// fmt.Printf("Pos Abs: %v  Pos Rel: %v \n", pers.Abs.Pos, pers.Rel.Pos)
	sc.Camera.Pose.Pos = pers.Rel.Pos.Add(camOff)
	gm.updateCursorPosition()
	// text.TrackCamera(&sc.Scene)
	// text.Pose.Pos.Z -= 10
	// text.Pose.Pos.X += 2

	gm.World.UpdateEnd(wupdt)
	gm.View.UpdatePose()
	sc.UpdateSig()
}

func (ev *Game) WorldStep(specialCheck bool) (stillNessecary bool) {
	// fmt.Printf("In world step! \n")
	ev.World.WorldRelToAbs()
	// var contacts eve.Contacts
	cts := ev.World.WorldCollide(eve.DynsTopGps)
	// fmt.Printf("Children: %v \n", ev.World.Children())
	// fmt.Printf("Cts: %v \n", cts)
	if cts == nil && specialCheck {
		ev.PersHitWall = false
		// fmt.Printf("We are safe now! \n")
	}
	for _, cl := range cts {
		// fmt.Printf("Cl: %v \n", cl)
		if len(cl) >= 1 {
			for _, c := range cl {
				if c.A.Name() == "FirstPerson" {
					// contacts = cl
					// fmt.Printf("Contacts: %v \n", contacts)
					// fmt.Printf("Contact: %v \n", c)
					name := c.B.Name()
					if strings.Contains(name, "wall") || strings.Contains(name, "Wall") {
						ev.PersHitWall = true
						// fmt.Printf("Hit wall! \n")
						if specialCheck {
							return true
						}
					}
					if strings.Contains(name, "Lava") {
						POINTS -= 1
						ev.removeHealthPoints(200, "Lava")
					}
				}
				// fmt.Printf("A: %v  B: %v\n", c.A.Name(), c.B.Name())
			}
		}
	}
	return false
	// if len(contacts) > 0 { // turn around
	// 	fmt.Printf("hit wall: turn around!\n")
	// }
	// ev.View.UpdatePose()
	// ev.Snapshot()
}