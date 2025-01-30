package bson

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/FaisonZ/bson/internal/bit"
	"github.com/FaisonZ/bson/internal/nums"
)

const (
	BSON_VERSION  = 0b0001
	OBJECT_TOKEN  = 0b001
	ARRAY_TOKEN   = 0b010
	STRING_TOKEN  = 0b011
	INTEGER_TOKEN = 0b100
	BOOLEAN_TOKEN = 0b101
	NULL_TOKEN    = 0b110
	FLOAT_TOKEN   = 0b111

	FALSE = 0b0
	TRUE  = 0b1

	MAX_CHUNK_LENGTH = 0b1_1111

	INT8_TOKEN  = 0b00
	INT16_TOKEN = 0b01
	INT32_TOKEN = 0b10
	INT64_TOKEN = 0b11
)

func EncodeJson(j []byte, bb *bit.BitBuilder) error {
	var a any
	err := json.Unmarshal(j, &a)

	if err != nil {
		return fmt.Errorf("JSON error: %v", err)
	}

	writeVersion(bb)
	//fmt.Printf("v: %08b\n", bb.Bytes)
	encodeValue(a, bb)
	return nil
}

func writeVersion(bb *bit.BitBuilder) {
	bb.AddBits(BSON_VERSION, 4)
}

func writeToken(tokenByte byte, bb *bit.BitBuilder) error {
	bb.AddBits(tokenByte, 3)
	return nil
}

func writeLength(lengthByte byte, bb *bit.BitBuilder) error {
	bb.AddBits(lengthByte, 5)
	return nil
}

func writeTokenWithLength(
	tokenByte byte,
	lengthByte byte,
	bb *bit.BitBuilder,
) error {
	writeToken(tokenByte, bb)
	writeLength(lengthByte, bb)
	return nil
}

func writeString(s string, bb *bit.BitBuilder) error {
	bb.AddBytes([]byte(s))

	return nil
}

func encodeValue(a any, bb *bit.BitBuilder) error {
	if o, ok := a.(map[string]any); ok {
		encodeObject(o, bb)
	} else if s, ok := a.(string); ok {
		encodeString(s, bb)
	} else if ar, ok := a.([]any); ok {
		encodeArray(ar, bb)
	} else if b, ok := a.(bool); ok {
		encodeBoolean(b, bb)
	} else if n, ok := a.(float64); ok {
		encodeNumber(n, bb)
	} else if a == nil {
		encodeNull(bb)
	}

	return nil
}

func encodeString(s string, bb *bit.BitBuilder) error {
	writeToken(STRING_TOKEN, bb)
	encodeStringChunks(s, bb)

	return nil
}

func encodeStringChunks(s string, bb *bit.BitBuilder) error {
	strToWrite := s
	strRemaining := ""

	if len(strToWrite) >= MAX_CHUNK_LENGTH {
		strToWrite = strToWrite[:MAX_CHUNK_LENGTH]
		strRemaining = s[MAX_CHUNK_LENGTH:]
	}

	//fmt.Printf("Len: %d\n", len(strToWrite))
	writeLength(byte(len(strToWrite)), bb)
	writeString(strToWrite, bb)

	if len(strToWrite) >= MAX_CHUNK_LENGTH {
		encodeStringChunks(strRemaining, bb)
	}

	return nil
}

func encodeNull(bb *bit.BitBuilder) error {
	writeToken(NULL_TOKEN, bb)
	return nil
}

func encodeBoolean(b bool, bb *bit.BitBuilder) error {
	writeToken(BOOLEAN_TOKEN, bb)

	if b {
		bb.AddBits(TRUE, 1)
	} else {
		bb.AddBits(FALSE, 1)
	}
	return nil
}

func encodeArray(a []any, bb *bit.BitBuilder) error {
	writeToken(ARRAY_TOKEN, bb)
	return encodeArrayChunk(a, bb)
}

func encodeArrayChunk(a []any, bb *bit.BitBuilder) error {
	arrToWrite := a
	arrRemaining := []any{}

	if len(arrToWrite) >= MAX_CHUNK_LENGTH {
		arrToWrite = arrToWrite[:MAX_CHUNK_LENGTH]
		arrRemaining = a[MAX_CHUNK_LENGTH:]
	}

	//fmt.Printf("Len: %d\n", len(strToWrite))
	writeLength(byte(len(arrToWrite)), bb)
	for _, val := range arrToWrite {
		encodeValue(val, bb)
	}

	if len(arrToWrite) >= MAX_CHUNK_LENGTH {
		encodeArrayChunk(arrRemaining, bb)
	}

	return nil
}

func encodeObject(o map[string]any, bb *bit.BitBuilder) error {
	writeToken(OBJECT_TOKEN, bb)

	keys := make([]string, 0, len(o))
	for key := range o {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	return encodeObjectChunk(keys, o, bb)
}

func encodeObjectChunk(
	keys []string,
	o map[string]any,
	bb *bit.BitBuilder,
) error {
	keysToWrite := keys
	keysRemaining := []string{}

	if len(keysToWrite) >= MAX_CHUNK_LENGTH {
		keysToWrite = keysToWrite[:MAX_CHUNK_LENGTH]
		keysRemaining = keys[MAX_CHUNK_LENGTH:]
	}

	//fmt.Printf("Len: %d\n", len(strToWrite))
	writeLength(byte(len(keysToWrite)), bb)
	for _, key := range keysToWrite {
		encodeString(key, bb)
		encodeValue(o[key], bb)
	}

	if len(keysToWrite) >= MAX_CHUNK_LENGTH {
		encodeObjectChunk(keysRemaining, o, bb)
	}

	return nil
}

func encodeNumber(n float64, bb *bit.BitBuilder) error {
	if nums.IsInt(n) {
		return encodeInt(n, bb)
	} else {
		fmt.Println("I haven't implemented floats yet")
	}

	return nil
}

func encodeInt(n float64, bb *bit.BitBuilder) error {
	writeToken(INTEGER_TOKEN, bb)
	intBytes := make([]byte, 0, 4)

	intSize, err := nums.MinIntSize(n)
	if err != nil {
		return nil
	}

	var intSizeToken byte

	switch intSize {
	case 8:
		intSizeToken = INT8_TOKEN
		intBytes = append(intBytes, byte(int(n)))
	case 16:
		intSizeToken = INT16_TOKEN
		// Note, AppendUint16 does not convert an int to an unsigned int, it
		// just takes the bytes and appends them, only having regard to the
		// number of bytes
		intBytes = binary.BigEndian.AppendUint16(intBytes, uint16(n))
	case 32:
		intSizeToken = INT32_TOKEN
		intBytes = binary.BigEndian.AppendUint32(intBytes, uint32(n))
	case 64:
		intSizeToken = INT64_TOKEN
		intBytes = binary.BigEndian.AppendUint64(intBytes, uint64(n))
	default:
		return fmt.Errorf("Invalid Int size found: %d", intSize)
	}

	bb.AddBits(intSizeToken, 2)
	bb.AddBytes(intBytes)

	return nil
}
