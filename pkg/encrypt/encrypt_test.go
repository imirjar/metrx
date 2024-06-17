package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestEncryptSHA256(t *testing.T) {
	tests := []struct {
		name  string
		value string
		key   string
		want  string
	}{
		{
			name:  "Empty strings",
			value: "",
			key:   "",
			want:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:  "Non-empty value and empty key",
			value: "test",
			key:   "",
			want:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		{
			name:  "Empty value and non-empty key",
			value: "",
			key:   "key",
			want:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:  "Non-empty value and key",
			value: "test",
			key:   "key",
			want:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncryptSHA256(tt.value, tt.key)
			if hex.EncodeToString(got) != tt.want {
				t.Errorf("EncryptSHA256() = %v, want %v", hex.EncodeToString(got), tt.want)
			}
		})
	}
}

func TestEncryptSHA256WithSHA256Package(t *testing.T) {
	tests := []struct {
		name  string
		value string
		key   string
	}{
		{
			name:  "Empty strings",
			value: "",
			key:   "",
		},
		{
			name:  "Non-empty value and empty key",
			value: "test",
			key:   "",
		},
		{
			name:  "Empty value and non-empty key",
			value: "",
			key:   "key",
		},
		{
			name:  "Non-empty value and key",
			value: "test",
			key:   "key",
		},
		{
			name:  "Long value and key",
			value: "this is a very long value to test the encryption function",
			key:   "this is a very long key to test the encryption function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := sha256.Sum256([]byte(tt.value))
			got := EncryptSHA256(tt.value, tt.key)
			if !compareHashes(got, want[:]) {
				t.Errorf("EncryptSHA256() = %v, want %v", hex.EncodeToString(got), hex.EncodeToString(want[:]))
			}
		})
	}
}

// compareHashes is a helper function to compare two byte slices
func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
