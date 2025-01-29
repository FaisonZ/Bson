package main

import (
	"fmt"
	"io"
	"os"

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
	var jsonBlob = []byte(`["a", null, "b", true, "ab"]`)

	bb := bit.NewBitBuilder()
	err := bson.EncodeJson(jsonBlob, bb)
	if err != nil {
		fmt.Printf("Bson Encoding error: %q\n", err)
		return
	}

	fmt.Println(bb)
}

func runCheckCmd(iao InsAndOuts) {
	jsonBytes, err := io.ReadAll(iao.in)
	if err != nil {
		fmt.Fprintf(iao.err, "Error reading JSON: %q\n", err)
		return
	}

	bb := bit.NewBitBuilder()
	err = bson.EncodeJson(jsonBytes, bb)
	if err != nil {
		fmt.Fprintf(iao.err, "Error encoding: %q\n", err)
		return
	}

	fmt.Fprintf(
		iao.out,
		"Json size: %d\nBson size: %d\n",
		len(jsonBytes),
		len(bb.Bytes),
	)
	fmt.Fprintf(iao.out, "diff: %d\n", len(jsonBytes)-len(bb.Bytes))
}

func runEncodeCmd(iao InsAndOuts) {
	jsonBytes, err := io.ReadAll(iao.in)
	if err != nil {
		fmt.Fprintf(iao.err, "Error reading JSON: %q\n", err)
		return
	}

	bb := bit.NewBitBuilder()
	err = bson.EncodeJson(jsonBytes, bb)
	if err != nil {
		fmt.Fprintf(iao.err, "Error encoding: %q\n", err)
		return
	}

	_, err = bb.WriteTo(iao.out)
	if err != nil {
		fmt.Fprintf(iao.err, "Error writing bson to target: %q\n", err)
		return
	}
}

func getCommand(args []string) (cmd string, err error) {
	if len(args) == 2 {
		switch args[1] {
		case "encode":
			return "encode", nil
		case "decode":
			return "decode", nil
		case "check":
			return "check", nil
		}
	}

	return "", fmt.Errorf(
		"Invalid command. Valid commands: bson encode, bson decode, bson check",
	)
}

type InsAndOuts struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func main() {
	cmd, err := getCommand(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
		return
	}

	insOuts := InsAndOuts{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}

	switch cmd {
	case "encode":
		runEncodeCmd(insOuts)
	case "decode":
	case "check":
		runCheckCmd(insOuts)
	}
}
