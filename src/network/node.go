package network

import (
	"net"
	"log"
	"encoding/json"
)

const (
	SuccessMessage = "success"
	MaxMesaageLen  = 1024 * 1000
)

var P2PNode *Node

type Node struct {
	nodes []string // 节点
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
		if err != nil {
			log.Println(err)
		}
		go n.handle(conn)
	}
}

func (n *Node) handle(conn net.Conn) {
	remoteAddr := conn.RemoteAddr()
	log.Println(remoteAddr, "connected.")
	defer func() {
		conn.Close()
		log.Println(remoteAddr, "closed.")
	}()
	b := make([]byte, 1)
	var message []byte
	for {
		_, err := conn.Read(b)
		if err != nil {
			break
		}
		message = append(message, b...)
		if len(message) > MaxMesaageLen {
			break
		}
	}
	if m := Unmarshal(message); m != nil {
		if result := m.Resolve(); result != nil {
			conn.Write(result)
		} else {
			conn.Write([]byte(SuccessMessage))
		}
	}
}

// 添加新节点
func (n *Node) AddNode(addr string) {
	n.nodes = append(n.nodes, addr)
}

// 网络发现
func (n *Node) Discovery() {

}

// 广播消息
func (n *Node) Broadcasting(m Message) {
	for _, node := range n.nodes {
		conn, err := net.Dial("tcp4", node)
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
