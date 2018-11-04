// Package apl implements an APL interpreter.
package apl

import (
	"fmt"
	"io"

	"github.com/ktye/iv/apl/scan"
	// _ "github.com/ktye/iv/apl/funcs" // Register default funcs
)

// New starts a new interpreter.
func New(w io.Writer) *Apl {
	a := Apl{
		stdout:     w,
		vars:       make(map[string]Value),
		Origin:     1,
		primitives: make(map[Primitive][]FunctionHandle),
		operators:  make(map[string]Operator),
		symbols:    make(map[rune]string),
		doc:        make(map[string]string),
	}
	a.parser.a = &a
	return &a
}

// Apl stores the interpreter state.
type Apl struct {
	scan.Scanner
	parser
	stdout     io.Writer
	vars       map[string]Value
	primitives map[Primitive][]FunctionHandle
	operators  map[string]Operator
	symbols    map[rune]string
	doc        map[string]string
	initscan   bool
	format     Format
	Origin     int
	debug      bool
}

// Parse parses a line of input into the current context.
// It returns a Program which can be evaluated.
func (a *Apl) Parse(line string) (Program, error) {

	// Before the scanner is used for the first time,
	// tell it about all registered symbols.
	if a.initscan == false {
		m := make(map[rune]string)
		for r, s := range a.symbols {
			m[r] = s
		}
		a.SetSymbols(m)
		a.initscan = true
	}

	tokens, err := a.Scan(line)
	if a.debug {
		fmt.Fprintf(a.stdout, "%s\n", scan.PrintTokens(tokens))
	}

	if err != nil {
		return nil, err
	}

	p, err := a.parse(tokens)
	if a.debug {
		fmt.Fprintf(a.stdout, "%s\n", p.String(a))
	}

	if err != nil {
		return nil, err
	} else {
		return p, nil
	}
}

func (a *Apl) ParseAndEval(line string) error {
	if p, err := a.Parse(line); err != nil {
		return err
	} else {
		return a.Eval(p)
	}
}

func (a *Apl) SetDebug(d bool) {
	a.debug = d
}

func (a *Apl) SetOutput(w io.Writer) {
	a.stdout = w
}