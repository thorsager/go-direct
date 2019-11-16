package utl

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

func AsBase64String(number uint64) string {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, number)
	return base64.RawURLEncoding.EncodeToString(buffer)
}

func AsUint64(uint64base string) (uint64, error) {
	if len(uint64base) != 11 {
		return 0, fmt.Errorf("invalid string length")
	}
	buffer, err := base64.RawURLEncoding.DecodeString(uint64base)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buffer), nil
}
