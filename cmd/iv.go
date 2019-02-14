package cmd

import (
	"io"

	"github.com/ktye/iv/apl"
)

func Iv(a *apl.Apl, p string, w io.Writer) error {
	a.SetOutput(w)
	if err := a.ParseAndEval(`r←{<⍤⍵<0}⋄s←{⍵⍴<⍤0<0}`); err != nil {
		return err
	}
	return a.ParseAndEval(p)
}
