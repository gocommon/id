package id

import "testing"

func Test_ID(t *testing.T) {

	var pre int64
	for i := 0; i < 100000; i++ {
		id := MustID()
		t.Log(id)
		if id < pre {
			t.Fatal("id < pre")
		}
		pre = id
	}
}
