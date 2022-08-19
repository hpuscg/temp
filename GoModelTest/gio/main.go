package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("煮蛋计时器"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)

		var ops op.Ops

		var startButton widget.Clickable

		th := material.NewTheme(gofont.Collection())

		for e := range w.Events() {
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &startButton, "Start")
						return btn.Layout(gtx)
					},
					),
					layout.Rigid(
						layout.Spacer{Height: unit.Dp(25)}.Layout,
					),
				)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
