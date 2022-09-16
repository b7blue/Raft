package raft

import (
	"fmt"
	"math/rand"
	"time"
)

// 发送选票
func (n *Node) handleAskvote(mess Message) {
	// 决定是否要投票：1、一个term内只能投给一个人 2、先到先得
	// 1、两个要求，假如自己是candidate就不投票、假如已经投给别的同term节点就不投了
	if n.Status != "candidate" && n.Votedone[mess.Term] == 0 {
		n.setVotedone(mess.Term)
		m := voteMessage(n.Id, n.PublicKey)
		m.ECDSAsign(n.PrivateKey)
		byte_m := packMessage(m)
		dst_addr := node_table[mess.Id-1]
		sendMessage(byte_m, dst_addr)
	}
}

// 接收选票
func (n *Node) handleVote(mess Message) {
	n.addVote()
	fmt.Printf("本节点收到节点<%d>投的一票\n\n", mess.Id)
}

//  接收心跳
func (n *Node) handeleHeartbeat(mess Message) {
	// 判断发送心跳的节点的term，只有接受的term >= 本身的term才有用
	// fmt.Println("接收心跳，未判断是否有效")
	if mess.Term >= n.Term {
		// 通知心跳定时器接收到了心跳
		n.Heartbeat <- 1
		fmt.Printf("本节点收到了来自leader<%d>的一次心跳,此时Term为%d\n\n", mess.Id, mess.Term)

		if mess.Term == n.Term {
			// 2. = (1)candidate收到了新leader的心跳 (2)follower收到leader正常的心跳
			if n.Status == "candidate" {
				// 停止选举
				n.setStatus("follower")
			}
		} else {
			// 1. > follower、leader、candidate收到了新leader的心跳
			n.setTerm(mess.Term)

			switch n.Status {
			case "leader":
				// 关闭定时发送心跳
				n.setStatus("follower")
			case "candidate":
				// 停止选举
				n.setStatus("follower")
			}
		}
	}
}

// 发起选举
func (n *Node) startElection() {
	// term + 1, change status to "candidate"
	n.addTerm()
	n.setStatus("candidate")
	n.initVote()

	fmt.Printf("没有按时收到leader的心跳,本节点成为candidate,发起选举,此时term为%d\n\n", n.Term)

	// broadcast askvote message
	m := askvoteMessage(n.Id, n.Term, n.PublicKey)
	m.ECDSAsign(n.PrivateKey)
	bytes_m := packMessage(m)
	dst_addrs := n.BroadcastAddr()
	sendMessage(bytes_m, dst_addrs...)

	// listening, recv vote......
	ticker := time.NewTicker(100 * time.Millisecond)
	i := 0
	// 为了防止平票，设置随机选举定时器
	rand.Seed(time.Now().Unix())
	time := rand.Intn(10) + 1
	for i < time {

		<-ticker.C
		i++

		// 假如获得半数以上选票，就成为leader
		if n.Vote >= node_num/2 {
			fmt.Printf("本节点成为新leader,停止选举!此时term为%d\n\n", n.Term)
			n.setStatus("leader")
			go n.BroadcastHeartbeat()
			go n.talk2client()
			break
		}
		// 假如收到新的leader的心跳(也就说成为follower)，停止选举
		if n.Status == "follower" {
			fmt.Printf("新leader已经出现,停止选举!\n\n")
			break
		}
	}

	// 说明选举定时器时间已到
	if i == 10 {
		// 开启新的选举
		n.startElection()
	}
}

// 广播心跳
func (n *Node) BroadcastHeartbeat() {
	ticker := time.NewTicker(300 * time.Millisecond)
	for {
		// 检查是否是leader,否则不广播心跳或者停止广播
		if n.Status != "leader" {
			break
		}

		select {
		case <-ticker.C:
			// package heartbeat message
			m := heartbeatMessage(n.Id, n.Term, n.PublicKey)
			m.ECDSAsign(n.PrivateKey)
			bytes_m := packMessage(m)
			dst_addrs := n.BroadcastAddr()
			go sendMessage(bytes_m, dst_addrs...)
			fmt.Printf("leader<%d>广播了一次心跳\n\n", n.Id)
		default:
		}
	}
}

// 心跳定时器
func (n *Node) HeartbeatTicker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		reset := false
		for i := 0; i < 30; i++ {
			select {
			case <-n.Heartbeat:
				// 假如收到心跳,就跳出内层循环,定时器重新定时一次.
				reset = true
				break
			case <-ticker.C:
				// 什么也不做
			}
		}

		// 假如300ms已到,follower还没收到心跳,就要转为候选人状态,开启选举
		if !reset && n.Status != "leader" {
			n.startElection()
		}
	}

}
