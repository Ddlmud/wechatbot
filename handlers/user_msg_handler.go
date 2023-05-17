package handlers

import (
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)
var ServerStartTime int64 = 0

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
	gpt     gtp.GPTChat
	history []string
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() && msg.CreateTime > ServerStartTime {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	chat := gtp.New()
	return &UserMessageHandler{gpt: chat, history: make([]string, 0)}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)
	if /*msg.IsSendBySelf() || */ !strings.Contains(msg.Content, "[旺柴][旺柴]") {
		return nil
	}
	if strings.Contains(msg.Content, "[旺柴][旺柴]新会话") {
		g.history = make([]string, 0)
	}

	// 向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")
	requestText = strings.ReplaceAll(msg.Content, "[旺柴][旺柴]", "")
	if requestText == "" {
		return nil
	}
	g.history = append(g.history, requestText)
	reply, err := g.gpt.Completion(g.history...)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
