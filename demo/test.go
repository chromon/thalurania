package main

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"log"
)

func main() {
	salt := []byte{0xe8, 0x9d, 0xb2, 0x77, 0xc3, 0xfe, 0xa5, 0xde}

	dk, err := scrypt.Key([]byte("dddddd"), salt, 1<<15, 5, 3, 32)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(dk))
	fmt.Println(string(dk[:]))
}
