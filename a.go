package main

import (
	"fmt"
	"sync"
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

func (q *Queue) Add(item QueueItem) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
	if len(q.items) > 100 {
		// Discard the oldest item.
		q.items = q.items[1:]
	}
}

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// Initialize the queue.
	queue := &Queue{}

	bot.MessageHandler = func(msg *openwechat.Message) {
		sender, err2 := msg.Sender()
		userID := sender.ID()
		if err2 != nil{
			fmt.Println(err2)
			return
}
		fmt.Println("User ID:", userID)

		item := QueueItem{
			SenderNickName: sender.NickName, 
			MessageCreateTime: msg.CreateTime,
			MessageID: msg.MsgId,
			MessageContent: msg.Content,
		}

		queue.Add(item)
		fmt.Printf("%+v\n", item)
	}

	// Rest of your code...


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
