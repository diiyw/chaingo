package network

import (
	"net"
	"log"
	"io/ioutil"
	"encoding/json"
	"time"
)

const (
	SuccessMessage = "success"
	timeOut        = 1e9 * 30
)

var P2PNode *Node

type Node struct {
	nodes []Node // 伙伴节点
	addr  *net.TCPAddr
}

func NewNode(ip string, port int) *Node {
	P2PNode = &Node{
		nodes: nil,
		addr: &net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		},
	}
	return P2PNode
}

// 监听网络服务
func (n *Node) Serving() {
	listener, err := net.ListenTCP("tcp4", n.addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		conn.SetDeadline(time.Now().Add(timeOut))
		if err != nil {
			log.Println(err)
		}
		message, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Println(err)
		}
		var m Message
		if json.Unmarshal(message, m) == nil {
			if result := m.Resolve(); result != nil {
				conn.Write(result)
			} else {
				conn.Write([]byte(SuccessMessage))
			}
		}
		conn.Close()
	}
}

// 添加新节点
func (n *Node) AddPartner(addr *net.TCPAddr) {
	n.nodes = append(n.nodes, Node{addr: addr})
}

// 网络发现
func (n *Node) Discovery() {

}

// 广播消息
func (n *Node) Broadcasting(m Message) {
	for _, node := range n.nodes {
		conn, err := net.Dial("tcp4", node.addr.String())
		if err != nil {
			log.Println(err)
			continue
		}
		b, err := json.Marshal(m)
		_, err = conn.Write(b)
		if err != nil {
			log.Println(err)
		}
		conn.Close()
	}
}

func (n *Node) String() string {
	return n.addr.String()
}
