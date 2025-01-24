package crypto

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestNewEncryptionKey(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "Generate valid encryption key",
			wantLen: 32, // AES-256 key length
			wantErr: false,
		},
	}

	// Run multiple times to ensure consistent behavior
	for i := 0; i < 100; i++ {
		key1, err1 := NewEncryptionKey()
		key2, err2 := NewEncryptionKey()

		if err1 != nil || err2 != nil {
			t.Errorf("NewEncryptionKey() failed to generate keys: %v, %v", err1, err2)
		}

		// Keys should be unique
		if bytes.Equal(key1, key2) {
			t.Error("NewEncryptionKey() generated duplicate keys")
		}

		// Keys should be the correct length
		if len(key1) != 32 || len(key2) != 32 {
			t.Errorf("NewEncryptionKey() generated keys with wrong length: %d, %d", len(key1), len(key2))
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewEncryptionKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptionKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(key) != tt.wantLen {
				t.Errorf("NewEncryptionKey() key length = %v, want %v", len(key), tt.wantLen)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	tests := []struct {
		name      string
		plaintext []byte
		wantErr   bool
	}{
		{
			name:      "Encrypt and decrypt small text",
			plaintext: []byte("Hello, World!"),
			wantErr:   false,
		},
		{
			name:      "Encrypt and decrypt empty text",
			plaintext: []byte(""),
			wantErr:   false,
		},
		{
			name:      "Encrypt and decrypt large text",
			plaintext: []byte(strings.Repeat("Large text with lots of data. ", 100)),
			wantErr:   false,
		},
		{
			name:      "Encrypt and decrypt binary data",
			plaintext: []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			wantErr:   false,
		},
		{
			name:      "Encrypt and decrypt JSON data",
			plaintext: []byte(`{"key":"value","nested":{"array":[1,2,3]}}`),
			wantErr:   false,
		},
		{
			name:      "Encrypt and decrypt with special characters",
			plaintext: []byte("Special chars: !@#$%^&*()_+-=[]{}|;:,.<>?"),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the plaintext
			encrypted, err := Encrypt(key, tt.plaintext)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// Verify the encrypted text is base64 encoded
			if _, err := base64.StdEncoding.DecodeString(encrypted); err != nil {
				t.Errorf("Encrypted text is not valid base64: %v", err)
			}

			// Verify encrypted text is different from plaintext
			if string(tt.plaintext) == encrypted {
				t.Error("Encrypted text matches plaintext")
			}

			// Decrypt the ciphertext
			decrypted, err := Decrypt(key, encrypted)
			if err != nil {
				t.Fatalf("Failed to decrypt: %v", err)
			}

			// Compare the decrypted text with original plaintext
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Decrypted text does not match original\nGot: %v\nWant: %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptionErrors(t *testing.T) {
	validKey, _ := NewEncryptionKey()
	invalidKey := make([]byte, 16) // Wrong key size

	tests := []struct {
		name        string
		key         EncryptionKey
		plaintext   []byte
		ciphertext  string
		wantEncErr  bool
		wantDecErr  bool
		errContains string
	}{
		{
			name:        "Invalid key size for encryption",
			key:         invalidKey,
			plaintext:   []byte("test"),
			wantEncErr:  true,
			errContains: "invalid key size",
		},
		{
			name:        "Invalid key size for decryption",
			key:         invalidKey,
			ciphertext:  "invalid",
			wantDecErr:  true,
			errContains: "invalid key size",
		},
		{
			name:        "Invalid base64 for decryption",
			key:         validKey,
			ciphertext:  "invalid base64!@#$",
			wantDecErr:  true,
			errContains: "failed to decode base64",
		},
		{
			name:        "Malformed ciphertext",
			key:         validKey,
			ciphertext:  base64.StdEncoding.EncodeToString([]byte("too short")),
			wantDecErr:  true,
			errContains: "malformed ciphertext",
		},
		{
			name:        "Tampered ciphertext",
			key:         validKey,
			ciphertext:  base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x00}, 32)),
			wantDecErr:  true,
			errContains: "failed to decrypt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.plaintext != nil {
				_, err := Encrypt(tt.key, tt.plaintext)
				if !tt.wantEncErr && err != nil {
					t.Errorf("Encrypt() unexpected error: %v", err)
				}
				if tt.wantEncErr && err == nil {
					t.Error("Encrypt() expected error but got none")
				}
				if tt.wantEncErr && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Encrypt() error = %v, want error containing %v", err, tt.errContains)
				}
			}

			if tt.ciphertext != "" {
				_, err := Decrypt(tt.key, tt.ciphertext)
				if !tt.wantDecErr && err != nil {
					t.Errorf("Decrypt() unexpected error: %v", err)
				}
				if tt.wantDecErr && err == nil {
					t.Error("Decrypt() expected error but got none")
				}
				if tt.wantDecErr && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Decrypt() error = %v, want error containing %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestEncryptionKeyReuse(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test that the same key can be reused for multiple encryptions
	plaintext1 := []byte("First message")
	plaintext2 := []byte("Second message")

	encrypted1, err := Encrypt(key, plaintext1)
	if err != nil {
		t.Fatalf("Failed to encrypt first message: %v", err)
	}

	encrypted2, err := Encrypt(key, plaintext2)
	if err != nil {
		t.Fatalf("Failed to encrypt second message: %v", err)
	}

	// Verify that encrypting the same plaintext twice produces different ciphertexts
	if encrypted1 == encrypted2 {
		t.Error("Encrypting same plaintext twice produced identical ciphertexts")
	}

	// Verify both messages can be decrypted correctly
	decrypted1, err := Decrypt(key, encrypted1)
	if err != nil {
		t.Fatalf("Failed to decrypt first message: %v", err)
	}

	decrypted2, err := Decrypt(key, encrypted2)
	if err != nil {
		t.Fatalf("Failed to decrypt second message: %v", err)
	}

	if !bytes.Equal(decrypted1, plaintext1) {
		t.Error("First decrypted message does not match original")
	}

	if !bytes.Equal(decrypted2, plaintext2) {
		t.Error("Second decrypted message does not match original")
	}
}

func TestEncryptDecryptShare(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	type Share struct {
		Index     int      `json:"index"`
		Value     []byte   `json:"value"`
		Metadata  string   `json:"metadata"`
		Timestamp int64    `json:"timestamp"`
		Tags      []string `json:"tags"`
	}

	share := Share{
		Index:     1,
		Value:     []byte{0x01, 0x02, 0x03},
		Metadata:  "test share",
		Timestamp: 1234567890,
		Tags:      []string{"tag1", "tag2"},
	}

	// Marshal share to JSON
	shareBytes, err := json.Marshal(share)
	if err != nil {
		t.Fatalf("Failed to marshal share: %v", err)
	}

	// Encrypt share
	encrypted, err := Encrypt(key, shareBytes)
	if err != nil {
		t.Fatalf("Failed to encrypt share: %v", err)
	}

	// Decrypt share
	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt share: %v", err)
	}

	// Unmarshal decrypted share
	var decryptedShare Share
	if err := json.Unmarshal(decrypted, &decryptedShare); err != nil {
		t.Fatalf("Failed to unmarshal decrypted share: %v", err)
	}

	// Compare original and decrypted shares
	if !reflect.DeepEqual(share, decryptedShare) {
		t.Errorf("Decrypted share does not match original\nGot: %+v\nWant: %+v", decryptedShare, share)
	}
}
