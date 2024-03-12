package crypt

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
)

//Decrypt 加密字符串
func Decrypt(strTobeDeCrypted string, args ...string) (string, error) {
	var strKey = defaultKey
	var iv = defaultIV
	if len(args) > 0 {
		strKey = args[0]
	}
	if len(args) > 1 {
		iv = args[1]
	}
	p := paddingKey(strKey, 64, "D")
	key := hexMustDecode(p)
	key = key[:24]
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, hexMustDecode(iv))
	crypted := make([]byte, len(strTobeDeCrypted))
	if n, err := base64.StdEncoding.Decode(crypted, []byte(strTobeDeCrypted)); err != nil {
		return "", err
	} else {
		crypted = crypted[:n]
	}
	fmt.Println(len(strTobeDeCrypted), len(crypted))
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = defaultPadder.unpadding(origData)
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

type padder interface {
	unpadding([]byte) []byte
	padding([]byte, int) []byte
}
