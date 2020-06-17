package json

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"
)

func TestUint64_MarshalJSON(t *testing.T) {
	for _, v := range []Uint64{0, 1234, math.MaxUint64} {
		data1, err1 := json.Marshal(v)

		data2, err2 := json.Marshal(uint64(v))
		data2 = []byte(`"` + string(data2) + `"`)

		if !bytes.Equal(data1, data2) || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
			t.Errorf("failed, value: %v, have(%s, %v), want(%s, %v)", v, data1, err1, data2, err2)
			return
		}
		t.Log(string(data1), err1)
	}
}

func TestUint64_UnmarshalJSON(t *testing.T) {
	// 不带引号
	{
		type T1 struct {
			X Uint64 `json:"x"`
		}
		type T2 struct {
			X uint64 `json:"x"`
		}
		tests := [][]byte{
			[]byte(`{"x":null}`),
			[]byte(`{"x":0}`),
			[]byte(`{"x":1234}`),
			[]byte(`{"x":18446744073709551615}`),
			[]byte(`{"x":18446744073709551616}`),
			[]byte(`{"x":-1234}`),
		}
		for _, data := range tests {
			var v1 T1
			v1.X = 100
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			v2.X = 100
			err2 := json.Unmarshal(data, &v2)

			if uint64(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
	// 带引号
	{
		type T1 struct {
			X Uint64 `json:"x"`
		}
		type T2 struct {
			X uint64 `json:"x,string"`
		}
		tests := [][]byte{
			[]byte(`{"x":"0"}`),
			[]byte(`{"x":"1234"}`),
			[]byte(`{"x":"18446744073709551615"}`),
			[]byte(`{"x":"18446744073709551616"}`),
			[]byte(`{"x":"-1234"}`),
			[]byte(`{"x":""}`),
			[]byte(`{"x":"abc"}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if uint64(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
}
