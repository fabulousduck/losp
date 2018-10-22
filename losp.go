package losp

import (
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
)

//Losp : Defines the global attributes of the interpreter
type Losp struct {
	Tokens   []*token
	HadError bool
}

//NewLosp : Creates a new Losp instance
func NewLosp() *Losp {
	return new(Losp)
}

//RunFile : Interprets a given file
func (losp *Losp) RunFile(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	losp.run(string(file))
	if losp.HadError {
		os.Exit(65)
	}
}

func (losp *Losp) run(sourceCode string) {
	l := new(lexer)
	l.lex(sourceCode)
	spew.Dump(l.tokens)
}