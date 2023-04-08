package utils

import (
	"fmt"
	"testing"
)

func TestCrypto(t *testing.T) {
	Key := `Cushion Key`
	Raw := "Hello World!"
	enc := EncodeAESWithKey(Raw, Key)
	dec := DecodeAESWithKey(enc, Key)
	fmt.Printf("Raw: %s\nEnc: %s\nDec: %s\n", Raw, enc, dec)
}
