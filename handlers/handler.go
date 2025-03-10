package handlers

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/pkg/logger"
	"github.com/eatmoreapple/openwechat"
	"github.com/patrickmn/go-cache"
	"github.com/skip2/go-qrcode"
)

var c = cache.New(config.LoadConfig().SessionTimeout, time.Minute*5)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle() error
	ReplyText() error
}

// QrCodeCallBack 登录扫码回调，
func QrCodeCallBack(uuid string) {
	if runtime.GOOS == "windows" {
		// 运行在Windows系统上
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		log.Println("login in linux")
		url := "https://login.weixin.qq.com/l/" + uuid
		log.Printf("如果二维码无法扫描，请缩小控制台尺寸，或更换命令行工具，缩小二维码像素")
		err := qrcode.WriteFile(url, qrcode.Medium, 256, "qr.png")
		fmt.Println(err)
		q, _ := qrcode.New(url, qrcode.High)
		fmt.Println(q.ToSmallString(true))
	}
}

func NewHandler() (msgFunc func(msg *openwechat.Message), err error) {
	//基于这个回调函数，可以对消息进行多样化处理
	dispatcher := openwechat.NewMessageMatchDispatcher()

	// 清空会话 // 注册消息处理函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return strings.Contains(message.Content, config.LoadConfig().SessionClearToken)
	}, TokenMessageContextHandler())

	// 处理群消息 // 注册消息处理函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsSendByGroup()
	}, GroupMessageContextHandler())

	// 好友申请
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsFriendAdd()
	}, func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		if config.LoadConfig().AutoPass {
			_, err := msg.Agree("")
			if err != nil {
				logger.Warning(fmt.Sprintf("add friend agree error : %v", err))
				return
			}
		}
	})

	// 私聊
	// 获取用户消息处理器
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return !(strings.Contains(message.Content, config.LoadConfig().SessionClearToken) || message.IsSendByGroup() || message.IsFriendAdd())
	}, UserMessageContextHandler())

	// 返回回调函数
	return openwechat.DispatchMessage(dispatcher), nil
}
