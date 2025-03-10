package bootstrap

import (
	"github.com/869413421/wechatbot/handlers"
	"github.com/869413421/wechatbot/ui"
	"github.com/eatmoreapple/openwechat"
	"log"
	"time"
)

func Run(app *ui.App) {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	handlers.ServerStartTime = time.Now().Unix()

	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	//bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	bot.UUIDCallback = app.ShowLoginQrCode
	bot.OnLogin(app.OnLogin)

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	} else {
		bot.LoginCallBack(nil)
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	//bot.Block()
}
