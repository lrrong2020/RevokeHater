package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
)

type QueueItem struct {
	SenderNickName string
	MessageCreateTime int64
	MessageID string
	MessageContent string
}

type Queue struct {
	items []QueueItem
	mu    sync.Mutex
}

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式


func (q *Queue) Add(item QueueItem) {
    q.mu.Lock()
    defer q.mu.Unlock()

    q.items = append(q.items, item)
    if len(q.items) > 100 {
        // Discard the oldest item.
        q.items = q.items[1:]
    }
}

// Initialize the queue.
queue := &Queue{}

bot.MessageHandler = func(msg *openwechat.Message) {
    item := QueueItem{
        SenderNickName: msg.NickName, 
        MessageCreateTime: msg.CreateTime,
        MessageID: msg.MsgId,
        MessageContent: msg.Content,
    }

    queue.Add(item)
		fmt.Printf("%+v\n", *item)
}



    if msg.IsRecalled() {
        revokeMsg, err := msg.RevokeMsg()
		if err != nil{
            fmt.Println(err)
            return
  	}
		fmt.Println("Replace: ", revokeMsg.RevokeMsg.ReplaceMsg)

        // Here you can access the fields of revokeMsg and reply accordingly
        // For example, you can reply with the content of the recalled message
  	// msg.ReplyText(fmt.Sprintf("You've recalled a message with ID: %d", revokeMsg.RevokeMsg.ReplaceMsg))
		// msg.ReplyText(fmt.Sprintf("UserID: %d", userID))	
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
