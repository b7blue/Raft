package main

import (
	"fmt"

	"./raft"
)

func main() {

	id := 0
	fmt.Print("请输入节点的编号：")
	fmt.Scan(&id)

	n := raft.Node{
		Id:        id,
		Status:    "follower",
		Heartbeat: make(chan int),
		Votedone:  make(map[int]int),
	}
	// fmt.Printf("节点<%d>的初始状态如下:\n%+v\n", n.Id, n)
	n.ECDSAgenerate()
	// fmt.Printf("节点<%d>的初始状态如下:\n%+v\n", n.Id, n)

	n.Register()

}
