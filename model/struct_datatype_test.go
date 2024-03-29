package model

import "testing"

func TestEmptyStructIsWellGenerated(t *testing.T) {
	expected := "{\n\n}"
	emptyStructType := NewStructDataType("emptyStructType")
	generated := emptyStructType.Generate(nil)
	if expected != generated {
		t.Errorf("Failed to generate type random example. Expected:\n%s\n, but got:\n%s", expected, generated)
	}
}
