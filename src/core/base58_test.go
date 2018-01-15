package core

import (
	"testing"
	"fmt"
)

func TestBase58Encode(t *testing.T) {
	base := []byte("1K1G5J3jotjHawA7ewGRDoEgGGiWf1cG7z")
	fmt.Println(base)
	fmt.Println(Base58Decode(base))
}
