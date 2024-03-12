package crypt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func New(mode string, attrs map[string]interface{}) (c *Cipher, err error) {
	mode = strings.ToLower(mode)
	modes := strings.Split(mode, "/")
	var cry crypter
	if cry, err = newCrypter(modes[0], attrs); err == nil {
		var blocker blocker
		if blocker, err = newBlocker(modes[1], attrs); err == nil {
			var padder padder
			if padder, err = newPadder(modes[2], attrs); err == nil {
				var c Cipher
				c.cry = cry
				c.bm = blocker
				c.pad = padder
				return &c, nil
			}
		}
	}
	return nil, err
}

type Cipher struct {
	cry crypter
	bm  blocker
	pad padder
}

func (c *Cipher) EncryptToString(bts []byte) (out string, err error) {
	if b, err := c.Encrypt(bts); err == nil {
		return base64.StdEncoding.EncodeToString(b), nil
	}
	return "", nil
}

func (c *Cipher) DecryptFromString(src string) (out []byte, err error) {
	var bts []byte
	if bts, err = base64.StdEncoding.DecodeString(src); err == nil {
		if bts, err = c.Decrypt(bts); err == nil {
			return bts, nil
		}
	}
	return nil, err
}

func (c *Cipher) Encrypt(bts []byte) (out []byte, err error) {
	if block, err := c.cry.createBlock(); err == nil {
		src := c.pad.padding(bts, block.BlockSize())
		out_size := len(src) / block.BlockSize() * block.outBlockSize()
		if len(src)-(len(src)/block.BlockSize()*block.BlockSize()) != 0 {
			out_size += block.outBlockSize()
		}
		crypted := make([]byte, out_size)
		blockMode := c.bm.createEncrypter(block)
		blockMode.CryptBlocks(crypted, src)
		return crypted, nil
	} else {
		return nil, err
	}
}

func (c *Cipher) Decrypt(src []byte) (out []byte, err error) {
	if block, err := c.cry.createBlock(); err == nil {
		fmt.Println(len(src))
		out_size := len(src) / block.outBlockSize() * block.BlockSize()
		if len(src)-(len(src)/block.BlockSize()*block.BlockSize()) != 0 {
			out_size += block.outBlockSize()
		}
		crypted := make([]byte, out_size)
		blockMode := c.bm.createDecrypter(block)
		blockMode.CryptBlocks(crypted, src)
		crypted = defaultPadder.unpadding(crypted)
		return crypted, nil
	} else {
		return nil, err
	}
}

type crypter interface {
	createBlock() (block, error)
}

func newCrypter(mode string, attrs map[string]interface{}) (crypter, error) {
	switch mode {
	case "rsa":
		rsa := new(rsaDescrypter)
		err := rsa.init(attrs)
		return rsa, err
	case "desede", "tripledes":
		des := new(TripleDesCrypter)
		err := des.init(attrs)
		return des, err
	case "des":
		des := new(DesCrypter)
		err := des.init(attrs)
		return des, err
	}
	return nil, errors.New("unsupported crypter type")
}

func newBlocker(mode string, attrs map[string]interface{}) (blocker, error) {
	switch mode {
	case "none":
		des := new(noBM)
		err := des.init(attrs)
		return des, err
	case "cbc":
		des := new(cbcBM)
		err := des.init(attrs)
		return des, err
	}
	return nil, errors.New("unsupported blocker type")
}

func newPadder(mode string, attrs map[string]interface{}) (padder, error) {
	switch mode {
	case "nopadding":
		des := new(noPadder)
		return des, nil
	case "pkcs5padding":
		des := new(pkcs5Padder)
		// err := des.init(attrs)
		return des, nil
	}
	return nil, errors.New("unsupported padder type")
}

func get_bytes_attr(name string, attrs map[string]interface{}) ([]byte, error) {
	if iv, ok := attrs[name]; ok {
		switch ti := iv.(type) {
		case []byte:
			return ti, nil
		case string:
			if encoding, ok := attrs["encoding_bytes"]; ok {
				if se, ok := encoding.(string); ok {
					switch se {
					case "raw":
						return []byte(ti), nil
					default:
						return hexMustDecode(ti), nil
					}
				} else {
					return nil, errors.New("encoding bytes must be string")
				}
			} else {
				return nil, errors.New("no encoding_bytes attr for decoding bytes from string")
			}
		}
	}
	return nil, errors.New("no " + name + " attr for des crypter")
}

func get_bool_attr(name string, attrs map[string]interface{}) (bool, error) {
	if iv, ok := attrs[name]; ok {
		switch ti := iv.(type) {
		case bool:
			return ti, nil
		case string:
			return strings.ToLower(ti) == "true", nil
		}
	}
	return false, errors.New("no " + name + " attr for des crypter")
}
