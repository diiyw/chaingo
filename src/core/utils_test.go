package core

import (
	"testing"
	"fmt"
)

func TestReverseBytes(t *testing.T) {
	data := []byte{1, 2, 3}
	ReverseBytes(data)
	fmt.Println(data)
}
