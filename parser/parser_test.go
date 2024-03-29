package parser

import (
	"errors"
	"log"
	"testing"

	"github.com/plbalbi/json-example-generator/model"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		testDescription            string
		input                      string
		expectedError              error
		expectedStructDeclarations []string
		resultPredicate            func(*Result) bool
	}{
		{
			"single struct is parser correctly",
			"type perro struct { }",
			nil,
			[]string{"perro"},
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"two structs are parsed correctly",
			"type perro struct { } type gato struct { }",
			nil,
			[]string{"perro", "gato"},
			func(result *Result) bool {
				return model.CountStructDataTypes(result.typesRepository) == 2 &&
					result.typesRepository["perro"] != nil &&
					result.typesRepository["gato"] != nil
			},
		},
		{
			"single structs single list field is parsed correctly",
			"type perro struct { hola []int }",
			nil,
			[]string{"perro"},
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"slices are not registered in repository",
			"type gato struct { hola []string }",
			nil,
			[]string{"gato"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"missing type declarations fails",
			`type test struct {
				nombre string
				edad int
				gustosDeHelado []gusto
			}`,
			errors.New("Type 'gusto' was not declared"),
			[]string{"test"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"declaring type later it valid",
			`type test struct {
				nombre string
				edad int
				gustosDeHelado []gusto
			}
			type gusto struct {
				nombre string
				granizado bool
			}`,
			nil,
			[]string{"test", "gusto"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"slices can be nested",
			`type test struct {
				nombre string
				edad int
				gustosDeHelado [][][][][]gusto
			}
			type gusto struct {
				nombre []string
				granizado bool
			}`,
			nil,
			[]string{"test", "gusto"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"test float64",
			`
			type persona struct {
				nacionalidad pais
				ventas []float64
			  }
			  type pais struct {
				nombre string
			  }
			`,
			nil,
			[]string{"persona", "pais"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"test underscore in identifier is valid",
			`
			type persona struct {
				color_pelo string
			  }
			`,
			nil,
			[]string{"persona"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"invalid character in identifier",
			`
			type persona struct {
				color-pelo! string
			  }
			`,
			errors.New("syntax error"),
			nil,
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"inline struct declaration",
			`
			type persona struct {
				nombre string
				primo struct{
					nombre string
				}
				hermano struct{
					edad int
				}
			  }
			`,
			nil,
			[]string{"persona"},
			func(result *Result) bool {
				return result.typesRepository["[]string"] == nil
			},
		},
		{
			"multiple declarations of struct is not valid",
			`
			type persona struct {
				nombre string
			}
			type persona struct {
				edad int
			}
			`,
			errors.New("Multiple declarations of type 'persona'"),
			[]string{"persona", "persona"},
			nil,
		},
		{
			"autoreferencing struct is invalid",
			`
			type persona struct {
				hermano persona
			}
			`,
			errors.New("Circular definition of type 'persona'"),
			[]string{"persona"},
			nil,
		},
		{
			"syntax error over semantic error",
			`
			type persona struct {
				primo persona
			}
			c
			`,
			errors.New("syntax error"),
			[]string{"persona"},
			nil,
		},
		/*
			NOTE: this test is not deterministic because of the error message
			{
				"indirect struct circular definition is invalid",
				`
				type persona struct {
					auto transporte
				}
				type transporte struct {
					asiento butaca
				}
				type butaca struct {
					pasajero persona
				}
				`,
				errors.New("Circular definition of type 'persona'"),
				[]string{"persona", "transporte", "butaca"},
				nil,
			},
		*/
	}

	for _, testCase := range tests {
		t.Run(testCase.testDescription, func(t *testing.T) {
			result, err := Parse(testCase.input)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedStructDeclarations, result.declaredStructs)
			if testCase.resultPredicate != nil && !testCase.resultPredicate(&result) {
				t.Errorf("Failed to evaluate test predicate")
			}
		})
	}
}

func TestMapStuff(t *testing.T) {
	repository := model.GetDefaultDataTypeRepository()
	nonExistingDataType := repository["caca"]
	if nonExistingDataType != nil {
		t.Errorf("Expected repository to return nil on non existing data type")
		log.Print("hola")
	}
}
