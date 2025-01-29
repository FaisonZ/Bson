package main

import (
	"fmt"

	"github.com/FaisonZ/bson/internal/bit"
	"github.com/FaisonZ/bson/internal/bson"
)

func runBitBuilder() {
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
	fmt.Println(b)
}

func runEncodeJson() {
	//var jsonBlob = []byte(`[{"foo":"bar"}, "foobar", true, false, null]`)
	//var jsonBlob = []byte(`{"bing": null, "baz": true, "bar": false}`)
	var jsonBlob = []byte(`["a", null, "b", true, "ab"]`)

	bb := bit.NewBitBuilder()
	err := bson.EncodeJson(jsonBlob, bb)
	if err != nil {
		fmt.Printf("Bson Encoding error: %q\n", err)
		return
	}

	fmt.Println(bb)
}

func main() {
	runEncodeJson()
	runBitBuilder()
}
