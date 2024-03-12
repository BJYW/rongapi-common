package crypt

import (
	"crypto/cipher"
	"crypto/des"
	"errors"
	"fmt"
	"strings"
)

type wrappedCipherBlock struct {
	cipher.Block
}

func (w wrappedCipherBlock) outBlockSize() int {
	return w.Block.BlockSize()
}

type TripleDesCrypter struct {
	iv  []byte
	key []byte
}

func (d *TripleDesCrypter) createBlock() (block, error) {
	if res, err := des.NewTripleDESCipher(d.key); err == nil {
		return wrappedCipherBlock{res}, nil
	} else {
		return nil, err
	}
}

func (d *TripleDesCrypter) init(attrs map[string]interface{}) (err error) {
	d.iv, err = get_bytes_attr("iv", attrs)
	if key, ok := attrs["key"]; ok {
		switch ti := key.(type) {
		case []byte:
			d.key = ti
		case string:
			auto_expend := false
			if auto_expend_s, ok := attrs["auto_expend_key"]; ok {
				switch aes := auto_expend_s.(type) {
				case string:
					auto_expend = strings.ToLower(aes) == "true"
				case bool:
					auto_expend = aes
				}
			}
			if encoding, ok := attrs["encoding_bytes"]; ok {
				if se, ok := encoding.(string); ok {
					switch se {
					case "raw":
						if auto_expend {
							ti = paddingKey(ti, 64, string([]byte{0xD}))
						}
						d.key = []byte(ti)
					default:
						if auto_expend {
							ti = paddingKey(ti, 64, "D")
						}
						d.key = hexMustDecode(ti)
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
	if len(d.key) > 24 {
		d.key = d.key[:24]
	} else {
		return errors.New("triple des need 24 bits of key,using auto_expend_key(bool) to expend the key of string type")
	}
	fmt.Println(d)
	return nil
}

type DesCrypter struct {
	key []byte
}

func (d *DesCrypter) createBlock() (block, error) {
	if res, err := des.NewCipher(d.key); err == nil {
		return wrappedCipherBlock{res}, nil
	} else {
		return nil, err
	}
}

func (d *DesCrypter) init(attrs map[string]interface{}) (err error) {
	if key, ok := attrs["key"]; ok {
		switch ti := key.(type) {
		case []byte:
			d.key = ti
		case string:
			auto_expend := false
			if auto_expend_s, ok := attrs["auto_expend_key"]; ok {
				switch aes := auto_expend_s.(type) {
				case string:
					auto_expend = strings.ToLower(aes) == "true"
				case bool:
					auto_expend = aes
				}
			}
			if encoding, ok := attrs["encoding_bytes"]; ok {
				if se, ok := encoding.(string); ok {
					switch se {
					case "raw":
						if auto_expend {
							ti = paddingKey(ti, 64, string([]byte{0}))
						}
						d.key = []byte(ti)
					default:
						if auto_expend {
							ti = paddingKey(ti, 64, "D")
						}
						d.key = hexMustDecode(ti)
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
	if len(d.key) > 8 {
		d.key = d.key[:8]
	} else {
		return errors.New("triple des need 8 bits of key,using auto_expend_key(bool) to expend the key of string type")
	}
	return nil
}
