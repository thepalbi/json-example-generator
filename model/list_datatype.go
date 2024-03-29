package model

import (
	"bytes"
	"math/rand"
	"time"
)

var MaxRandomListSize = 10

//ListDataType represents a list type, containing multiple elements of a certain DataType.
type ListDataType struct {
	name      string
	innerType DataType
}

//NewListDataType returns a brand new ListDataType, with a given name and inner type.
func NewListDataType(aName string, aType DataType) *ListDataType {
	return &ListDataType{
		name:      aName,
		innerType: aType,
	}
}

//GetName shows the datatype's name
func (data *ListDataType) GetName() string {
	return data.name
}

//IsSimple return true if the datatype is a ListDataType.
func (data *ListDataType) IsSimple() bool {
	return false
}

//IsList return true if the datatype is a ListDataType.
func (data *ListDataType) IsList() bool {
	return true
}

//IsStruct return true if the datatype is a ListDataType.
func (data *ListDataType) IsStruct() bool {
	return false
}

//Generate generates a random example of this datatype.
func (data *ListDataType) Generate(repository DataTypeRepository) string {
	if MaxRandomListSize < 1 {
		return "[]"
	}
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	//Use as list size a random number, mod 10, + 1 to assure that is >= 1
	listSize := randomGenerator.Int()%MaxRandomListSize + 1
	MaxRandomListSize--
	var randomListBuffer bytes.Buffer
	randomListBuffer.WriteString("[")
	for i := 1; i <= listSize; i++ {
		randomListBuffer.WriteString(data.innerType.Generate(repository))
		if i < listSize {
			randomListBuffer.WriteString(", ")
		}
	}
	randomListBuffer.WriteString("]")
	return randomListBuffer.String()
}
