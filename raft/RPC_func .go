package raft

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

func sendMessage(Message []byte, dst_addrs ...string) {
	// 变参函数，一个个向目标地址发消息
	var mu sync.Mutex
	num := 0
	sent := make(chan int, len(dst_addrs))
	for _, addr := range dst_addrs {
		// 多线程发送，快
		go func(addr string) {
			conn, err := net.Dial("tcp", "localhost:"+addr)

			if err != nil {
				// 说明目的Node已下线
			} else {
				mu.Lock()
				num++
				mu.Unlock()
				// 发送数据

				conn.Write(Message)
				// fmt.Printf("向%s发送数据%s\n", conn.RemoteAddr(), string(Message))

				conn.Close()

			}
			sent <- 1

		}(addr)
	}

	for i := 0; i < len(dst_addrs); i++ {
		<-sent
	}
	node_num = num + 1
}

func packMessage(m Message) []byte {

	jsonbytes, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return jsonbytes
}

func (n *Node) BroadcastAddr() []string {
	dst_addr := make([]string, 0)
	my_id := n.Id
	for i := 0; i < init_num; i++ {
		if my_id != i+1 {
			dst_addr = append(dst_addr, node_table[i])
		}
	}
	return dst_addr
}
