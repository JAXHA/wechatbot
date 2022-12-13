package bootstrap

import (
	"fmt"

	"github.com/869413421/wechatbot/handlers"
	"github.com/869413421/wechatbot/pkg/logger"
	"github.com/eatmoreapple/openwechat"
)

func Run() {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	handler, err := handlers.NewHandler()
	if err != nil {
		logger.Danger("register error: %v", err)
		return
	}
	bot.MessageHandler = handler

	// 注册登陆二维码回调
	bot.UUIDCallback = handlers.QrCodeCallBack

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录
	err = bot.HotLogin(reloadStorage, true)
	if err != nil {
		logger.Warning(fmt.Sprintf("login error: %v ", err))
		return
	}
	//GetCurrentUser 获取当前的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		return
	}
	fmt.Println(self.NickName)

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println("好友列表:", friends, err)
	fmt.Println("好友的数量:", friends.Count())

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println("群列表:", groups, err)

	// Mps 获取所有的公众号
	mps, err := self.Mps()
	fmt.Println("公众号列表:", mps, err)

	//给文件助手发消息
	fh := openwechat.NewFriendHelper(self)
	fmt.Println(fh)
	self.SendTextToFriend(fh, "机器人开始运行")

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
