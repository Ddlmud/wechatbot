package ui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/869413421/wechatbot/ui/view/values"
	"github.com/eatmoreapple/openwechat"
	"image/color"
	"log"
	"os"
)

type App struct {
	currentView int
	window      *app.Window
	views       map[values.ViewIndexType]values.View
	data        interface{}
}

func New(w *app.Window) App {
	return App{
		window:      w,
		currentView: values.ViewMin,
		views:       make(map[values.ViewIndexType]values.View, 0),
	}
}

func (a *App) Views(index values.ViewIndexType) values.View {
	return a.views[index]
}

func (a *App) CurrentView() values.ViewIndexType {
	return a.currentView
}

func (a *App) SetCurrentView(currentView values.ViewIndexType) {
	a.currentView = currentView
}

func (a *App) RegisterView(index values.ViewIndexType, view values.View) {
	a.views[index] = view
}

func (a *App) OnLogin(body openwechat.CheckLoginResponse) {
	oldIndex := a.currentView
	for _, view := range a.views {
		view.OnLogin(body)
	}

	if values.ViewMin == oldIndex && a.currentView == values.ViewStart {
		go func() {
			if err := a.Loop(a.window); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
	}
}

func (a *App) ShowLoginQrCode(uuid string) {
	a.SetCurrentView(values.ViewLogin)
	a.data = openwechat.GetQrcodeUrl(uuid)

	go func() {
		if err := a.Loop(a.window); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

func (a *App) Loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	for {
		e := <-w.Events()
		view, ok := a.views[a.currentView]
		if len(a.views) > 0 && ok {
			view.OnEvent(e, a.data)
			continue
		}

		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			l := material.H1(th, "wechat robot")
			lb := material.Body1(th, "    this is no view,please check your code!")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			l.Color = maroon
			l.Alignment = text.Middle
			lb.Color = maroon
			lb.Alignment = text.Start
			l.Layout(gtx)
			lb.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
