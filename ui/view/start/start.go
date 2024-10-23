package start

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/869413421/wechatbot/ui"
	"github.com/869413421/wechatbot/ui/view/values"
	"image/color"
	"os"
)

type Impl struct {
	app     *ui.App
	window  *app.Window
	index   values.ViewIndexType
	showBtn bool
}

func New(
	app *ui.App,
	w *app.Window,
) values.View {
	svi := &Impl{
		app:    app,
		window: w,
		index:  values.ViewStart,
	}
	return svi
}

func (impl *Impl) ViewIndex() values.ViewIndexType {
	return impl.index
}

func (impl *Impl) Weight() values.WeightType {
	return 0
}

func (impl *Impl) Show(event event.Event) {
	var ops op.Ops
	e := event.(system.FrameEvent)
	th := material.NewTheme(gofont.Collection())
	gtx := layout.NewContext(&ops, e)

	l := material.H1(th, "start view")
	lb := material.Body1(th, "    this is start view,please design your content!")
	maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	l.Color = maroon
	l.Alignment = text.Middle
	lb.Color = maroon
	lb.Alignment = text.Start
	l.Layout(gtx)
	lb.Layout(gtx)

	e.Frame(gtx.Ops)
}

func (impl *Impl) OnEvent(e event.Event, data interface{}) {
	switch e := e.(type) {
	case system.DestroyEvent:
		os.Exit(0)
		return
	case system.FrameEvent:
		impl.Show(e)
	}
}

func (impl *Impl) OnLogin(body []byte) {
	impl.app.SetCurrentView(values.ViewStart)
}
