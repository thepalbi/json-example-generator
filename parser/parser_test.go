package parser

import (
	"testing"
)

// TODO: Make a go-like test infrastructure, in which one is able to just define
// the string to be parsed, and the expected error (be it nil or the error itself)

func TestParserCountsASingleStruct(t *testing.T) {
	result, err := Parse("type perro struct { }")
	if err != nil {
		t.Errorf("Test failed with error: %s", err.Error())
	}
	if result.structsCount != 1 {
		t.Errorf("Expected %d structs to be parsed, but found %d", 1, result.structsCount)
	}
}

func TestParserCountsTwoStructs(t *testing.T) {
	result, err := Parse("type perro struct { } type gato struct { }")
	if err != nil {
		t.Errorf("Test failed with error: %s", err.Error())
	}
	if result.structsCount != 2 {
		t.Errorf("Expected %d structs to be parsed, but found %d", 2, result.structsCount)
	}
}

func TestStructWithOneFieldIsCorrectlyIdentified(t *testing.T) {
	result, err := Parse("type perro struct { perro gato }")
	if err != nil {
		t.Errorf("Test failed with error: %s", err.Error())
	}
	if result.structsCount != 1 {
		t.Errorf("Expected %d structs to be parsed, but found %d", 1, result.structsCount)
	}
}

func TestStructWithTwoFieldsIsCorrectlyIdentified(t *testing.T) {
	result, err := Parse("type perro struct { perro gato atr perro }")
	if err != nil {
		t.Errorf("Test failed with error: %s", err.Error())
	}
	if result.structsCount != 1 {
		t.Errorf("Expected %d structs to be parsed, but found %d", 1, result.structsCount)
	}
}
