package smol

import (
	"io/ioutil"
	"os"

	"github.com/fabulousduck/smol/bytecode"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/ir"
	"github.com/fabulousduck/smol/lexer"
)

//Smol : Defines the global attributes of the interpreter
type Smol struct {
	Tokens   []*lexer.Token
	HadError bool //TODO: use this
}

//NewSmol : Creates a new Smol instance
func NewSmol() *Smol {
	return new(Smol)
}

//RunFile : Interprets a given file
func (smol *Smol) RunFile(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	smol.Run(string(file), filename)
	if smol.HadError {
		os.Exit(65)
	}
}

//Run exectues a given script
func (smol *Smol) Run(sourceCode string, filename string) {
	l := lexer.NewLexer(filename, sourceCode)
	l.Lex()
	p := ast.NewParser(filename, l.Tokens)
	//We can ignore the second return value here as it is the amount of tokens consumed.
	//We do not need this here
	p.Ast, _ = p.Parse("")
	g := ir.NewGenerator(filename)
	g.Generate(p.Ast)
	bg := bytecode.Init(g, filename)
	bg.CreateRom()
	return
}
