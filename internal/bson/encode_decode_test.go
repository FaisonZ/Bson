package bson

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/FaisonZ/bson/internal/bit"
)

func TestEncodeDecodeObject(t *testing.T) {
	inJson := []byte(
		`{"foo":"bar","true":true,"false":false,"null":null,"ar":["a", "b"],"ob":{"f":"b"}}`,
	)
	var inJsonObj any
	err := json.Unmarshal(inJson, &inJsonObj)
	if err != nil {
		t.Errorf("invalid test json string: %q", err)
	}

	bb := bit.NewBitBuilder()
	err = EncodeJson(inJson, bb)
	if err != nil {
		t.Errorf("Failed to encode JSON: %q", err)
	}

	fmt.Printf("Json size: %d\nBson size: %d\n", len(inJson), len(bb.Bytes))

	decoded, err := Decode(bb.Bytes)
	if err != nil {
		t.Errorf("Failed to decode JSON: %q", err)
	}

	outJson, err := json.Marshal(decoded)
	if err != nil {
		t.Errorf("Failed to Marshal JSON: %q", err)
	}

	inMap := inJsonObj.(map[string]any)
	outMap := decoded.(map[string]any)

	inKeys := make([]string, 0, len(inMap))
	for k := range inMap {
		inKeys = append(inKeys, k)
	}
	outKeys := make([]string, 0, len(outMap))
	for k := range outMap {
		outKeys = append(outKeys, k)
	}

	slices.Sort(inKeys)
	slices.Sort(outKeys)

	if !slices.Equal(inKeys, outKeys) {
		t.Errorf("Expected object keys:\n%v\nGot:\n%v", outKeys, inKeys)
	}

	// Need more work to fully test in and out maps

	fmt.Printf("%s\n", inJson)
	fmt.Printf("%s\n", outJson)
}

func TestEncodeDecodeArray(t *testing.T) {
	inJson := []byte(
		`["a", "b", null, true, false]`,
	)
	var inJsonArr any
	err := json.Unmarshal(inJson, &inJsonArr)
	if err != nil {
		t.Errorf("invalid test json string: %q", err)
	}

	bb := bit.NewBitBuilder()
	err = EncodeJson(inJson, bb)
	if err != nil {
		t.Errorf("Failed to encode JSON: %q", err)
	}

	fmt.Printf("Json size: %d\nBson size: %d\n", len(inJson), len(bb.Bytes))

	decoded, err := Decode(bb.Bytes)
	if err != nil {
		t.Errorf("Failed to decode JSON: %q", err)
	}

	outJson, err := json.Marshal(decoded)
	if err != nil {
		t.Errorf("Failed to Marshal JSON: %q", err)
	}

	inArr := inJsonArr.([]any)
	outArr := decoded.([]any)

	if !slices.Equal(inArr, outArr) {
		t.Errorf("In Array doesn't match out array")
	}

	fmt.Printf("%s\n", inJson)
	fmt.Printf("%s\n", outJson)
}
