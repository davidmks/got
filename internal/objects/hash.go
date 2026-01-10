package objects

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// HashLength is the length of a SHA-256 hash in hexadecimal characters.
const HashLength = 64

// Hash computes the SHA-256 hash of content and returns it as a hex string.
func Hash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// HashReader computes the SHA-256 hash from an io.Reader.
// This is useful for hashing large files without loading them entirely into memory.
func HashReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
