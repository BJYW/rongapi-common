package crypt

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func RsaSignWithPrivateKey(origData, private_key_bts []byte) (bts []byte, err error) {
	var private *rsa.PrivateKey
	if block, _ := pem.Decode(private_key_bts); block != nil {
		var priInterface interface{}
		if priInterface, err = x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			private = priInterface.(*rsa.PrivateKey)
		}
	} else {
		err = errors.New("private key error")
	}
	resbyte := md5.Sum(origData)
	return rsa.SignPKCS1v15(nil, private, crypto.MD5, resbyte[:])
}
