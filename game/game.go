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

type CurPosition struct {
	Username string
	Pos      mat32.Vec3
	Points   int
}

type Weapon struct {
	Name     string
	MinD     float32 `desc:"Minimum amount of damage that this weapon will do"`
	MaxD     float32 `desc:"Maximum amount of damage that this weapon will do"`
	FireRate float32 `desc:"How many times this weapon can be fired in a second"`
}

type FireEventInfo struct {
	Creator string
	Damage  int
	Origin  mat32.Vec3
	Dir     mat32.Vec3
}

type Game struct {
	World        *eve.Group
	View         *evev.View
	Scene        *Scene
	Map          Map
	MapObjs      map[string]bool
	OtherPos     map[string]*CurPosition
	PosUpdtChan  chan bool  `desc:"channel connecting server pos updates with world state update"`
	PosMu        sync.Mutex `desc:"protects updates to OtherPos map"`
	WorldMu      sync.Mutex `desc:"protects updates to World physics and view"`
	GameOn       bool       // starts on when game turned out, turn off when close game
	Winner       string
	PersHitWall  bool
	Gravity      float32
	AbleToFire   bool
	FireEvents   map[int]*FireEventInfo
	FireEventMu  sync.Mutex
	FireUpdtChan chan bool
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
	ogp.Initial.Pos.Set(0.25, 1, 40)

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
	gamerow := gi.AddNewLayout(playTab, "gamerow", gi.LayoutVert)
	gamerow.SetStretchMaxWidth()
	gamerow.SetStretchMaxHeight()

	sc := AddNewScene(gamerow, "scene")
	gm.Scene = sc
	sc.SetStretchMaxWidth()
	sc.SetStretchMaxHeight()
	sc.Win.OSWin.SetCursorEnabled(false, true)

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

	text := gi3d.AddNewText2D(&sc.Scene, &sc.Scene, "CrossText", "+")
	text.SetProp("color", "white")
	// text.SetProp("background-color", "black")
	text.SetProp("text-align", "center")
	text.Pose.Scale.SetScalar(0.1)
	text.Pose.Pos = sc.Camera.Pose.Pos
	text.Pose.Pos.Z -= 10
	text.Pose.Pos.X += 2

	gi.AddNewSpace(gamerow, "space1")

	brow := gi.AddNewLayout(playTab, "brow", gi.LayoutHoriz)
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
	clearAllBullets()
	floorp := gi3d.AddNewPlane(&sc.Scene, "floor-plane", 200, 200)
	floor := gi3d.AddNewSolid(&sc.Scene, &sc.Scene, "floor", floorp.Name())
	floor.Pose.Pos.Set(0, 0, 0)
	// floor.Mat.Emissive.SetString("green", nil)
	grtx := gi3d.AddNewTextureFile(&sc.Scene, "ground", "objs/grass.jpg")
	floor.Mat.SetTexture(&sc.Scene, grtx)
	floor.Mat.Tiling.Repeat.Set(50, 50)

	gi.FilterLaggyKeyEvents = true // fix key lag

	gm.PosUpdtChan = make(chan bool) // todo: close channel when ending game, will terminate goroutines
	gm.FireUpdtChan = make(chan bool)

	gm.OtherPos = make(map[string]*CurPosition)
	gm.FireEvents = make(map[int]*FireEventInfo)
	gm.GameOn = true
	RayGroup := gm.Scene.Scene.ChildByName("RayGroup", 0).(*gi3d.Group)
	for i := 0; i < 30; i++ {
		color, _ := gi.ColorFromName("red")
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
		_, ok := <-gm.FireUpdtChan // we wait here to receive channel message sent when positions have been updated
		if !ok {                   // this means channel was closed, we need to bail, game over!
			return
		}
		gm.FireEventMu.Lock()
		for i := 0; i < 30; i++ {
			if gm.FireEvents[i] == nil {
				rayObj := RayGroup.ChildByName(fmt.Sprintf("bullet_arrow_enemy%v", i), 0).(*gi3d.Solid)
				rayObj.SetInvisible()
				gi3d.SetLineStartEnd(rayObj, mat32.Vec3{500, 500, 500}, mat32.Vec3{500, 500, 500})
			}
		}

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
					if d1.Body.Name() == "FirstPerson" {
						gm.FireEventMu.Unlock()
						gm.removeHealthPoints(d.Damage, d.Creator)
						rayObj.SetInvisible()
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
		}

		gm.FireEventMu.Unlock()
	}
}
func (gm *Game) fireWeapon() { // standard event for what happens when you fire
	// Currently just deal damage to yourself, at interval
	// todo: actually fire and deal damage to others
	if !gm.AbleToFire {
		return
	}
	// what to do on fire in here:

	cursor := gm.Scene.Scene.ChildByName("CrossText", 0).(*gi3d.Text2D)
	endPos := cursor.Pose
	rayPos := cursor.Pose
	// endPos.MoveOnAxis(0, 0, -1, 100)
	rayPos.Pos = mat32.Vec3{0, 0, 0}
	rayPos.MoveOnAxis(0, 0, -1, 1)
	ray := mat32.NewRay(cursor.Pose.Pos, rayPos.Pos)
	// fmt.Printf("Ray: %v \n", ray)
	cts := gm.World.RayBodyIntersections(*ray)
	var closest *eve.BodyPoint
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
		rayPos.MoveOnAxis(0, 0, -1, 100)
		endPos.Pos = rayPos.Pos
	}
	var index int
	for index = 0; index < 29; index++ {
		if gm.FireEvents[index] == nil {
			break
		}
	}
	RayGroup := gm.Scene.Scene.ChildByName("RayGroup", 0).(*gi3d.Group)
	// fmt.Printf("Name: bullet_arrow_enemy%v \n", index)
	rayObj := RayGroup.ChildByName(fmt.Sprintf("bullet_arrow_enemy%v", index), 0).(*gi3d.Solid)
	gi3d.SetLineStartEnd(rayObj, cursor.Pose.Pos, endPos.Pos)
	rayObj.ClearInvisible()
	// bullet = gi3d.AddNewLine(&gm.Scene.Scene, RayGroup, "bullet_arrow_you", cursor.Pose.Pos, endPos.Pos, .05, color)
	go gm.removeBulletLoop(rayObj, cursor.Pose.Pos, rayPos.Pos)
	// done with what to fire
	gm.AbleToFire = false
	addFireEventToDB(USER, generateDamageAmount(WEAPON), cursor.Pose.Pos, rayPos.Pos)
	numOfSeconds := TheWeapons[WEAPON].FireRate
	time.Sleep(time.Duration(1/numOfSeconds) * time.Second)
	gm.AbleToFire = true
}

func (gm *Game) removeBulletLoop(bullet *gi3d.Solid, origin mat32.Vec3, dir mat32.Vec3) {
	gm.FireEventMu.Lock()
	time.Sleep(300 * time.Millisecond)
	removeBulletFromDB(origin, dir)
	for k, d := range gm.FireEvents {
		if d.Origin == origin && d.Dir == dir {
			delete(gm.FireEvents, k)
			break
		}
	}
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
		pers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)
		camOff := gm.Scene.Camera.Pose.Pos.Sub(pers.Rel.Pos) // currrent offset of camera vs. person
		pers.Rel.Pos = mat32.Vec3{1000, 1, 1000}
		updatePosition("pos", pers.Rel.Pos)
		gm.Scene.Camera.Pose.Pos = pers.Rel.Pos.Add(camOff)
		gm.World.WorldRelToAbs()
		gm.Scene.UpdateSig()
		gm.Scene.Win.OSWin.SetCursorEnabled(true, false)
		resultText.SetText("<b>You were killed by " + from + " - Respawning in 5</b>")
		resultText.SetFullReRender()
		updateBattlePoints(from, gm.OtherPos[from].Points+1)
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
	resultButton := gi.AddNewButton(resultRow, "resultButton")
	resultButton.Text = "<b>Respawn</b>"
	resultButton.SetProp("horizontal-align", "center")
	resultButton.SetProp("font-size", "40px")
	resultButton.SetFullReRender()
	rec := ki.Node{}
	rec.InitName(&rec, "rec")
	resultButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			clearAllBullets()
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
			gm.Scene.Win.OSWin.SetCursorEnabled(false, true)
			pers.Rel.Pos = mat32.Vec3{0, 1, 50}
			updatePosition("pos", pers.Rel.Pos)
			gm.Scene.Camera.Pose.Pos = pers.Rel.Pos.Add(camOff)
			gm.World.WorldRelToAbs()
			gm.Scene.UpdateSig()

		}
	})
	resultText.SetFullReRender()
}

func (gm *Game) updateCursorPosition() {
	cursor := gm.Scene.Scene.ChildByName("CrossText", 0).(*gi3d.Text2D)
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
		gm.PosMu.Lock()
		gm.WorldMu.Lock()
		pers := gm.World.ChildByName("FirstPerson", 0).(*eve.Group)

		camOff := gm.Scene.Camera.Pose.Pos.Sub(pers.Rel.Pos) // currrent offset of camera vs. person

		if pers.Rel.Pos.Y != 1 {
			pers.Rel.LinVel.Y = pers.Rel.LinVel.Y - gm.Gravity

			pers.Rel.Pos.Y += pers.Rel.LinVel.Y
			if pers.Rel.Pos.Y <= 1 {
				pers.Rel.Pos.Y = 1
				pers.Rel.LinVel.Y = 0
			}
			statement := fmt.Sprintf("UPDATE players SET posY = '%v' WHERE username='%v'", pers.Rel.Pos.Y, USER)
			_, err := db.Exec(statement)
			if err != nil {
				fmt.Printf("DB err: %v \n", err)
			}
		}
		gm.World.WorldRelToAbs()
		gm.WorldMu.Unlock()
		gm.PosMu.Unlock()
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
	pGp := gm.World.ChildByName("PeopleGroup", 0).(*eve.Group)
	pgt := gm.Scene.Scene.ChildByName("PeopleTextGroup", 0)
	uk := playTab.ChildByName("usernameKey", 0)
	for i := 0; true; i++ {
		_, ok := <-gm.PosUpdtChan // we wait here to receive channel message sent when positions have been updated
		if !ok {                  // this means channel was closed, we need to bail, game over!
			return
		}
		gm.PosMu.Lock()
		gm.WorldMu.Lock()
		keys := make([]string, len(gm.OtherPos))
		ctr := 0
		for k := range gm.OtherPos {
			keys[ctr] = k
			ctr++
		}
		sort.Strings(keys) // it is "key" to have others in same order so if there are no changes, nothing happens
		config := kit.TypeAndNameList{}
		config1 := kit.TypeAndNameList{}
		for _, k := range keys {
			config.Add(eve.KiT_Group, k)
			if uk.ChildByName("ukt_"+k, 0) != nil {
				config1.Add(gi.KiT_Label, "ukt_"+k)
			}
			if uk.ChildByName("ukt_"+k+"_button", 0) != nil {
				config1.Add(gi.KiT_Button, "ukt_"+k+"_button")
			}
			if i >= 1 {
				config1.Add(gi.KiT_Label, "ukt_"+USER)
				config1.Add(gi.KiT_Button, "ukt_"+USER+"_button")
			}
		}
		mods, updt := pGp.ConfigChildren(config, ki.NonUniqueNames)
		mods1, updt1 := pgt.ConfigChildren(config, ki.NonUniqueNames)
		mods2, updt2 := uk.ConfigChildren(config1, ki.NonUniqueNames)
		if !mods {
			updt = pGp.UpdateStart() // updt is automatically set if mods = true, so we're just doing it here
		}
		if !mods1 {
			updt1 = pgt.UpdateStart() // updt is automatically set if mods = true, so we're just doing it here
		}
		if !mods2 {
			updt2 = uk.UpdateStart() // updt is automatically set if mods = true, so we're just doing it here
		}
		// now, the children of pGp are the keys of OtherPos in order
		for i, k := range keys {
			ppos := gm.OtherPos[k]
			pers := pGp.Child(i).(*eve.Group) // this is guaranteed to be for person "k"
			if !pers.HasChildren() {          // if has not already been made
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
				addPointsButton := gi.AddNewButton(uk, uktn+"_button")
				addPointsButton.SetText("Plus 1 point")
				addPointsButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
					if sig == int64(gi.ButtonClicked) {
						if gm.OtherPos[k] != nil {
							gm.OtherPos[k].Points = gm.OtherPos[k].Points + 1
							updateBattlePoints(k, gm.OtherPos[k].Points)
							if gm.OtherPos[k].Points >= 10 {
								gm.setGameOver(k)
							}
						} else {
							POINTS = POINTS + 1
							updateBattlePoints(USER, POINTS)
						}

					}
				})

				// text.Pose.Pos.X = text.Pose.Pos.X - 0.2

				// fmt.Printf("Text: %v    Pos: %v    Text: %v\n", text, text.Pose.Pos, text.Text)
			} else {
				text1 := gm.Scene.Scene.ChildByName(k+"Text", 0)
				if text1 == nil {
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
				text.Pose.Pos.X = text.Pose.Pos.X - 0.2
				if gm.OtherPos[k].Points >= 10 {
					gm.PosMu.Unlock()
					gm.WorldMu.Unlock()
					gm.setGameOver(k)
					gm.PosMu.Lock()
					gm.WorldMu.Lock()
				}
			}
			pers.Rel.Pos = ppos.Pos
		}

		_, err := uk.ChildByNameTry("ukt_"+USER, 0)
		if err != nil {
			// fmt.Printf("Points: %v", POINTS)
			ukt := gi.AddNewLabel(uk, "ukt_"+USER, "")
			ukt.SetText(fmt.Sprintf("<b>%v:</b>         %v kills              ", USER, POINTS))
			ukt.SetProp("font-size", "30px")
			ukt.Redrawable = true
			addPointsButton := gi.AddNewButton(uk, "ukt_"+USER+"_button")
			addPointsButton.SetText("Plus 1 point")
			addPointsButton.ButtonSig.Connect(rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.ButtonClicked) {
					POINTS = POINTS + 1
					updateBattlePoints(USER, POINTS)
					if POINTS >= 10 {
						gm.setGameOver(USER)
					}
				}

			})
		} else {
			ukt := uk.ChildByName("ukt_"+USER, 0).(*gi.Label)
			ukt.SetText(fmt.Sprintf("<b>%v:</b>         %v kills            ", USER, POINTS))
		}
		if POINTS >= 10 {
			gm.setGameOver(USER)
		}
		if mods {
			gm.View.Sync() // if something was created or destroyed, it must use Sync to update Scene
		} else {
			gm.View.UpdatePose() // UpdatePose is much faster and assumes no changes in objects
		}
		gm.PosMu.Unlock()
		// so now everyone's updated
		gm.World.WorldRelToAbs()
		pGp.UpdateEnd(updt)
		pgt.UpdateEnd(updt1)
		uk.UpdateEnd(updt2)

		gm.WorldMu.Unlock()
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
