package raft

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var t = [3]string{"9000", "9001", "9002"}
var node_table = t[:]
var node_num int
var init_num int = 3
var client_addr string = "9003"
var server_addr string = "9004"

func (n *Node) Register() {

	// 开启监听，说明在线
	listener, err := net.Listen("tcp", "localhost:"+node_table[n.Id-1])
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("节点<", n.Id, ">上线。。。")

	// 监听心跳
	go n.HeartbeatTicker()
	// 广播心跳
	go n.BroadcastHeartbeat()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			fmt.Println(err)
		} else {
			go handleMessage(n, conn)
		}
		// fmt.Println("A connection is closed")
	}

}

func handleMessage(n *Node, c net.Conn) {
	input := bufio.NewScanner(c)

	for input.Scan() {
		jsonbytes := input.Bytes()
		var m Message

		// 用json格式传输，所以收到之后先转换回Message
		json.Unmarshal(jsonbytes, &m)
		// 要验证一下数字签名
		if m.ECDSAverify() {

			// fmt.Printf("收到一条消息,消息内容如下:\n%+v\n", m)
			req := m.Req
			switch req {
			case "comment":
				n.handleComment(m)
			case "apply":
				n.handleApply(m)
			case "commit":
				n.handleCommit(m)
			case "already_commit":
				n.handleAlready(m)
			case "heartbeat":
				n.handeleHeartbeat(m)
			case "askvote":
				n.handleAskvote(m)
			case "vote":
				n.handleVote(m)
			}

			fmt.Println()
		} else {
			fmt.Println("签名验证失败，消息不安全")
		}
	}
}

func (n *Node) talk2client() {
	// 开启与client通信的端口，说明在线
	listener, err := net.Listen("tcp", "localhost:"+server_addr)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("server开启\n\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			fmt.Println(err)
		} else {

			go handleMessage(n, conn)
		}
	}
}
