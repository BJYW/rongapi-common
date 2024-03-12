package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type rsaDescrypter struct {
	public                    *rsa.PublicKey
	private                   *rsa.PrivateKey
	use_public_key_to_encrypt bool
}

func (d *rsaDescrypter) createBlock() (block, error) {
	return d, nil
}

func (d *rsaDescrypter) BlockSize() int {
	if d.use_public_key_to_encrypt {
		return (d.public.N.BitLen()+7)/8 - 11
	} else {
		return (d.private.N.BitLen()+7)/8 - 11
	}
}

func (d *rsaDescrypter) outBlockSize() int {
	return d.BlockSize() + 11
}

func (d *rsaDescrypter) Encrypt(dst, src []byte) {
	if d.use_public_key_to_encrypt {
		if out, err := rsa.EncryptPKCS1v15(rand.Reader, d.public, src); err == nil {
			copy(dst, out)
		}
	} else {
		if out, err := priKeyEncrypt(rand.Reader, d.private, src); err == nil {
			copy(dst, out)
		}
	}
}

func (d *rsaDescrypter) Decrypt(dst, src []byte) {
	if out, err := d.private.Decrypt(rand.Reader, src, new(rsa.PKCS1v15DecryptOptions)); err == nil {
		copy(dst, out)
	}
}

func (d *rsaDescrypter) init(attrs map[string]interface{}) (err error) {
	var public []byte
	var private []byte
	if public, err = get_bytes_attr("public_key", attrs); err != nil {
		return err
	}
	if private, err = get_bytes_attr("private_key", attrs); err != nil {
		return err
	}
	if d.use_public_key_to_encrypt, err = get_bool_attr("public_key_encrypt", attrs); err != nil {
		return err
	}
	if decoded, _ := pem.Decode(public); decoded != nil {
		var pubInterface interface{}
		if pubInterface, err = x509.ParsePKIXPublicKey(decoded.Bytes); err == nil {
			d.public = pubInterface.(*rsa.PublicKey)
		}
	} else {
		err = errors.New("public key error")
	}
	if block, _ := pem.Decode(private); block != nil {
		var priInterface interface{}
		if priInterface, err = x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			d.private = priInterface.(*rsa.PrivateKey)
		}
	} else {
		err = errors.New("private key error")
	}
	return
}
