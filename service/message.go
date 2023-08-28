package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"simple-demo/controller"
	"sync"
)

// chatConnMap用于存储用户之间的聊天连接，键为聊天对方的用户ID，值为net.Conn类型的连接对象。
var chatConnMap = sync.Map{}

// RunMessageServer 函数用于启动消息服务器，监听TCP连接，并处理聊天消息。
func RunMessageServer() {
	listener, err := net.Listen("tcp", "112.193.83.145:8080")
	if err != nil {
		fmt.Printf("Run message server failed: %v\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept conn failed: %v\n", err)
			continue
		}

		go process(conn)
	}
}

// process函数处理单个客户端的聊天消息。
func process(conn net.Conn) {
	defer conn.Close()

	var buf [256]byte
	for {
		n, err := conn.Read(buf[:])
		if n == 0 {
			if err == io.EOF {
				// 客户端连接关闭时，结束处理该连接的聊天消息。
				break
			}
			fmt.Printf("Read message failed: %v\n", err)
			continue
		}

		var event = controller.MessageSendEvent{}
		_ = json.Unmarshal(buf[:n], &event)
		fmt.Printf("Receive Message：%+v\n", event)

		// 生成用户之间的聊天键，格式为“发送者ID_接收者ID”。
		fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)

		// 如果消息内容为空，表示客户端发送的是连接消息，将客户端连接存储在chatConnMap中，以便后续消息推送。
		if len(event.MsgContent) == 0 {
			chatConnMap.Store(fromChatKey, conn)
			continue
		}

		// 如果消息内容不为空，表示客户端发送的是聊天消息。

		// 生成接收者的聊天键，格式为“接收者ID_发送者ID”。
		toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)

		// 根据接收者的聊天键查找对方用户的连接。
		writeConn, exist := chatConnMap.Load(toChatKey)

		// 如果对方用户不在线（即找不到对方用户的连接），则打印一条提示信息，然后继续处理下一条消息。
		if !exist {
			fmt.Printf("User %d offline\n", event.ToUserId)
			continue
		}

		// 如果对方用户在线，则根据接收到的消息内容构造推送事件，并通过writeConn将推送事件发送给对方用户，实现消息推送功能。

		pushEvent := controller.MessagePushEvent{
			FromUserId: event.UserId,
			MsgContent: event.MsgContent,
		}

		pushData, _ := json.Marshal(pushEvent)
		_, err = writeConn.(net.Conn).Write(pushData)
		if err != nil {
			fmt.Printf("Push message failed: %v\n", err)
		}
	}
}
