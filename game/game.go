package main

import (
	"fmt"

	"github.com/goki/gi/gi"
	"github.com/goki/gi/gi3d"
	"github.com/goki/gi/oswin"
	"github.com/goki/gi/oswin/key"
	"github.com/goki/gi/oswin/mouse"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

type Scene struct {
	gi3d.Scene
}

var KiT_Scene = kit.Types.AddType(&Scene{}, nil)

func startGame() {
	scrow := gi.AddNewLayout(signUpTab, "scrow", gi.LayoutHoriz)
	scrow.SetStretchMaxWidth()
	scrow.SetStretchMaxHeight()

	sc := AddNewScene(scrow, "scene")
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

	cbm := gi3d.AddNewBox(&sc.Scene, "cube1", 1, 1, 1)
	// cbm.Segs.Set(10, 10, 10) // not clear if any diff really..

	rcb := sc.AddNewObject("red-cube", cbm.Name())
	rcb.Pose.Pos.Set(-1, 0, 0)
	rcb.Mat.Color.SetString("red", nil)
	rcb.Mat.Shiny = 500
}

func AddNewScene(parent ki.Ki, name string) *Scene {
	sc := parent.AddNewChild(KiT_Scene, name).(*Scene)
	sc.Defaults()
	return sc
}

func (sc *Scene) Render2D() {
	if sc.PushBounds() {
		if !sc.NoNav {
			sc.NavEvents()
		}
		if gi.Render2DTrace {
			fmt.Printf("3D Render2D: %v\n", sc.PathUnique())
		}
		sc.Render()
		sc.PopBounds()
	} else {
		sc.DisconnectAllEvents(gi.RegPri)
	}
}

func (sc *Scene) NavEvents() {
	// sc.ConnectEvent(oswin.MouseMoveEvent, gi.RegPri, func(recv, send ki.Ki, sig int64, d interface{}) {
	// 	me := d.(*mouse.MoveEvent)
	// 	me.SetProcessed()
	// 	ssc := recv.Embed(KiT_Scene).(*Scene)
	// 	orbDel := float32(.2)
	// 	panDel := float32(.05)
	//
	// 	del := me.Where.Sub(me.From)
	// 	dx := float32(del.X)
	// 	dy := float32(del.Y)
	// 	switch {
	// 	case key.HasAllModifierBits(me.Modifiers, key.Shift):
	// 		ssc.Camera.Pan(dx*panDel, -dy*panDel)
	// 	case key.HasAllModifierBits(me.Modifiers, key.Control):
	// 		ssc.Camera.PanAxis(dx*panDel, -dy*panDel)
	// 	case key.HasAllModifierBits(me.Modifiers, key.Alt):
	// 		ssc.Camera.PanTarget(dx*panDel, -dy*panDel, 0)
	// 	default:
	// 		if mat32.Abs(dx) > mat32.Abs(dy) {
	// 			dy = 0
	// 		} else {
	// 			dx = 0
	// 		}
	// 		ssc.Camera.Orbit(-dx*orbDel, -dy*orbDel)
	// 	}
	// 	ssc.UpdateSig()
	//
	// })
	sc.ConnectEvent(oswin.MouseScrollEvent, gi.RegPri, func(recv, send ki.Ki, sig int64, d interface{}) {
		me := d.(*mouse.ScrollEvent)
		me.SetProcessed()
		ssc := recv.Embed(KiT_Scene).(*Scene)
		if ssc.SetDragCursor {
			oswin.TheApp.Cursor(ssc.Viewport.Win.OSWin).Pop()
			ssc.SetDragCursor = false
		}
		zoom := float32(me.NonZeroDelta(false))
		zoomPct := float32(.05)
		zoomDel := float32(.05)
		switch {
		case key.HasAllModifierBits(me.Modifiers, key.Alt):
			ssc.Camera.PanTarget(0, 0, zoom*zoomDel)
		default:
			ssc.Camera.Zoom(zoomPct * zoom)
		}
		ssc.UpdateSig()
	})
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
	orbDeg := float32(5)
	panDel := float32(.1)
	zoomPct := float32(.05)
	switch ch {
	case "UpArrow":
		sc.Camera.Orbit(0, orbDeg)
		kt.SetProcessed()
	case "Shift+UpArrow":
		sc.Camera.Pan(0, panDel)
		kt.SetProcessed()
	case "Control+UpArrow":
		sc.Camera.PanAxis(0, panDel)
		kt.SetProcessed()
	case "Alt+UpArrow":
		sc.Camera.PanTarget(0, panDel, 0)
		kt.SetProcessed()
	case "DownArrow":
		sc.Camera.Orbit(0, -orbDeg)
		kt.SetProcessed()
	case "Shift+DownArrow":
		sc.Camera.Pan(0, -panDel)
		kt.SetProcessed()
	case "Control+DownArrow":
		sc.Camera.PanAxis(0, -panDel)
		kt.SetProcessed()
	case "Alt+DownArrow":
		sc.Camera.PanTarget(0, -panDel, 0)
		kt.SetProcessed()
	case "LeftArrow":
		sc.Camera.Orbit(orbDeg, 0)
		kt.SetProcessed()
	case "Shift+LeftArrow":
		sc.Camera.Pan(-panDel, 0)
		kt.SetProcessed()
	case "Control+LeftArrow":
		sc.Camera.PanAxis(-panDel, 0)
		kt.SetProcessed()
	case "Alt+LeftArrow":
		sc.Camera.PanTarget(-panDel, 0, 0)
		kt.SetProcessed()
	case "RightArrow":
		sc.Camera.Orbit(-orbDeg, 0)
		kt.SetProcessed()
	case "Shift+RightArrow":
		sc.Camera.Pan(panDel, 0)
		kt.SetProcessed()
	case "Control+RightArrow":
		sc.Camera.PanAxis(panDel, 0)
		kt.SetProcessed()
	case "Alt+RightArrow":
		sc.Camera.PanTarget(panDel, 0, 0)
		kt.SetProcessed()
	case "Alt++", "Alt+=":
		sc.Camera.PanTarget(0, 0, panDel)
		kt.SetProcessed()
	case "Alt+-", "Alt+_":
		sc.Camera.PanTarget(0, 0, -panDel)
		kt.SetProcessed()
	case "+", "=":
		sc.Camera.Zoom(-zoomPct)
		kt.SetProcessed()
	case "-", "_":
		sc.Camera.Zoom(zoomPct)
		kt.SetProcessed()
	case " ":
		sc.Camera.DefaultPose()
		kt.SetProcessed()
	case "w":
		sc.Camera.Zoom(-zoomPct)
		kt.SetProcessed()
	case "s":
		sc.Camera.Zoom(zoomPct)
		kt.SetProcessed()
	case "a":
		sc.Camera.Pan(panDel, 0)
		kt.SetProcessed()
	case "d":
		sc.Camera.Pan(-panDel, 0)
		kt.SetProcessed()
	case "t":
		kt.SetProcessed()
		obj := sc.Child(0).(*gi3d.Object)
		fmt.Printf("updated obj: %v\n", obj.PathUnique())
		obj.UpdateSig()
		return
	}
	sc.UpdateSig()
}
