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

func Test_SliceUniqueString(t *testing.T) {

	in := make([]string, 0)
	in = append(in, "hello", "world")
	t.Log(SliceUniqueString(in))
}

func Test_SliceMerge(t *testing.T) {

	in := make([]interface{}, 0)
	in = append(in, 2, 8, 3)
	in1 := make([]interface{}, 0)
	in1 = append(in, 2, 8, 3)
	t.Log(SliceMerge(in, in1))
}

func Test_SliceSumInt64(t *testing.T) {

	in := make([]int64, 0)
	in = append(in, 2, 8, 3, 5, 7, 0)
	t.Log(SliceSumInt64(in))
}
