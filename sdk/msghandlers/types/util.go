package types

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func getSHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}
