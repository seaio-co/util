package conv

import "testing"

func Test_SliceContains(t *testing.T) {

	in := make([]interface{}, 0)
	in = append(in, 2, "Go", 8, "language", 'a', false, "A", 3.14)
	t.Log(SliceContains(in, 2))
}

func Test_SliceContainsInt(t *testing.T) {

	in := make([]int, 0)
	in = append(in, 2, 8, 3)
	t.Log(SliceContainsInt(in, 2))
}
