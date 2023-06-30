package run

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/869413421/wechatbot/bootstrap"
	"github.com/869413421/wechatbot/ui"
	"github.com/869413421/wechatbot/ui/view/login"
	"github.com/869413421/wechatbot/ui/view/start"
	"github.com/869413421/wechatbot/ui/view/values"
)

func Run() {
	go func() {
		window := app.NewWindow(
			app.Title("wechat robot"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)
		weApp := ui.New(window)

		weApp.RegisterView(values.ViewStart, start.New(&weApp, window))
		weApp.RegisterView(values.ViewLogin, login.New(&weApp, window))
		bootstrap.Run(&weApp)
	}()
	app.Main()
}
