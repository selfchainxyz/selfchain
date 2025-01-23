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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the plaintext
			encrypted, err := Encrypt(key, tt.plaintext)
			if err != nil {
				t.Fatalf("Failed to encrypt: %v", err)
			}

			// Verify the encrypted text is base64 encoded
			if _, err := base64.StdEncoding.DecodeString(encrypted); err != nil {
				t.Errorf("Encrypted text is not valid base64: %v", err)
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

func TestEncryptDecryptShare(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatal("Failed to generate key:", err)
	}

	type Share struct {
		Index     int      `json:"index"`
		Value     []byte   `json:"value"`
		Threshold int      `json:"threshold"`
		Parties   []string `json:"parties"`
	}

	share := Share{
		Index:     1,
		Value:     []byte("test share value"),
		Threshold: 2,
		Parties:   []string{"party1", "party2", "party3"},
	}

	// Serialize share to JSON
	shareBytes, err := json.Marshal(share)
	if err != nil {
		t.Fatal("Failed to marshal share:", err)
	}

	// Encrypt share
	encrypted, err := Encrypt(key, shareBytes)
	if err != nil {
		t.Fatal("Failed to encrypt share:", err)
	}

	// Decrypt share
	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatal("Failed to decrypt share:", err)
	}

	// Unmarshal decrypted data
	var decryptedShare Share
	if err := json.Unmarshal(decrypted, &decryptedShare); err != nil {
		t.Fatal("Failed to unmarshal decrypted share:", err)
	}

	// Compare original and decrypted shares
	if !reflect.DeepEqual(share, decryptedShare) {
		t.Errorf("Decrypted share does not match original\nOriginal: %+v\nDecrypted: %+v", share, decryptedShare)
	}
}

func TestEncryptionErrors(t *testing.T) {
	tests := []struct {
		name        string
		key         EncryptionKey
		ciphertext  string
		expectError string
	}{
		{
			name:        "Invalid key length",
			key:         make([]byte, 31),
			ciphertext:  "test",
			expectError: "invalid key size: must be 32 bytes",
		},
		{
			name:        "Invalid base64 ciphertext",
			key:         make([]byte, 32),
			ciphertext:  "invalid base64",
			expectError: "failed to decode base64",
		},
		{
			name:        "Ciphertext too short",
			key:         make([]byte, 32),
			ciphertext:  base64.StdEncoding.EncodeToString([]byte("short")),
			expectError: "malformed ciphertext: too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.key, tt.ciphertext)
			if err == nil {
				t.Error("Expected error but got nil")
			} else if !strings.Contains(err.Error(), tt.expectError) {
				t.Errorf("Expected error containing %q but got %q", tt.expectError, err.Error())
			}
		})
	}
}

func TestEncryptionKeyReuse(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatal("Failed to generate key:", err)
	}

	plaintexts := [][]byte{
		[]byte("message 1"),
		[]byte("message 2"),
		[]byte("message 3"),
	}

	// Encrypt multiple messages with the same key
	var ciphertexts []string
	for _, plaintext := range plaintexts {
		ciphertext, err := Encrypt(key, plaintext)
		if err != nil {
			t.Fatalf("Failed to encrypt: %v", err)
		}
		ciphertexts = append(ciphertexts, ciphertext)
	}

	// Decrypt and verify each message
	for i, ciphertext := range ciphertexts {
		decrypted, err := Decrypt(key, ciphertext)
		if err != nil {
			t.Fatalf("Decryption %d failed: %v", i, err)
		}
		if !bytes.Equal(decrypted, plaintexts[i]) {
			t.Errorf("Decrypted text %d does not match original\nGot: %v\nWant: %v", i, decrypted, plaintexts[i])
		}
	}
}
