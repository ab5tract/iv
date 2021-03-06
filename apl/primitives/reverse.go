package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⌽",
		doc:    "reverse",
		Domain: Monadic(nil),
		fn:     revLast,
		sel:    selection(revLast),
	})
	register(primitive{
		symbol: "⊖",
		doc:    "reverse first",
		Domain: Monadic(nil),
		fn:     revFirst,
		sel:    selection(revFirst),
	})
	// TODO reverse with axis

	register(primitive{
		symbol: "⌽",
		doc:    "rotate",
		Domain: Dyadic(Split(ToIndexArray(nil), nil)),
		fn:     rotLast,
		sel:    selection(rotLast),
	})
	register(primitive{
		symbol: "⊖",
		doc:    "rotate first",
		Domain: Dyadic(Split(ToIndexArray(nil), nil)),
		fn:     rotFirst,
		sel:    selection(rotFirst),
	})
	// TODO rotate with axis
}

func revLast(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return reverse(a, R, -1)
}
func revFirst(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return reverse(a, R, 0)
}
func reverse(a *apl.Apl, R apl.Value, axis int) (apl.Value, error) {
	if _, ok := R.(apl.Axis); ok {
		if r, n, err := splitAxis(a, R); err != nil {
			return nil, err
		} else {
			R = r
			if len(n) != 1 {
				return nil, fmt.Errorf("reverse with axis: axis must be a scalar or length 1")
			}
			axis = n[0]
		}
	}

	ar, ok := R.(apl.Array)

	// Scalar values are returned as scalars.
	if ok == false {
		return R, nil
	}

	shape := ar.Shape()
	if axis < 0 {
		axis = len(shape) + axis
	}
	if axis < 0 || axis >= len(shape) {
		return nil, fmt.Errorf("reverse: axis out of range: %d  (rank %d)", axis, len(shape))
	}

	res := apl.MakeArray(ar, nil)
	ic, src := apl.NewIdxConverter(shape)
	dst := make([]int, len(shape)) // dst index vector
	for i := 0; i < res.Size(); i++ {
		copy(src, dst) // sic: copy dst over src
		src[axis] = shape[axis] - src[axis] - 1
		res.Set(i, ar.At(ic.Index(src)).Copy())
		apl.IncArrayIndex(dst, shape)
	}
	return res, nil
}

func rotLast(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	return rotate(a, L, R, -1)
}
func rotFirst(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	return rotate(a, L, R, 0)
}

// rotate R around axis by L.
// L is a Vector which should be convertable to index.
// It must have the shape of R, with the axis removed.
// It it is a single element array, it is repeated to conform.
func rotate(a *apl.Apl, L, R apl.Value, axis int) (apl.Value, error) {
	if _, ok := R.(apl.Axis); ok {
		if r, n, err := splitAxis(a, R); err != nil {
			return nil, err
		} else {
			R = r
			if len(n) != 1 {
				return nil, fmt.Errorf("rotate with axis: axis must be a scalar or length 1")
			}
			axis = n[0]
		}
	}

	ar, ok := R.(apl.Array)

	// Scalar R are returned as scalars.
	if ok == false {
		return R, nil
	}

	shape := ar.Shape()

	al := L.(apl.IntArray)
	lshape := al.Shape()

	rot := func(i, n, size int) int {
		k := (i + n) % size
		if k < 0 {
			k = size + k
		}
		return k
	}

	// If R is a vector, shortcut.
	if len(shape) == 1 {
		if len(lshape) != 1 || lshape[0] != 1 {
			return nil, fmt.Errorf("rotate: wrong shape of L for vector R: %v", lshape)
		}
		n := int(al.At(0).(apl.Int))
		size := shape[0]

		res := apl.MakeArray(ar, []int{shape[0]})
		for i := 0; i < size; i++ {
			res.Set(i, ar.At(rot(i, n, size)).Copy())
		}
		return res, nil
	}

	if axis < 0 {
		axis = len(shape) + axis
	}
	if axis < 0 || axis >= len(shape) {
		return nil, fmt.Errorf("rotate: illeal axis: %d (rank: %d)", axis, len(shape))
	}

	// Extend L to conform, if it is a single element array.
	if len(lshape) == 1 && len(shape) > 2 {
		newshape := make([]int, len(shape)-1)
		for i := range newshape {
			newshape[i] = lshape[0]
		}
		if rs, ok := L.(apl.Reshaper); ok {
			newl := rs.Reshape(newshape)
			al = newl.(apl.IntArray)
			lshape = al.Shape()
		} else {
			return nil, fmt.Errorf("rotate: cannot reshape L") // this should not happen
		}
	}

	if len(lshape) != len(shape)-1 {
		return nil, fmt.Errorf("rotate L: has wrong rank: %d (R: %d)", len(lshape), len(shape))
	}
	for i := range lshape {
		k := i
		if i >= axis {
			k = i + 1
		}
		if shape[k] != lshape[i] {
			return nil, fmt.Errorf("rotate L: has wrong shape: %v (R: %v)", lshape, shape)
		}
	}

	res := apl.MakeArray(ar, nil)
	lic, idx := apl.NewIdxConverter(lshape)
	ric, src := apl.NewIdxConverter(shape)
	dst := make([]int, len(shape))
	axsize := shape[axis]
	for i := 0; i < res.Size(); i++ {
		// Calculate the rotation number n.
		// Copy dst over idx, omitting axis
		copy(idx, dst[:axis])
		copy(idx[axis:], dst[axis+1:])
		n := int(al.At(lic.Index(idx)).(apl.Int))
		copy(src, dst)                        // sic: copy dst over src
		src[axis] = rot(dst[axis], n, axsize) // replace the axis by it's rotation
		res.Set(i, ar.At(ric.Index(src)).Copy())
		apl.IncArrayIndex(dst, shape)
	}
	return res, nil
}
