package raft

import (
	"fmt"
)

// 处理客户端的请求
func (n *Node) handleComment(mess Message) {
	n.initAlready()

	fmt.Printf("leader< %d >收到client的请求%+v，开始同步消息\n", n.Id, mess.Com)

	// 	打包指令成log entry，
	entry := n.PackEntry(mess.Com)

	// 	append log entry（commit），log的term更新、index+1，
	n.Commit(entry)

	// 	广播AppendEntries RPC，(广播)
	m := commitMessage(n.Id, entry, n.PublicKey) //构造commit消息
	m.ECDSAsign(n.PrivateKey)
	byte_m := packMessage(m)
	dst_addrs := n.BroadcastAddr() //获取广播的地址
	sendMessage(byte_m, dst_addrs...)

	// 	接收已commit消息，等消息数到节点数的一半时，（监听）
	for {
		//循环监听消息

		fmt.Println("n.Already_commit", n.Already_commit)
		if n.Already_commit >= node_num/2 {
			// 	log entry提交给复制状态机（apply），
			n.RCM(entry)

			// 	将结果返回client，（发送）
			m = doneMessage(n.Id, n.PublicKey)
			m.ECDSAsign(n.PrivateKey)
			byte_m = packMessage(m)
			dst_addr := client_addr
			sendMessage(byte_m, dst_addr)

			// 	广播apply log entry的消息。（广播）
			m = applyMessage(n.Id, entry, n.PublicKey)
			m.ECDSAsign(n.PrivateKey)
			byte_m = packMessage(m)
			dst_addrs = n.BroadcastAddr() //获取广播的地址
			sendMessage(byte_m, dst_addrs...)

			fmt.Println("leader<", n.Id, ">同步消息已完成，此时leader的本地内容如下:")
			fmt.Printf("%+v\n", n)

			break
		}
	}

}

// 处理apply要求
func (n *Node) handleApply(mess Message) {
	n.RCM(mess.Entry)

	fmt.Println("节点<", n.Id, ">收到来自leader<", mess.Id, ">的apply要求, apply后节点本地内容如下：")
	fmt.Printf("%+v\n", n)
}

// 处理commit要求
func (n *Node) handleCommit(mess Message) {

	n.Commit(mess.Entry)

	// 回复already commit给leader
	m := already_commitMessage(n.Id, n.PublicKey)
	m.ECDSAsign(n.PrivateKey)
	byte_m := packMessage(m)
	dst_addr := node_table[mess.Id-1] //获取目标地址
	sendMessage(byte_m, dst_addr)

	fmt.Println("follower<", n.Id, ">收到来自leader<", mess.Id, ">的commit要求，并且返回响应，commit后节点本地内容如下：")
	fmt.Printf("%+v\n", n)
}

// 接收到follower已经commit的消息
func (n *Node) handleAlready(mess Message) {
	n.addAlready()

	fmt.Println("leader<", n.Id, ">收到来自follower<", mess.Id, ">的already commit回复")
}
