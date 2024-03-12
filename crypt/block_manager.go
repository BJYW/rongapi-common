package crypt

import (
	"crypto/cipher"
	"errors"
	"fmt"
)

type blocker interface {
	createEncrypter(block) cipher.BlockMode
	createDecrypter(block) cipher.BlockMode
}
type block interface {
	cipher.Block
	outBlockSize() int
}

type cbcBM struct {
	iv []byte
}

func (c *cbcBM) createEncrypter(b block) cipher.BlockMode {
	return cipher.NewCBCEncrypter(b, c.iv)
}

func (c *cbcBM) createDecrypter(b block) cipher.BlockMode {
	return cipher.NewCBCDecrypter(b, c.iv)
}

func (c *cbcBM) outSize(src int) int {
	return src
}

func (c *cbcBM) init(attrs map[string]interface{}) error {
	if iv, ok := attrs["iv"]; ok {
		switch ti := iv.(type) {
		case []byte:
			c.iv = ti
		case string:
			if encoding, ok := attrs["encoding_bytes"]; ok {
				if se, ok := encoding.(string); ok {
					switch se {
					case "raw":
						c.iv = []byte(ti)
					default:
						c.iv = hexMustDecode(ti)
					}
				} else {
					return errors.New("encoding bytes must be string")
				}
			} else {
				return errors.New("no encoding_bytes attr for decoding bytes from string")
			}
		}
	} else {
		return errors.New("no iv attr for des crypter")
	}
	fmt.Println(c)
	return nil
}

type noBM struct {
	block cipher.Block
}
type noBM_bm struct {
	block   cipher.Block
	outSize int
}

func (c *noBM) createEncrypter(b block) cipher.BlockMode {
	return &noBM_bm{b, b.outBlockSize()}
}

func (c *noBM) createDecrypter(b block) cipher.BlockMode {
	return &noBM_bm{b, b.outBlockSize()}
}

func (c *noBM) init(attrs map[string]interface{}) error {
	return nil
}

func (c *noBM_bm) BlockSize() int {
	return c.block.BlockSize()
}

func (c *noBM_bm) CryptBlocks(dst, src []byte) {
	temp_src := src
	handled_length := 0
	for len(temp_src) > 0 {
		var temp []byte
		max := c.BlockSize()
		if len(temp_src) > max {
			temp = temp_src[:max]
			temp_src = temp_src[max:]
		} else {
			temp = temp_src[:]
			temp_src = temp_src[:0]
		}
		if handled_length+c.outSize < len(dst) {
			c.block.Encrypt(dst[handled_length:handled_length+c.outSize], temp)
		} else {
			c.block.Encrypt(dst[handled_length:], temp)
		}
		handled_length += c.outSize
	}
}
