package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		sender, err2 := msg.Sender()
		userID := sender.ID()
		fmt.Println("User ID:", userID)
    if msg.IsRecalled() {
        revokeMsg, err := msg.RevokeMsg()
		if err != nil{
            fmt.Println(err)
            return
  	}
        // Here you can access the fields of revokeMsg and reply accordingly
        // For example, you can reply with the content of the recalled message
  	msg.ReplyText(fmt.Sprintf("You've recalled a message with ID: %d", revokeMsg.RevokeMsg.ReplaceMsg))
		msg.ReplyText(fmt.Sprintf("UserID: %d", userID))	
    } else if msg.IsText() && msg.Content == "ping" {
        msg.ReplyText("pong")
    }
}

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
