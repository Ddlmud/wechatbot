package login

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/869413421/wechatbot/ui"
	"github.com/869413421/wechatbot/ui/view/values"
	"image"
	"net/http"
)

type LoginView interface {
	ShowLoginQrCode(uuid string)
}

type Impl struct {
	app    *ui.App
	window *app.Window
	index  values.ViewIndexType
}

func New(
	app *ui.App,
	w *app.Window,
) values.View {
	svi := &Impl{
		app:    app,
		window: w,
		index:  values.ViewLogin,
	}
	return svi
}

func (impl *Impl) OnLogin(body []byte) {

}

func (impl *Impl) OnEvent(e event.Event, data interface{}) {
	// 图片URL
	imageURL := data.(string)

	// 创建一个图片Widget
	img := &widget.Image{
		Fit:   widget.Fill,
		Scale: 1.0,
	}

	var ops op.Ops

	switch e := e.(type) {
	case system.DestroyEvent:
	case system.FrameEvent:
		// 加载并显示图片
		if err := loadImage(imageURL, img); err != nil {
			// 处理加载图片错误
		}

		gtx := layout.NewContext(&ops, e)
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return img.Layout(gtx)
		})
		e.Frame(gtx.Ops)
	}
}

func (impl *Impl) Weight() values.WeightType {
	return 0
}

func (impl *Impl) ViewIndex() values.ViewIndexType {
	return impl.index
}

func loadImage(url string, img *widget.Image) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imgData, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}

	img.Src = paint.NewImageOp(imgData)

	return nil
}
