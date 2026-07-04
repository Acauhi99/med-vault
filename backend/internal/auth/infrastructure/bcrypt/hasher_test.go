package bcrypt

import "testing"

func TestHashAndCompare(t *testing.T) {
	h := NewHasher(12)

	hash, err := h.Hash("password12345")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if hash == "" {
		t.Fatal("empty hash")
	}

	if err := h.Compare(hash, "password12345"); err != nil {
		t.Errorf("compare should succeed: %v", err)
	}

	if err := h.Compare(hash, "wrongpassword"); err == nil {
		t.Error("compare should fail for wrong password")
	}
}
