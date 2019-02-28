package opentoken

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"strings"
	"time"
)

// DecodeToken is used to decode the open token into multiple key value paairs. Decode is capitalized so it's now exportable in the package namespace
func DecodeToken(token string, cipherID int, password string) (string, error) {

	// Base64 decode the otk
	data, err := base64.URLEncoding.DecodeString(AddBase64Padding(token))
	if err != nil {
		fmt.Println("error:", err)
		return "", errors.New("Could not decode the input Base 64 String properly")
	}

	fmt.Printf("data is of type %T\n", data)
	fmt.Println(data)
	fmt.Printf("%q\n", data)

	// Validate the OTK header literal
	var index = 0
	var otkHeader = string(data[index : index+3])
	fmt.Println(otkHeader)
	if otkHeader != "OTK" {
		return "", errors.New("Not a valid OTK Header")
	}

	// Validate the OTK version
	index += 3
	var otkVersion = int(data[index])
	index++
	fmt.Println(otkVersion)
	if otkVersion != 1 {
		return "", errors.New("Not a valid version " + string(otkVersion) + " must be 1.")
	}

	// Extract cipher, mac and iv information.
	var oktCipherID = int(data[index])
	index++
	fmt.Println(oktCipherID)
	var cipherName = ciphers[oktCipherID].name
	fmt.Println(cipherName)
	fmt.Println(index)
	if cipherID != oktCipherID {
		return "", errors.New("cipherId " + string(cipherID) + " argument doesn't match the encoding cipher " + string(oktCipherID))
	}

	// Extract mac and iv information.
	var hMac = data[index : index+20]
	fmt.Println(hMac)

	index += 20
	var ivlength = int(data[index])
	fmt.Println(ivlength)

	index++
	var iv []byte
	if ivlength > 0 {
		iv = data[index : index+ivlength]
		fmt.Println(iv)
		index += ivlength
	} else {
		iv = nil
	}
	fmt.Println(index)

	// Extract the Key Info (if present) and select a key for decryption.
	var keyInfo []byte
	var keyInfoLen = int(data[index])
	index = index + 1
	if keyInfoLen != 0 {
		keyInfo = data[index : index+keyInfoLen]
		index += keyInfoLen
		fmt.Println(keyInfo)
	}

	// Decrypt the payload cipher-text using the selected cipher suite
	var payloadCipherText []byte
	var payloadLength = int(binary.BigEndian.Uint16(data[index : index+2]))
	index += 2
	fmt.Println(payloadLength)
	payloadCipherText = data[index : index+payloadLength]
	fmt.Println(payloadCipherText)

	var decryptionKey = generateKey(password, nil, cipherID)
	fmt.Printf("decryptionKey is of type %T\n", decryptionKey)
	fmt.Println(decryptionKey)

	var block cipher.Block

	if cipherName == "aes-128-cbc" || cipherName == "aes-256-cbc" {
		block, err = aes.NewCipher(decryptionKey)
		if err != nil {
			panic(err)
		}
	}
	if cipherName == "3des" {
		block, err = des.NewTripleDESCipher(decryptionKey)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("***********************")
	fmt.Printf("block is of type %T\n", block)
	fmt.Println("***********************")
	fmt.Println(block)
	fmt.Println("***********************")
	mode := cipher.NewCBCDecrypter(block, iv)
	fmt.Println(mode)
	fmt.Println("***********************")

	origData := make([]byte, len(payloadCipherText))
	mode.CryptBlocks(origData, payloadCipherText)
	fmt.Println(origData)

	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return "", err
	}
	fmt.Println(origData)

	var zd1 = append(origData)
	var zdb = append(zd1, payloadCipherText...)
	// fmt.Println(zdb)

	// b := bytes.NewReader(origData)
	b := bytes.NewReader(zdb)
	z, err := zlib.NewReader(b)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer z.Close()
	decodedPayload, err := ioutil.ReadAll(z)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(decodedPayload))

	var hmacDigest = InitializeMac(oktCipherID, cipherID, decryptionKey, iv, keyInfo, decodedPayload)
	if bytes.Compare(hmacDigest, hMac) != 0 {
		return "", errors.New("HMAC does not match")
	}

	return string(decodedPayload), nil
}

// InitializeMac is used to check the integrity of the decoded message by comparing the hmacTest with hmac received in the opentoken
func InitializeMac(oktCipherID, cipherID int, decryptionKey, iv, keyInfo, decodedPayload []byte) []byte {

	// Initialize Hmac and check the integrity of the decoded message
	var hmacTest hash.Hash
	if oktCipherID == 0 {
		hmacTest = sha1.New()
	} else {
		hmacTest = hmac.New(sha1.New, decryptionKey)
	}
	// _otkVersion := []byte(strconv.Itoa(otkVersion))
	hmacTest.Write([]byte{1})

	// _oktCipherID := []byte(strconv.Itoa(oktCipherID))
	if cipherID == 0 {
		hmacTest.Write([]byte{0})
	} else if cipherID == 1 {
		hmacTest.Write([]byte{1})
	} else if cipherID == 2 {
		hmacTest.Write([]byte{2})
	} else {
		hmacTest.Write([]byte{3})
	}

	if len(iv) > 0 {
		hmacTest.Write(iv)
	}

	if len(keyInfo) > 0 {
		hmacTest.Write(keyInfo)
	}
	hmacTest.Write(decodedPayload)

	var hmacDigest = hmacTest.Sum(nil)
	fmt.Println(hmacDigest)
	return hmacDigest
}

// MakeStringBase64Compatible is used to make a string base64 compatibele by replacing '*' with '='
func MakeStringBase64Compatible(token *string) {
	// Replace trailing "*" pad characters with standard Base64 "=" characters
	*token = strings.Replace(*token, "*", "=", -1)
}

// PKCS5UnPadding is used to to PKCS5 unpadding
// PKCS5UnPadding unpad data after AES or 3DES block
// decrypting.
func PKCS5UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length <= 0 {
		return nil, fmt.Errorf("invalid byte blob lenght: expecting > 0 having %d", length)
	}
	unpadding := int(src[length-1])
	delta := length - unpadding
	if delta < 0 {
		return nil, fmt.Errorf("invalid padding delta lenght: expecting >= 0 having %d", delta)
	}
	return src[:delta], nil
}

// AddBase64Padding is used to add Base64 Padding if the length of the string is not a multiple of 4
func AddBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}
	return value
}

// ProcessPayload is used to process the payload.
func ProcessPayload(message string) (map[string]string, error) {
	keyValuePairs := make(map[string]string)
	var splitedMessage = strings.Split(message, "\n")
	fmt.Println(splitedMessage)

	for index := 0; index < len(splitedMessage); index++ {
		splitSinglePair := strings.Split(splitedMessage[index], "=")
		keyValuePairs[splitSinglePair[0]] = splitSinglePair[1]
	}
	fmt.Println(keyValuePairs)
	if (len(keyValuePairs["subject"])) == 0 {
		return nil, errors.New("OpenToken missing 'subject'")
	}

	var now = time.Now()
	var tolerance = 120 * 1000
	layout := "yyyy-MM-dd'T'HH:mm:ss'Z'"
	var notBefore, _ = time.Parse(layout, keyValuePairs["not-before"])
	var notOnOrAfter, _ = time.Parse(layout, keyValuePairs["not-on-or-after"])
	var renewUntil, _ = time.Parse(layout, keyValuePairs["renew-until"])

	if notBefore.UnixNano() > notOnOrAfter.UnixNano() {
		return nil, errors.New("'not-on-or-after' should be above 'not-before'")
	}

	if notBefore.UnixNano() > now.UnixNano() && notBefore.UnixNano() > int64(tolerance) {
		return nil, errors.New("Must not use this token before " + notBefore.String())
	}

	if now.UnixNano() > notOnOrAfter.UnixNano() {
		return nil, errors.New("Must not use this token before " + notOnOrAfter.String())
	}

	if now.UnixNano() > renewUntil.UnixNano() {
		return nil, errors.New("Must not use this token before " + renewUntil.String())
	}

	return keyValuePairs, nil
}
