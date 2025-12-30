package core_test

import (
	"fmt"
	"testing"

	"github.com/gurshaan17/redis/core"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}

	for k, v := range cases {
		value, err := core.Decode([]byte(k))
		if err != nil {
			t.Fatal(err)
		}
		if value != v {
			t.Fatalf("expected %v, got %v", v, value)
		}
	}
}

func TestBulkStringDecode(t *testing.T) {
	cases := map[string]string{
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n":      "",
	}

	for k, v := range cases {
		value, err := core.Decode([]byte(k))
		if err != nil {
			t.Fatal(err)
		}
		if value != v {
			t.Fatalf("expected %v, got %v", v, value)
		}
	}
}

func TestArrayDecode(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n": {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n": {
			"hello", "world",
		},
		"*3\r\n:1\r\n:2\r\n:3\r\n": {
			int64(1), int64(2), int64(3),
		},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+hello\r\n-world\r\n": {
			[]interface{}{int64(1), int64(2), int64(3)},
			[]interface{}{"hello", "world"},
		},
	}

	for k, v := range cases {
		value, err := core.Decode([]byte(k))
		if err != nil {
			t.Fatal(err)
		}

		array := value.([]interface{})
		if len(array) != len(v) {
			t.Fatalf("length mismatch")
		}

		for i := range array {
			if fmt.Sprintf("%v", array[i]) != fmt.Sprintf("%v", v[i]) {
				t.Fatalf("expected %v, got %v", v[i], array[i])
			}
		}
	}
}
