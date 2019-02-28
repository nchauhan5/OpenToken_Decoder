package opentoken

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

type cipherTypes struct {
	id       int
	name     string
	cipher   string
	keysize  int
	mode     string
	padding  string
	ivlength int
}

var ciphers = []cipherTypes{
	{
		id:       0,
		name:     "",
		cipher:   "",
		keysize:  0,
		mode:     "",
		padding:  "",
		ivlength: 0,
	},
	{
		id:       1,
		name:     "aes-256-cbc",
		cipher:   "AES",
		keysize:  256,
		mode:     "CBC",
		padding:  "PKCS 5",
		ivlength: 16,
	},
	{
		id:       2,
		name:     "aes-128-cbc",
		cipher:   "AES",
		keysize:  128,
		mode:     "CBC",
		padding:  "PKCS 5",
		ivlength: 16,
	},
	{
		id:       3,
		name:     "3des",
		cipher:   "3DES",
		keysize:  168,
		mode:     "CBC",
		padding:  "PKCS 5",
		ivlength: 8,
	},
}

func generateKey(password string, salt []byte, cipherID int) []byte {

	var passwordBytes = []byte(password)
	fmt.Println(passwordBytes)

	// decodedPass, _ := base64.URLEncoding.DecodeString(password)
	// fmt.Println(decodedPass)

	var derivedKey []byte
	if cipherID == 1 || cipherID == 2 {
		salt = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
		fmt.Println("The SALT is *********** ")
		fmt.Println(salt)
		if cipherID < 0 || cipherID > len(ciphers) {
			return nil
		}

		var cipher = ciphers[cipherID]
		var iterations = 1000
		derivedKey = pbkdf2.Key(passwordBytes, salt, iterations, cipher.keysize/8, sha1.New)
	}
	if cipherID == 3 {
		ede2Key := passwordBytes
		derivedKey = append(derivedKey, ede2Key[:16]...)
		derivedKey = append(derivedKey, ede2Key[:8]...)
	}
	return derivedKey
}
