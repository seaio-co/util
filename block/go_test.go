// combination.go

package hmath

import (
	"math/rand"
	"testing"
)

func combination(m, n int) int {
	if n > m-n {
		n = m - n
	}

	c := 1
	for i := 0; i < n; i++ {
		c *= m - i
		c /= i + 1
	}

	return c
}

// combination_test.go

package hmath

import (
"math/rand"
"testing"
)


func TestCombination(t *testing.T) {

	for _, unit := range []struct {
		m        int
		n        int
		expected int
	}{
		{1, 0, 1},
		{4, 1, 4},
		{4, 2, 6},
		{4, 3, 4},
		{4, 4, 1},
		{10, 1, 10},
		{10, 3, 120},
		{10, 7, 120},
	} {

		if actually := combination(unit.m, unit.n); actually != unit.expected {
			t.Errorf("combination: [%v], actually: [%v]", unit, actually)
		}
	}
}


func BenchmarkCombination(b *testing.B) {

	for i := 0; i < b.N; i++ {
		combination(i+1, rand.Intn(i+1))
	}
}


func BenchmarkCombinationParallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := rand.Intn(100) + 1
			n := rand.Intn(m)
			combination(m, n)
		}
	})
}