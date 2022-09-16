package raft

type Message struct {
	Id    int       //表示消息发送者身份
	Req   string    //comment、done、commit、already_commit、apply、heartbeat、askvote、vote
	Entry Log_entry //entry同步的时候用
	Term  int       //leader发送心跳的时候用
	Com   Comment   //comment的时候用
	Key   []byte    //发送者的公钥
	Sign  []byte    //数字签名

}

func commitMessage(id int, e Log_entry, key []byte) Message {
	m := Message{
		Id: id, Key: key,
		Req:   "commit",
		Entry: e,
	}
	return m
}

func already_commitMessage(id int, key []byte) Message {
	m := Message{
		Id:  id,
		Key: key,
		Req: "already_commit",
	}
	return m
}

func applyMessage(id int, e Log_entry, key []byte) Message {
	m := Message{
		Id:    id,
		Key:   key,
		Req:   "apply",
		Entry: e,
	}
	return m
}

func doneMessage(id int, key []byte) Message {
	m := Message{
		Id:  id,
		Key: key,
		Req: "done",
	}
	return m
}

func commentMessage(c Comment, key []byte) Message {
	m := Message{
		Key: key,
		Req: "comment",
		Com: c,
	}
	// fmt.Println(nil == m.Sign)
	return m
}

func heartbeatMessage(id int, term int, key []byte) Message {
	m := Message{
		Id:   id,
		Key:  key,
		Term: term,
		Req:  "heartbeat",
		// Beat: b,
	}
	return m
}

func askvoteMessage(id int, term int, key []byte) Message {
	m := Message{
		Id:   id,
		Key:  key,
		Term: term,
		Req:  "askvote",
	}
	return m
}

func voteMessage(id int, key []byte) Message {
	m := Message{
		Id:  id,
		Key: key,
		Req: "vote",
	}
	return m
}
