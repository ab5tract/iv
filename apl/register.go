package apl

import "unicode/utf8"

// RegisterPrimitive registeres the function handle for the symbol s
// as a primitive function.
// Multiple versions may be registered, which can handle different
// argument type combinations.
// When the function is applied, the last registered handle is tested first.
func (a *Apl) RegisterPrimitive(p Primitive, h FunctionHandle) {
	a.primitives[p] = append(a.primitives[p], h)
	a.registerSymbol(string(p))
}

// RegisterOperator registers s as the symbol for the operator.
func (a *Apl) RegisterOperator(s string, op Operator) {
	a.operators[s] = op
	a.registerSymbol(s)
}

// RegisterDoc adds help text to a key in the documentation.
// This should be called by packages, that add primitives or operators.
func (a *Apl) RegisterDoc(key, help string) {
	a.doc[key] += help
}

// registerSymbol adds single rune symbols for the parser.
func (a *Apl) registerSymbol(s string) {
	if r, w := utf8.DecodeRuneInString(s); w == len(s) {
		a.symbols[r] = s
	}
}