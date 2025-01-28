package main

import (
	"fmt"

	"github.com/FaisonZ/bson/internal/bit"
)

func main() {
	fmt.Println("Hello world")
	b := bit.NewBitBuilder()
	b.AddBits(0b0001, 4)
	b.AddBits(0b001, 3)
	b.AddBits(0b00001, 5)
	b.AddBits(0b011, 3)
	b.AddBits(0b00011, 5)
	b.AddBytes([]byte("foo"))
	b.AddBits(0b011, 3)
	b.AddBits(0b00011, 5)
	b.AddBytes([]byte("bar"))
}
