package crypt

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

const defaultKey = "wwwidtagcn"
const defaultIV = "WWWID5CNWWWID5CN"

var defaultPadder = pkcs5Padder{}

//EncryptPwd 加密密码
func EncryptPwd(password string, ivs ...string) (string, error) {
	encrypted_password := md5String(password)
	var iv = defaultIV
	if len(ivs) > 0 {
		iv = ivs[0]
	}
	return Encrypt(encrypted_password, password, iv)
}

//Encrypt 加密字符串
func Encrypt(strTobeEnCrypted string, args ...string) (string, error) {
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
	if block, err := des.NewTripleDESCipher(key); err == nil {
		src := defaultPadder.padding([]byte(strTobeEnCrypted), block.BlockSize())
		crypted := make([]byte, len(src))
		blockMode := cipher.NewCBCEncrypter(block, hexMustDecode(iv))
		blockMode.CryptBlocks(crypted, src)
		return base64.StdEncoding.EncodeToString(crypted), nil
	} else {
		return "", err
	}
}

func md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	bts := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(bts))
}

func paddingKey(encryptionKey string, keylength int, appStr string) string {

	if encryptionKey == "" || len(encryptionKey) == 0 {
		return encryptionKey
	}

	strkeylength := len(encryptionKey)

	if strkeylength-keylength > 0 {
		encryptionKey = encryptionKey[:47]
	}
	size := keylength - strkeylength
	for i := 0; i < size; i++ {
		encryptionKey += appStr
	}
	return encryptionKey
}

func hexMustDecode(str string) []byte {
	bts := make([]byte, len(str)/2)
	str = strings.ToLower(str)
	for index := 0; index < len(str); index = index + 2 {
		content := []byte(str[index : index+2])
		i := index / 2
		if content[0] < 'a' {
			bts[i] = byte(content[0]-'0') << 4
		} else {
			bts[i] = (content[0] - 'a' + 10) << 4
		}
		if content[1] < 'a' {
			bts[i] += content[1] - '0'
		} else {
			bts[i] += content[1] - 'a' + 10
		}
	}
	return bts
}
