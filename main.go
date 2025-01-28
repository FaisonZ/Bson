package main

import (
	"encoding/json"
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

func main() {
	var jsonBlob = []byte(`[{
    "foo": "bar"
},
"foo",
true,
false
]`)

	var data any
	bb := bit.NewBitBuilder()

	err := json.Unmarshal(jsonBlob, &data)

	if err != nil {
		fmt.Printf("Unmarshal error: %q", err)
		return
	}

	bson.EncodeJson(data, bb)
	fmt.Println(bb)

	runBitBuilder()
}
