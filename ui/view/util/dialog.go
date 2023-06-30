package util

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func ShowDialog(text, title string) {
	w := app.NewWindow(
		app.Title(title),
		app.Size(unit.Dp(400), unit.Dp(300)),
	)
	dlg := new(widget.Clickable)
	var ops op.Ops
	th := material.NewTheme(gofont.Collection())
	go func() {
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.DestroyEvent:
				return
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				for _, gtxEvent := range gtx.Events(0) {
					switch gtxE := gtxEvent.(type) {
					case pointer.Event:
						if gtxE.Type == pointer.Press && dlg.Clicked() {
							// 当弹窗关闭按钮被点击时，关闭弹窗
							return
						}
					default:
						layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							// 显示弹窗内容
							return material.Button(th, dlg, text).Layout(gtx)
						})
						e.Frame(gtx.Ops)
					}
				}
			case pointer.Event:
				if e.Type == pointer.Press && dlg.Clicked() {
					// 当弹窗关闭按钮被点击时，关闭弹窗
					return
				}
			}
		}
	}()

}
