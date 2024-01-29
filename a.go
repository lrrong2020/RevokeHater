package main

import (
	"fmt"
	"sync"
	"time"
	"strconv"
	"os"
	"github.com/eatmoreapple/openwechat"
)

type QueueItem struct {
	SenderNickName string
	Alias string
	MessageCreateTime string
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
	if len(q.items) > 500 {
		// Discard the oldest item.
		q.items = q.items[1:]
	}
}

func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func (q *Queue) FindByID(id string) (QueueItem, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, item := range q.items {
		if item.MessageID == id {
			return item, true
		}
	}
	return QueueItem{}, false
}


func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// Initialize the queue.
	queue := &Queue{}
	bot.MessageHandler = func(msg *openwechat.Message) {
		//debug
		fmt.Printf("[Message Info]%+v\n\n", *msg)

		sender, err2 := msg.SenderInGroup()
		if err2 != nil{
			fmt.Println(err2)
			return
		}
		// fmt.Printf("%+v\n\n", *sender)

		// fmt.Println("User ID:", userID)


		// if (msg.IsText() && strings.HasPrefix(msg.Content, ":-::---:::") ){
		// 	toUserName = msg.ToUserName
		// 	fmt.Println("Update ToUserName: ", toUserName)

		// 	var group *openwechat.User
		// 	var err3 error
		// 	if msg.IsSendByGroup() {
		// 		if msg.IsSendBySelf() {
		// 				group, err3 = msg.Receiver()
		// 		} else {
		// 			group, err3 = msg.Sender()
		// 		}
		// 	}
		// 	groupNickName := group.NickName
		// 	toUserName
		// }
		
		// fmt.Println("\n\n\ntoUserName: ", toUserName)
		// fmt.Println("\n\n\nmsg.ToUserName: ", msg.ToUserName)


		if (msg.IsText() && msg.IsSendByGroup()) {
			i, err := strconv.ParseInt(strconv.FormatInt(msg.CreateTime, 10), 10, 64)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}
			tm := time.Unix(i, 0)
	    	// Load the desired timezone
			location, err := time.LoadLocation("Asia/Shanghai")  // Asia/Shanghai is UTC+8
			if err != nil {
				fmt.Println("Error loading location:", err)
				return
			}
	
			// Convert the time value to the desired timezone
			tm = tm.In(location)
	
			// fmt.Println(tm)

			item := QueueItem{
				SenderNickName: sender.NickName,
				Alias: sender.Alias, 
				MessageCreateTime: tm.Format("1月2日 15:04"),
				MessageID: msg.MsgId,
				MessageContent: msg.Content,
			}
	
			queue.Add(item)
			// fmt.Printf("%+v\n", item)
			// size := queue.Size()
			// fmt.Println("Size of the queue:", size)
		}

		if msg.IsRecalled() {
			revokeMsg, err := msg.RevokeMsg()
			if err != nil {
				fmt.Println(err)
				return
			}
			recalledMsgId := strconv.FormatInt(revokeMsg.RevokeMsg.MsgId, 10)
			// Find the recalled message in the queue.
			if item, found := queue.FindByID(recalledMsgId); found {
				resString := "[" + item.MessageCreateTime + "] " + item.Alias + "(" + item.SenderNickName + "): " + item.MessageContent
		
				// Open the file in append mode, create it if it does not exist.
				f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer f.Close()
		
				// Write the result string to the file.
				if _, err := f.WriteString(resString + "\n"); err != nil {
					fmt.Println(err)
					return
				}
			} else {
				fmt.Println("The recalled message was not found in the queue.")
			}
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
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(friends)

	// 获取所有的群组
	groups, err := self.Groups()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(groups)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
