package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

// ` will preserve linebreaks and tabs
const usage = `
usage:
	hmac sign|verify <key> <value>
`

func main() {
	if len(os.Args) < 4 ||
		(os.Args[1] != "sign" && os.Args[1] != "verify") {
		fmt.Println(usage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	key := os.Args[2]
	value := os.Args[3]

	switch cmd {
	case "sign":
		v := []byte(value) // converted here because we use the byte slice twice
		h := hmac.New(sha256.New, []byte(key))
		h.Write(v)
		signature := h.Sum(nil)

		buffer := make([]byte, len(v)+len(signature))
		copy(buffer, v)                          // copies v into buffer
		copy(buffer[len(v):], []byte(signature)) // copies signature to end of buffer
		fmt.Println(base64.URLEncoding.EncodeToString(buffer))

	case "verify":
		buffer, err := base64.URLEncoding.DecodeString(value)
		if err != nil {
			fmt.Printf("error decoding: %v\n", err)
			os.Exit(1)
		}
		value := buffer[:len(buffer)-sha256.Size]
		signature := buffer[len(buffer)-sha256.Size:]

		h := hmac.New(sha256.New, []byte(key))
		h.Write(value)
		signature2 := h.Sum(nil)
		/* User sends us value/signature, we check if the value encoded by our key equals the signature, if so then they are
		 * the person they say they are
		 */
		if hmac.Equal(signature, signature2) { // Constant time comparison so no timing attacks, will take same time to verify regardless of how close it is
			fmt.Println("signature is valid")
		} else {
			fmt.Println("invalid signature")
		}

	}

}
