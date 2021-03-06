package diff

import "github.com/zclconf/go-cty/cty"

// longestCommonSubsequence finds a sequence of values that are common to both
// x and y, with the same relative ordering as in both collections.
//
// The approached used here is a "naive" one, assuming that both x and y will
// generally be small in the use-cases cty is aimed at. For larger lists
// the time/space usage may be sub-optimal.
//
// A pair of lists may have multiple longest common subsequences. In that
// case, the one selected by this function is undefined.
func longestCommonSubsequence(xs, ys []cty.Value) []cty.Value {
	if len(xs) == 0 || len(ys) == 0 {
		return make([]cty.Value, 0)
	}

	c := make([]int, len(xs)*len(ys))
	eqs := make([]bool, len(xs)*len(ys))
	w := len(xs)

	for y := 0; y < len(ys); y++ {
		for x := 0; x < len(xs); x++ {
			eqV := xs[x].Equals(ys[y])
			eq := false
			if eqV.IsKnown() && eqV.True() {
				eq = true
				eqs[(w*y)+x] = true // equality tests can be expensive, so cache it
			}
			if eq {
				// Sequence gets one longer than for the cell at top left,
				// since we'd append a new item to the sequence here.
				if x == 0 || y == 0 {
					c[(w*y)+x] = 1
				} else {
					c[(w*y)+x] = c[(w*(y-1))+(x-1)] + 1
				}
			} else {
				// We follow the longest of the sequence above and the sequence
				// to the left of us in the matrix.
				l := 0
				u := 0
				if x > 0 {
					l = c[(w*y)+(x-1)]
				}
				if y > 0 {
					u = c[(w*(y-1))+x]
				}
				if l > u {
					c[(w*y)+x] = l
				} else {
					c[(w*y)+x] = u
				}
			}
		}
	}

	// The bottom right cell tells us how long our longest sequence will be
	seq := make([]cty.Value, c[len(c)-1])

	// Now we will walk back from the bottom right cell, finding again all
	// of the equal pairs to construct our sequence.
	x := len(xs) - 1
	y := len(ys) - 1
	i := len(seq) - 1

	for x > -1 && y > -1 {
		if eqs[(w*y)+x] {
			// Add the value to our result list and then walk diagonally
			// up and to the left.
			seq[i] = xs[x]
			x--
			y--
			i--
		} else {
			// Take the path with the greatest sequence length in the matrix.
			l := 0
			u := 0
			if x > 0 {
				l = c[(w*y)+(x-1)]
			}
			if y > 0 {
				u = c[(w*(y-1))+x]
			}
			if l > u {
				x--
			} else {
				y--
			}
		}
	}

	if i > -1 {
		// should never happen if the matrix was constructed properly
		panic("not enough elements in sequence")
	}

	return seq
}
