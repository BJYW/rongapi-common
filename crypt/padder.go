package crypt

import "bytes"

type pkcs5Padder struct{}

func (p *pkcs5Padder) unpadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (p *pkcs5Padder) padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

type noPadder struct{}

func (p *noPadder) unpadding(origData []byte) []byte {
	return origData
}

func (p *noPadder) padding(ciphertext []byte, blockSize int) []byte {
	return ciphertext
}
