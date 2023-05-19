package handlers

import (
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
	gpt     gtp.GPTChat
	history map[string][]string
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() && msg.CreateTime > ServerStartTime {
		return g.ReplyText(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	chat := gtp.New()
	return &GroupMessageHandler{gpt: chat, history: make(map[string][]string, 0)}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收群消息
	sender, err := msg.Sender()
	group := openwechat.Group{sender}
	log.Printf("Received Group %v Text Msg : %v", group.NickName, msg.Content)

	// 不是@的不处理
	if !msg.IsAt() || !strings.Contains(msg.Content, "[旺柴][旺柴]") {
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

	// 替换掉@文本，然后向GPT发起请求
	replaceText := "@" + sender.Self.NickName
	requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	requestText = strings.ReplaceAll(msg.Content, "[旺财][旺财]", "")
	if requestText == "" {
		return nil
	}
	g.history[sender.NickName] = append(g.history[sender.NickName], requestText)
	reply, err := g.gpt.Completion(g.history[sender.NickName]...)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		log.Printf("get sender in group error :%v \n", err)
		return err
	}

	// 回复@我的用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	atText := "@" + groupSender.NickName
	replyText := atText + " " + reply
	_, err = msg.ReplyText(replyText)
	if err != nil {
		log.Printf("response group error: %v \n", err)
	}
	return err
}
