package main

import (
	"fmt"

	"./raft"
)

func main() {

	c := raft.Client{}
	c.ECDSAgenerate()

	go c.Running()

	for {
		operate := ""
		operand := 0

		fmt.Print("请输入操作数:")
		fmt.Scan(&operand)
		fmt.Print("请输入操作:")
		fmt.Scan(&operate)

		c.SendComment(operate, operand)
	}
}
