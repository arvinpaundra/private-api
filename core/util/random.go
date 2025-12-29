package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"time"
)

const (
	AlphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	LowercaseCharset    = "abcdefghijklmnopqrstuvwxyz"
	UppercaseCharset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumericCharset      = "0123456789"
)

func RandomAlphanumeric(length int) (string, error) {
	// Get current timestamp in nanoseconds for additional entropy
	timestamp := time.Now().UnixNano()
	timestampBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(timestamp))

	// Generate random bytes
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Combine timestamp and random bytes, then hash to mix entropy
	combined := append(timestampBytes, randomBytes...)
	hash := sha256.Sum256(combined)

	// Use the hash to select characters from charset
	result := make([]byte, length)
	for i := range result {
		// Use different parts of the hash for each character position
		hashIndex := i % len(hash)
		charsetIndex := int(hash[hashIndex]) % len(AlphanumericCharset)
		result[i] = AlphanumericCharset[charsetIndex]

		// Re-hash if we've used all hash bytes to get more entropy
		if (i+1)%len(hash) == 0 && i < length-1 {
			hash = sha256.Sum256(hash[:])
		}
	}

	return string(result), nil
}
