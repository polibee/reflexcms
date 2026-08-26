package adminhub

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func cryptoRandRead(b []byte) (int, error) {
	return rand.Read(b)
}

func randHex(n int) string {
	raw := make([]byte, n)
	if _, err := rand.Read(raw); err != nil {
		return time.Now().Format("150405.000000000")
	}
	return hex.EncodeToString(raw)
}

func timeNowHex() string {
	return time.Now().Format("150405.000000000")
}
