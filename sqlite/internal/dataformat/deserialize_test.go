package dataformat

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestDeserializeVarint(t *testing.T) {
	var testVal uint64 = 0b1111101000
	buf := []byte{0b10000111, 0b01101000}
	decodedVal, bytesRead := DeserializeVarint(buf)

	if testVal != decodedVal {
		t.Errorf(`Bytes Read = %d
		Decoded Value = %b
		Actual Value = %b`,
			bytesRead, decodedVal, testVal)
	}
}

func TestDeserializeInteger(t *testing.T) {
	var testVal int64 = 0b1111101000
	//
	buf := []byte{0b11, 0b11101000}
	decodedVal := DeserializeInteger(buf)

	if testVal != decodedVal {
		t.Errorf(`
		Decoded Value = %b 
		Actual Value = %b`,
			decodedVal, testVal)
	}
}

func TestDeserializeFloat(t *testing.T) {
	expected := 3.141592653589793 // Example float64 value
	bytes := make([]byte, 8)      // Create a byte slice to hold the binary representation

	// Serialize the float64 value into bytes
	binary.BigEndian.PutUint64(bytes, math.Float64bits(expected))

	result := DeserializeFloat(bytes)

	if result != expected {
		t.Errorf("DeserializeFloat failed: got %f, expected %f", result, expected)
	}
}
