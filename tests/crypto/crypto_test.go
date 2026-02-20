package crypto_test

import (
	"bytes"
	"testing"

	"GitSyncer/core/crypto"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}

	if len(salt1) != crypto.SaltLen {
		t.Fatalf("GenerateSalt() length = %d, want %d", len(salt1), crypto.SaltLen)
	}

	salt2, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() second call error: %v", err)
	}

	if bytes.Equal(salt1, salt2) {
		t.Fatal("GenerateSalt() produced identical salts on two calls")
	}
}

func TestDeriveKey(t *testing.T) {
	salt, _ := crypto.GenerateSalt()

	key := crypto.DeriveKey("test-password", salt)
	if len(key) != 32 {
		t.Fatalf("DeriveKey() length = %d, want 32", len(key))
	}

	// Same password + salt = same key
	key2 := crypto.DeriveKey("test-password", salt)
	if !bytes.Equal(key, key2) {
		t.Fatal("DeriveKey() with same inputs produced different keys")
	}

	// Different password = different key
	key3 := crypto.DeriveKey("different-password", salt)
	if bytes.Equal(key, key3) {
		t.Fatal("DeriveKey() with different password produced same key")
	}

	// Different salt = different key
	salt2, _ := crypto.GenerateSalt()
	key4 := crypto.DeriveKey("test-password", salt2)

	if bytes.Equal(key, key4) {
		t.Fatal("DeriveKey() with different salt produced same key")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := crypto.DeriveKey("my-password", []byte("1234567890123456"))

	tests := []struct {
		name      string
		plaintext string
	}{
		{"short string", "ghp_abc123token"},
		{"empty string", ""},
		{"long string", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7...a-very-long-key-value-here"},
		{"unicode", "p\u00e4ssw\u00f6rd-with-\u00fcn\u00efc\u00f6de"},
		{"special chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := crypto.Encrypt([]byte(tt.plaintext), key)
			if err != nil {
				t.Fatalf("Encrypt() error: %v", err)
			}

			// Ciphertext should differ from plaintext
			if tt.plaintext != "" && bytes.Equal(ciphertext, []byte(tt.plaintext)) {
				t.Fatal("Encrypt() ciphertext equals plaintext")
			}

			decrypted, err := crypto.Decrypt(ciphertext, key)
			if err != nil {
				t.Fatalf("Decrypt() error: %v", err)
			}

			if string(decrypted) != tt.plaintext {
				t.Fatalf("Decrypt() = %q, want %q", string(decrypted), tt.plaintext)
			}
		})
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key1 := crypto.DeriveKey("correct-password", []byte("1234567890123456"))
	key2 := crypto.DeriveKey("wrong-password", []byte("1234567890123456"))

	ciphertext, err := crypto.Encrypt([]byte("secret-token"), key1)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	_, err = crypto.Decrypt(ciphertext, key2)
	if err == nil {
		t.Fatal("Decrypt() with wrong key should return error")
	}

	if err != crypto.ErrDecryptionFailed {
		t.Fatalf("Decrypt() error = %v, want ErrDecryptionFailed", err)
	}
}

func TestDecryptCiphertextTooShort(t *testing.T) {
	key := crypto.DeriveKey("password", []byte("1234567890123456"))

	_, err := crypto.Decrypt([]byte("short"), key)
	if err == nil {
		t.Fatal("Decrypt() with short ciphertext should return error")
	}
}

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	key := crypto.DeriveKey("password", []byte("1234567890123456"))
	plaintext := []byte("same-input")

	ct1, _ := crypto.Encrypt(plaintext, key)
	ct2, _ := crypto.Encrypt(plaintext, key)

	if bytes.Equal(ct1, ct2) {
		t.Fatal("Encrypt() should produce different ciphertexts for same input (random nonce)")
	}
}

func TestHashPasswordAndVerify(t *testing.T) {
	salt, _ := crypto.GenerateSalt()
	password := "master-password-123"

	hash := crypto.HashPassword(password, salt)
	if len(hash) != 32 {
		t.Fatalf("HashPassword() length = %d, want 32", len(hash))
	}

	// Correct password verifies
	if !crypto.VerifyPassword(password, salt, hash) {
		t.Fatal("VerifyPassword() returned false for correct password")
	}

	// Wrong password does not verify
	if crypto.VerifyPassword("wrong-password", salt, hash) {
		t.Fatal("VerifyPassword() returned true for wrong password")
	}
}

func TestZeroBytes(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	crypto.ZeroBytes(data)

	for i, b := range data {
		if b != 0 {
			t.Fatalf("ZeroBytes() did not zero byte at index %d: got %d", i, b)
		}
	}
}
