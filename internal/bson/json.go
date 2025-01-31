package bson

import (
	"encoding/json"
	"fmt"
	"strings"
)

func jsonUnmarshal(b []byte) (any, error) {
	var j json.RawMessage
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return unmarshalRaw(j)
}

func unmarshalRaw(m json.RawMessage) (any, error) {
	t, err := getTypeFromRaw(m)
	if err != nil {
		return nil, err
	}

	switch t {
	default:
		return nil, fmt.Errorf("Unexpected Json type: %s", t)
	case "object":
		return unmarshalObject(m)
	case "array":
		return unmarshalArray(m)
	case "string":
		return unmarshalString(m)
	case "boolean":
		return unmarshalBoolean(m)
	case "integer":
		return unmarshalInteger(m)
	case "float":
		return unmarshalFloat(m)
	case "null":
		return nil, nil
	}
}

func unmarshalString(m json.RawMessage) (string, error) {
	var s string
	err := json.Unmarshal(m, &s)
	if err != nil {
		return "", err
	}

	return s, nil
}

func unmarshalBoolean(m json.RawMessage) (bool, error) {
	var b bool
	err := json.Unmarshal(m, &b)
	if err != nil {
		return false, err
	}

	return b, nil
}

func unmarshalArray(m json.RawMessage) ([]any, error) {
	var rawArr []json.RawMessage
	err := json.Unmarshal(m, &rawArr)
	if err != nil {
		return []any{}, nil
	}

	arr := make([]any, 0, len(rawArr))
	for _, rm := range rawArr {
		val, err := unmarshalRaw(rm)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	return arr, nil
}

func unmarshalObject(m json.RawMessage) (map[string]any, error) {
	var rawObj map[string]json.RawMessage
	err := json.Unmarshal(m, &rawObj)
	if err != nil {
		return map[string]any{}, nil
	}

	obj := make(map[string]any, len(rawObj))
	for key, rm := range rawObj {
		val, err := unmarshalRaw(rm)
		if err != nil {
			return obj, err
		}
		obj[key] = val
	}

	return obj, nil
}

func unmarshalInteger(m json.RawMessage) (int64, error) {
	var rawInt int64
	err := json.Unmarshal(m, &rawInt)
	if err != nil {
		return 0, err
	}

	return rawInt, nil
}

func unmarshalFloat(m json.RawMessage) (float64, error) {
	var f float64
	err := json.Unmarshal(m, &f)
	if err != nil {
		return 0, err
	}

	return f, nil
}

func getTypeFromRaw(m json.RawMessage) (string, error) {
	if len(m) == 0 {
		return "", fmt.Errorf("Invalid JSON value")
	}

	fChar := m[0]

	switch fChar {
	case []byte("{")[0]:
		return "object", nil
	case []byte("[")[0]:
		return "array", nil
	case []byte("\"")[0]:
		return "string", nil
	case []byte("t")[0]:
		return "boolean", nil
	case []byte("f")[0]:
		return "boolean", nil
	case []byte("n")[0]:
		return "null", nil
	default:
		return getNumberTypeFromRaw(m)
	}
}

func getNumberTypeFromRaw(m json.RawMessage) (string, error) {
	if len(m) == 0 {
		return "", fmt.Errorf("Invalid JSON value")
	}

	if strings.ContainsAny(string(m), ".eE") {
		return "float", nil
	}

	return "integer", nil
}
