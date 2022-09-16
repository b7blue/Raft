package raft

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Client struct {
	PrivateKey []byte
	PublicKey  []byte
}

// 生成公私钥,直接给node加上
func (c *Client) ECDSAgenerate() {
	// 生成私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		println(err)
	}

	// 将私钥序列化为PKIX格式DER编码
	privateBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		println(err)
	}

	// 取出公钥
	publicKey := privateKey.PublicKey

	// 将公钥序列化为PKIX格式DER编码
	publicBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		println(err)
	}

	c.PublicKey = publicBytes
	c.PrivateKey = privateBytes

}

func (c *Client) Running() {

	// 开启监听，说明在线
	listener, err := net.Listen("tcp", "localhost:"+client_addr)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("客户端已上线...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			fmt.Println(err)
		} else {
			go talk2server(conn)
		}
	}

}

func talk2server(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		b := input.Bytes()
		var m Message
		// 要验证一下数字签名
		// 用json格式传输，所以收到之后先转换回Message
		json.Unmarshal(b, &m)

		switch m.Req {
		case "done":
			fmt.Println("good job!")
		}
	}

}

func (c *Client) SendComment(operate string, operand int) {

	com := Comment{
		Operate: operate,
		Operand: operand,
	}

	m := commentMessage(com, c.PublicKey)
	m.ECDSAsign(c.PrivateKey)
	byte_m := packMessage(m)
	dst_addr := server_addr
	sendMessage(byte_m, dst_addr)

	fmt.Println("客户端发送请求:", operand, operate, operate)

}
