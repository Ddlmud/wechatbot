package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"log"
	"os"
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
	if msg.IsText() && msg.CreateTime > ServerStartTime /*&& !msg.IsSendBySelf()*/ {
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
	if err != nil {
		return err
	}
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)

	// 读配置
	cfg := config.LoadConfig()
	question := "default"
	answer, ok := cfg.AnswerMapStr[msg.Content]
	if ok {
		question = "text"
		return ReplyMessage(msg, question, answer)
	}
	answer, ok = cfg.AnswerMapPic[msg.Content]
	if ok {
		question = "pic"
		return ReplyMessage(msg, question, answer)
	}
	answer, ok = cfg.AnswerMapVideo[msg.Content]
	if ok {
		question = "video"
		return ReplyMessage(msg, question, answer)
	}
	return ReplyMessage(msg, question, cfg.QuestionAndAnswer)
}

func ReplyMessage(msg *openwechat.Message, question, answer string) error {
	switch question {
	case "pic":
		file, err := os.OpenFile(answer, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			msg.ReplyText("bug了，正在紧急修复中。")
			return err
		}
		_, err = msg.ReplyImage(file)
		if err != nil {
			msg.ReplyText("bug了，正在紧急修复中。")
			return err
		}
		return err
	case "video":
		file, err := os.OpenFile(answer, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			msg.ReplyText("bug了，正在紧急修复中。")
			return err
		}
		_, err = msg.ReplyVideo(file)
		if err != nil {
			msg.ReplyText("bug了，正在紧急修复中。")
			return err
		}
		return err
	case "text":
		_, err := msg.ReplyText(answer)
		if err != nil {
			msg.ReplyText("bug了，正在紧急修复中。")
			return err
		}
		return err
	default:
		_, err := msg.ReplyText(answer)
		if err != nil {
			log.Printf("response user error: %v \n", err)
		}
		return err
	}
}
