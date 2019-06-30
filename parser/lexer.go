package parser

import (
	"log"
	"strings"
	"text/scanner"

	"github.com/golang-collections/collections/stack"
)

const TYPE_TOKEN_STRING = "type"
const STRUCT_TOKEN_STRING = "struct"
const OPEN_CBRASE_TOKEN_STRING = "{"
const CLOSE_CBRASE_TOKEN_STRING = "}"

type lexer struct {
	result         Result
	scan           scanner.Scanner
	states         stack.Stack
	currentSymType *yySymType
}

//Result is the object in which the parser transmits the parsed text.
type Result struct {
	structsCount int
}

//stateFunction is the base type for lexer states. Contains all the logic for lexing the text,
//communicating found values with the parser, and choosing the following states.
type stateFunction func(*lexer) tokenType

type tokenType int

//Lex is somehow like the tokenStream.next() called it time it needs by the parser
func (lex *lexer) Lex(currentSymType *yySymType) int {
	lex.currentSymType = currentSymType
	return int(lex.scanUntilTokenFound())
}

// TODO: Implement lexer/parser error handling.
func (lex *lexer) Error(message string) {
}

func (lex *lexer) scanAndLog() string {
	lex.scan.Scan()
	log.Printf("%s: %s", lex.scan.Position, lex.scan.TokenText())
	return lex.scan.TokenText()
}

func (lex *lexer) scanUntilTokenFound() tokenType {
	stateToRun := lex.states.Pop()
	// Had to hardcode the stateFunction type in here, no way out
	if currentStateFunction, ok := stateToRun.(func(*lexer) tokenType); ok {
		// casting and calling current state
		return stateFunction(currentStateFunction)(lex)
	}
	log.Fatal("Failed to cast state to stateFunction type")
	return 0 // Horrible way, should return some other error instead of EOF
}

//Parse does the actual parsing
func Parse(inputStream string) (Result, error) {
	lex := newLexer(inputStream)
	yyParse(lex)
	return lex.result, nil
}

func newLexer(inputStream string) *lexer {
	brandNewLexer := &lexer{}
	brandNewLexer.scan.Init(strings.NewReader(inputStream))
	brandNewLexer.states.Push(typeDefinitionOpeningState)
	return brandNewLexer
}

//typeDefinitionOpeningState lexes the whole struct definition string 'type ID struct {'
func typeDefinitionOpeningState(lex *lexer) tokenType {
	// Lex 'type' type identifier
	if lex.scanAndLog() != TYPE_TOKEN_STRING {
		return 0
	}
	// Lex type identifier
	// TODO: Add keywords restrictions
	typeIdentifier := lex.scanAndLog()
	lex.currentSymType.value = typeIdentifier
	if lex.scanAndLog() != STRUCT_TOKEN_STRING {
		return 0
	}
	if lex.scanAndLog() != OPEN_CBRASE_TOKEN_STRING {
		return 0
	}
	// //TODO: Push real state. This is for testing
	lex.states.Push(typeDefinitionOpeningState)
	return TypeOpening
}
