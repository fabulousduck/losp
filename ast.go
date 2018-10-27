package losp

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type function struct {
	name   string
	params []string
	body   []node
}

func (f function) getNodeName() string {
	return "function"
}

type variable struct {
	name  string
	value string
}

type anb struct {
	lhs  string
	rhs  string
	body []node
}

func (anb anb) getNodeName() string {
	return "anb"
}

type statement struct {
	lhs string
	rhs string
}

func (s statement) getNodeName() string {
	return "statement"
}

func (v variable) getNodeName() string {
	return "variable"
}

type functionCall struct {
	name   string
	params []string
	body   []node
}

func (fc functionCall) getNodeName() string {
	return "functionCall"
}

type node interface {
	getNodeName() string
}

type parser struct {
	ast      []node
	filename string
}

func NewParser(filename string) *parser {
	p := new(parser)
	p.filename = filename
	return p
}

func (p *parser) parse(tokens []token) ([]node, int) {
	nodes := []node{}
	// fmt.Printf("FUUUUUUUUU")
	// spew.Dump(tokens)
	for i := 0; i < len(tokens); {
		switch tokens[i].Type {
		case "variable_assignment":
			node, tokensConsumed := p.createVariable(tokens, i)
			i += tokensConsumed
			nodes = append(nodes, node)
		case "function_definition":
			node, tokensConsumed := p.createFunctionHeader(tokens, i+1)
			i += tokensConsumed + 1
			body, consumed := p.parse(tokens[i:])
			node.body = body
			i += consumed
			nodes = append(nodes, node)
		case "left_not_right":
			node, tokensConsumed := p.createLNR(tokens, i+1)
			i += tokensConsumed + 1
			body, consumed := p.parse(tokens[i:])
			node.body = body
			i += consumed
			nodes = append(nodes, node)
		case "print_integer":
			node, tokensConsumed := p.createStatement(tokens, i+1, "PRI")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "print_ascii":
			node, tokensConsumed := p.createStatement(tokens, i+1, "PRA")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "increment_value":
			node, tokensConsumed := p.createStatement(tokens, i+1, "INC")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "close_block":
			spew.Dump("END BLOCK")
			return nodes, i
		case "string":
		case "CHAR":
		case "NUMB":
		case "LEFT_BRACKET":
		case "RIGHT_BRACKET":
		case "LEFT_ARROW":
		case "RIGHT_ARROW":
		case "DOUBLE_DOT":
		case "COMMA":
		case "SEMICOLON":
		default:
			spew.Dump("fuck")
		}
	}

	return nodes, len(tokens)
}

func (p *parser) createStatement(tokens []token, index int, t string) (*statement, int) {
	fmt.Printf("creating statment of type :")
	spew.Dump(t)
	s := new(statement)
	tokensConsumed := 0

	s.lhs = t

	p.expect([]string{"string", "CHAR"}, tokens[index+tokensConsumed])
	s.rhs = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"SEMICOLON"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return s, tokensConsumed
}

func (p *parser) createLNR(tokens []token, index int) (*anb, int) {
	anb := new(anb)
	tokensConsumed := 0

	p.expect([]string{"LEFT_BRACKET"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"CHAR", "NUMB", "string"}, tokens[index+tokensConsumed])
	anb.lhs = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"COMMA"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"CHAR", "NUMB", "string"}, tokens[index+tokensConsumed])
	anb.rhs = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"RIGHT_BRACKET"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"DOUBLE_DOT"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return anb, tokensConsumed
}

func (p *parser) createFunctionHeader(tokens []token, index int) (*function, int) {
	f := new(function)
	tokensConsumed := 0
	p.expect([]string{"string", "CHAR"}, tokens[index+tokensConsumed])
	f.name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"LEFT_ARROW", "DOUBLE_DOT"}, tokens[index+tokensConsumed])
	if tokens[index+tokensConsumed].Type == "DOUBLE_DOT" {
		tokensConsumed++
		return f, tokensConsumed
	}
	tokensConsumed++
	for currentToken := tokens[index+tokensConsumed]; currentToken.Type != "RIGHT_ARROW"; currentToken = tokens[index+tokensConsumed] {
		p.expect([]string{"string", "CHAR", "COMMA"}, currentToken)
		if currentToken.Type == "COMMA" {
			p.expect([]string{"CHAR", "string"}, tokens[index+tokensConsumed+1])
			tokensConsumed++
			continue
		}
		f.params = append(f.params, currentToken.Value)
		tokensConsumed++
	}
	tokensConsumed++
	p.expect([]string{"DOUBLE_DOT"}, tokens[index+tokensConsumed])
	tokensConsumed++
	return f, tokensConsumed
}

func (p *parser) createVariable(tokens []token, index int) (*variable, int) {
	variable := new(variable)
	tokensConsumed := 0
	expectedNameTypes := []string{
		"CHAR",
		"STRING",
	}
	p.expect(expectedNameTypes, tokens[index+1])
	variable.name = tokens[index+1].Value
	tokensConsumed++

	expectedValueTypes := []string{"NUMB"}
	p.expect(expectedValueTypes, tokens[index+2])
	variable.value = tokens[index+2].Value
	tokensConsumed++

	return variable, tokensConsumed
}

func (p *parser) expect(expectedValues []string, token token) {
	if !contains(token.Type, expectedValues) {
		throwSemanticError(&token, expectedValues, p.filename)
		os.Exit(65)
	}
}
