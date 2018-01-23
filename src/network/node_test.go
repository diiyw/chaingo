package network

import (
	"testing"
	"net"
	"log"
	"time"
	"encoding/json"
	"fmt"
)

func TestNode_Serving(t *testing.T) {
	NewNode("0.0.0.0", 1024)
	go P2PNode.Serving()
	time.Sleep(1e9 * 3)
	conn, err := net.Dial("tcp4", "127.0.0.1:1024")
	if err != nil {
		log.Println(err)
	}
	conn.Write([]byte("closed"))
	conn.Close()
	time.Sleep(1e9 * 2)
}

func TestNode_AddNode(t *testing.T) {
	NewNode("0.0.0.0", 1024)
	go P2PNode.Serving()
	time.Sleep(1e9 * 3)
	conn, err := net.Dial("tcp4", "127.0.0.1:1024")
	if err != nil {
		log.Println(err)
	}
	b, err := json.Marshal(RelNode{"127.0.0.1", 1025})
	if err != nil {
		log.Println(err)
	}
	if b == nil {
		t.Error("message error.")
	}
	conn.Write(b)
	conn.Close()
	time.Sleep(1e9 * 2)
	for _, node := range P2PNode.nodes {
		fmt.Println(node)
	}
}
