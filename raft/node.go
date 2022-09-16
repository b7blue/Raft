package raft

import (
	"fmt"
	"sync"
)

type Node struct {
	Term   int
	Log    Log
	Result int
	mu     sync.Mutex //为了互斥访问，加锁！

	Id             int // 表示节点的编号
	Status         string
	Vote           int
	Already_commit int
	Heartbeat      chan int
	Votedone       map[int]int
	PrivateKey     []byte
	PublicKey      []byte
}

// apply, commit Log entry to Replicated state machines
func (n *Node) RCM(comment Log_entry) {
	Result := 0
	switch comment.Operate {
	case "+":
		Result = comment.Operand + 1
	case "-":
		Result = comment.Operand - 1
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Result = Result
}

// pack client's comment to Log entry, return
func (n *Node) PackEntry(c Comment) Log_entry {
	entry := Log_entry(c)
	return entry
}

// commit, append Log entries
func (n *Node) Commit(entry Log_entry) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Log.Log_entries = append(n.Log.Log_entries, entry)
	n.Log.Index++
	fmt.Printf("已commit,此时节点状态为:\n%+v\n\n", n)
}

func (n *Node) setTerm(t int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Term = t
}

func (n *Node) addTerm() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Term++
}

func (n *Node) setStatus(s string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Status = s
}

func (n *Node) addAlready() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Already_commit++
}

func (n *Node) initAlready() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Already_commit = 0
}

func (n *Node) addVote() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Vote++
}

func (n *Node) initVote() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Vote = 0
}

func (n *Node) setVotedone(term int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Votedone[term]++
}
