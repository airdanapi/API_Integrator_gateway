package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordHasherUsesCostTenAndNeverReturnsPlaintext(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	hash, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("Hash() returned an unexpected error: %v", err)
	}
	if hash == "correct horse battery staple" {
		t.Fatal("Hash() returned plaintext")
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("read bcrypt cost: %v", err)
	}
	if cost != 10 {
		t.Fatalf("bcrypt cost = %d, want 10", cost)
	}
	if err := hasher.Compare(hash, "correct horse battery staple"); err != nil {
		t.Fatalf("Compare() rejected valid password: %v", err)
	}
	if err := hasher.Compare(hash, "wrong password"); err == nil {
		t.Fatal("Compare() accepted invalid password")
	}
}
