package raft

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"fmt"
)

// 生成公私钥,直接给node加上
func (n *Node) ECDSAgenerate() {
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

	n.PublicKey = publicBytes
	n.PrivateKey = privateBytes

}

// 咱就是说传进去一个message struct，出来里面已经加上签名了
func (m *Message) ECDSAsign(privatekey []byte) {
	// 先解DER编码为ecdsa privatekey形式
	priKey, err := x509.ParseECPrivateKey(privatekey)
	if err != nil {
		panic(err)
	}

	// 然后先把mes除了sig之外的部分打包成一个byte数组然后hash
	messbytes, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	hash := sha1.Sum(messbytes)

	// 然后用私钥对hash签名
	sign, err := ecdsa.SignASN1(rand.Reader, priKey, hash[:])
	if err != nil {
		fmt.Println(err)
	}
	m.Sign = sign
}

//  传一个message进去，返回bool说明message内容有没有变
func (m *Message) ECDSAverify() bool {
	//x509解码
	publicStream, err := x509.ParsePKIXPublicKey(m.Key)
	if err != nil {
		fmt.Println(err)
	}
	// 接口转换成公钥
	pubKey := publicStream.(*ecdsa.PublicKey)

	// 然后先把mes除了sg之外的部分打包成一个byte数组然后hash
	// woc一开始传地址，所以下面这个赋值nil直接把m的值也给改了。。
	m_nosign := *m
	m_nosign.Sign = nil //就是说把sign那部分弄成初始
	messbytes, err := json.Marshal(m_nosign)
	if err != nil {
		fmt.Println(err)
	}
	hash := sha1.Sum(messbytes)

	// 然后用公钥验证消息是否改变
	result := ecdsa.VerifyASN1(pubKey, hash[:], m.Sign)
	if err != nil {
		fmt.Println(err)
	}
	return result
}
