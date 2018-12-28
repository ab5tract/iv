package numbers

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// ComplexArray is a uniform array of complex128
type ComplexArray struct {
	Dims  []int
	Cmplx []complex128
}

func (f ComplexArray) String(a *apl.Apl) string {
	return apl.ArrayString(a, f)
}

func (f ComplexArray) At(i int) (apl.Value, error) {
	if i < 0 || i >= len(f.Cmplx) {
		return nil, fmt.Errorf("index out of range")
	}
	return Complex(f.Cmplx[i]), nil
}

func (f ComplexArray) Shape() []int {
	return f.Dims
}

func (f ComplexArray) Size() int {
	return len(f.Cmplx)
}

func (f ComplexArray) Zero() apl.Value {
	return Complex(0.0)
}

func (f ComplexArray) Set(i int, v apl.Value) error {
	if i < 0 || i > len(f.Cmplx) {
		return fmt.Errorf("index out of range")
	}
	if c, ok := v.(Complex); ok {
		f.Cmplx[i] = complex128(c)
		return nil
	}
	return fmt.Errorf("cannot assign %T to ComplexArray", v)
}

func makeComplexArray(v []apl.Value) ComplexArray {
	f := make([]complex128, len(v))
	for i, e := range v {
		f[i] = complex128(e.(Complex))
	}
	return ComplexArray{
		Dims:  []int{len(v)},
		Cmplx: f,
	}
}

func (f ComplexArray) Reshape(shape []int) apl.Value {
	res := ComplexArray{
		Dims:  shape,
		Cmplx: make([]complex128, prod(shape)),
	}
	k := 0
	for i := range res.Cmplx {
		res.Cmplx[i] = f.Cmplx[k]
		k++
		if k == len(f.Cmplx) {
			k = 0
		}
	}
	return res
}