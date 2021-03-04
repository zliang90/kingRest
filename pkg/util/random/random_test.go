package random

import (
	"testing"
)

func TestRandString(t *testing.T) {
	s, err := RandString(18, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(s) != 18 {
		t.Errorf("mismatch random string length")
	}
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := RandString(32, "")
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
		//fmt.Println(n)
	}
}
