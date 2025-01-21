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
			// Generate a new key for each test
			key, err := NewEncryptionKey()
			if err != nil {
				t.Fatalf("Failed to generate key: %v", err)
			}

			// Encrypt the plaintext
			encrypted, err := Encrypt(key, tt.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the encrypted text is base64 encoded
			if _, err := base64.StdEncoding.DecodeString(encrypted); err != nil {
				t.Errorf("Encrypted text is not valid base64: %v", err)
			}

			// Decrypt the ciphertext
			decrypted, err := Decrypt(key, encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the decrypted text with original plaintext
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptDecryptShare(t *testing.T) {
	// Generate a test data structure
	type TestShareData struct {
		ShareID string
		Data    []byte
		Meta    map[string]string
	}

	testData := &TestShareData{
		ShareID: "test-share-1",
		Data:    []byte("test share data"),
		Meta: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Generate encryption key
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatalf("Failed to generate encryption key: %v", err)
	}

	// Marshal test data to JSON
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// Encrypt the data
	encrypted, err := Encrypt(key, jsonData)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Decrypt the data
	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Unmarshal and verify
	var decryptedData TestShareData
	if err := json.Unmarshal(decrypted, &decryptedData); err != nil {
		t.Fatalf("Failed to unmarshal decrypted data: %v", err)
	}

	// Verify the decrypted data matches original
	if decryptedData.ShareID != testData.ShareID {
		t.Errorf("ShareID mismatch: got %v, want %v", decryptedData.ShareID, testData.ShareID)
	}
	if !bytes.Equal(decryptedData.Data, testData.Data) {
		t.Errorf("Data mismatch: got %v, want %v", decryptedData.Data, testData.Data)
	}
	if !reflect.DeepEqual(decryptedData.Meta, testData.Meta) {
		t.Errorf("Meta mismatch: got %v, want %v", decryptedData.Meta, testData.Meta)
	}
}

func TestEncryptionErrors(t *testing.T) {
	tests := []struct {
		name        string
		key         EncryptionKey
		ciphertext  string
		shouldError bool
	}{
		{
			name:        "Invalid key length",
			key:         make([]byte, 16), // Too short for AES-256
			ciphertext:  "",
			shouldError: true,
		},
		{
			name:        "Invalid base64 ciphertext",
			key:         make([]byte, 32),
			ciphertext:  "invalid base64!@#$",
			shouldError: true,
		},
		{
			name:        "Ciphertext too short",
			key:         make([]byte, 32),
			ciphertext:  base64.StdEncoding.EncodeToString([]byte("too short")),
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Try to decrypt with invalid input
			_, err := Decrypt(tt.key, tt.ciphertext)
			if (err != nil) != tt.shouldError {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.shouldError)
			}
		})
	}
}

func TestEncryptionKeyReuse(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test multiple encryptions with the same key
	plaintexts := [][]byte{
		[]byte("First message"),
		[]byte("Second message"),
		[]byte("Third message"),
	}

	for i, plaintext := range plaintexts {
		encrypted, err := Encrypt(key, plaintext)
		if err != nil {
			t.Errorf("Encryption %d failed: %v", i, err)
			continue
		}

		decrypted, err := Decrypt(key, encrypted)
		if err != nil {
			t.Errorf("Decryption %d failed: %v", i, err)
			continue
		}

		if !bytes.Equal(decrypted, plaintext) {
			t.Errorf("Test %d: got %v, want %v", i, decrypted, plaintext)
		}
	}
}
