package raft

type Comment struct {
	Operate string
	Operand int
}

type Log_entry Comment

type Log struct {
	Log_entries []Log_entry //结构体数组初始化需要注意
	// Term        int
	Index int
}
