package uuid

import (
	"testing"
)

func TestNewUUID(t *testing.T) {
	id := New()
	t.Logf("new uuid: %s", id)
}

func TestVerifyUUID(t *testing.T) {
	ok := VerifyUUID("ac690c4f-74a5-4a13-8b3e-98eb42b358ff")
	if !ok {
		t.Fatal("invalid uuid string")
	}
	t.Logf("verify: %v", ok)
}

func BenchmarkNewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := New()
		if id == "" {
			b.Fatal("new uuid is empty")
		}
	}
}
