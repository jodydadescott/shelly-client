package types

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
)

func getSHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

// getCnonce returns a random client nonce; Shelly's auth object takes it as a number.
func getCnonce() uint32 {
	b := make([]byte, 4)
	io.ReadFull(rand.Reader, b)
	return binary.BigEndian.Uint32(b)
}
