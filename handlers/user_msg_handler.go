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
	history map[string][]string
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
	return &UserMessageHandler{gpt: chat, history: make(map[string][]string, 0)}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)
	if /*msg.IsSendBySelf() || */ !strings.Contains(msg.Content, "[旺柴][旺柴]") {
		return nil
	}

	// 处理下上下文历史记录
	_, ok := g.history[sender.NickName]
	if !ok {
		g.history[sender.NickName] = make([]string, 0)
	}
	// 如果是新会话，则清空历史记录
	if strings.Contains(msg.Content, "[旺柴][旺柴]新会话") {
		g.history[sender.NickName] = make([]string, 0)
	}

	// 向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")
	requestText = strings.ReplaceAll(msg.Content, "[旺柴][旺柴]", "")
	if requestText == "" {
		return nil
	}
	g.history[sender.NickName] = append(g.history[sender.NickName], requestText)
	reply, err := g.gpt.Completion(g.history[sender.NickName]...)
	if err != nil {
		log.Printf("gpt request error: %v \n", err)
		msg.ReplyText("bug了，我一会发现了就去修。")
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
