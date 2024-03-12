package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestEnc(t *testing.T) {
	if c, err := New("des/none/PKCS5Padding", map[string]interface{}{
		"key":             "RNQDK",
		"encoding_bytes":  "raw",
		"auto_expend_key": true,
	}); err == nil {
		s, err := c.Encrypt([]byte("{\"SvcNum\":\"123\",\"optrid\":\"123\",\"password\":\"123\",\"imei\":\"123456789\"}"))
		fmt.Println(len(s))
		for _, v := range s {
			fmt.Printf("%02x", v)
		}
		fmt.Println(err)
		e, _ := Encrypt("test")
		fmt.Println(len(e))
		fmt.Println(e)
		fmt.Println(Decrypt(e))
	} else {
		t.Fatal(err)
	}
}

func TestDec(t *testing.T) {
	if c, err := New("desede/cbc/PKCS5Padding", map[string]interface{}{
		"key":             "test",
		"iv":              defaultIV,
		"encoding_bytes":  "md5",
		"auto_expend_key": true,
	}); err == nil {
		src := []byte("ohX0859JwNtruOX0uBcmO3TwU1o00cgJvXkaKDUJTHbOhkLLbjT6KA==")
		crypted := make([]byte, len(src))
		fmt.Println(len(crypted), "---1")
		if n, err := base64.StdEncoding.Decode(crypted, src); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(n)
			crypted = crypted[:n]
		}
		s, err := c.Decrypt(crypted)
		s, err = hex.DecodeString(strings.ToLower(string(s)))
		fmt.Println(string(s), err)
		if n, err := base64.StdEncoding.Decode(s, s); err == nil {
			fmt.Println(string(s[:n]))
		} else {
			fmt.Println(err)
		}
		// fmt.Println(Decrypt(string(src)))
	} else {
		t.Fatal(err)
	}
}

func TestPasswordEnc(t *testing.T) {
	pwd := "test"
	if c, err := New("desede/cbc/PKCS5Padding", map[string]interface{}{
		"key":             pwd,
		"iv":              defaultIV,
		"encoding_bytes":  "md5",
		"auto_expend_key": true,
	}); err == nil {
		md5pwd := md5String(pwd)
		src := []byte(md5pwd)
		if encPwd, err := c.EncryptToString(src); err == nil {
			if encPwd != "ohX0859JwNtruOX0uBcmO3TwU1o00cgJvXkaKDUJTHbOhkLLbjT6KA==" {
				t.Fatal("Encrypt error")
			}
		} else {
			t.Error(err)
		}
		// fmt.Println(Decrypt(string(src)))
	} else {
		t.Fatal(err)
	}
}

func TestRsa(t *testing.T) {
	if c, err := New("rsa/none/noPadding", map[string]interface{}{
		"iv":                 "WWWID5CNWWWID5CN",
		"public_key":         PublicKey,
		"private_key":        PrivateKey,
		"public_key_encrypt": false,
		"encoding_bytes":     "raw",
		"auto_expend_key":    true,
	}); err == nil {
		data := []byte(`河南省郑州市长春路与玉兰街交叉口中移在线`)
		s, err := c.Encrypt(data)
		fmt.Println(base64.StdEncoding.EncodeToString(s), err)
		if dec, err := RsaDecryptWithPrivateKey(s, []byte(PrivateKey)); err != nil {
			t.Fatal(err)
		} else {
			if string(dec) != string(data) {
				log.Println("decrypt false", string(dec))
				t.Failed()
			} else {
				log.Println("decrypt ok", string(dec))
			}
		}
	} else {
		t.Fatal(err)
	}
}

var PublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqiy0H5Zowr3HAQXfygHIXsM9LgcMn+CQaQ6pwqEUGKPV+JoPPAtaAenlVDt20EWiPRb+dPeNw9+QqOOnYbVHrMscWeQgBPwGvqURz+I5oT1KwlN02LUfB8soXOjWzGdsRACOHwUe/85way4p6XiA8Kc5YUy7zHKmYO4asSPwnT2RTylbD5HMaCZsDaKdSyFrr1yd1oGG9UGk+txOFskzmqioV0386iZoq9X+B3pnUbtpytvL7O2g/d8wGebMaqOz3xC53Dp+Z4a/qwBPWNafFtzwcbXOKgfcTVwH+Xtk/3zBV6E9b1URIl9l+tcyiX+/cJrsYCxeIsNNqwK+HAQLGwIDAQAB
-----END PUBLIC KEY-----
`
var PrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCqLLQflmjCvccBBd/KAchewz0uBwyf4JBpDqnCoRQYo9X4mg88C1oB6eVUO3bQRaI9Fv50943D35Co46dhtUesyxxZ5CAE/Aa+pRHP4jmhPUrCU3TYtR8Hyyhc6NbMZ2xEAI4fBR7/znBrLinpeIDwpzlhTLvMcqZg7hqxI/CdPZFPKVsPkcxoJmwNop1LIWuvXJ3WgYb1QaT63E4WyTOaqKhXTfzqJmir1f4HemdRu2nK28vs7aD93zAZ5sxqo7PfELncOn5nhr+rAE9Y1p8W3PBxtc4qB9xNXAf5e2T/fMFXoT1vVREiX2X61zKJf79wmuxgLF4iw02rAr4cBAsbAgMBAAECggEBAI0PCS+1yZjDHpYmfF5CJAkI5Uml5j0QNCQhV25Rfwr9o3uymDY5Yp57dvDfaZkyX00sJjhmi0h1pL+aFUDUt1jv88w/YJtBHYYh5iXfa83PtIxoCU8os1QOfqKlDofYotGrdCMaZ5u7T0xIs80nahoRQ2r35COQVYG7XLzBAc4TvtRjVznLkNE2jHHFZrzOKsL/2epJaYfzfBZTUD7qys0r7eM0P5SrcLDsUOBf1eQDM11VN3Z7QAxbKS8EIDijBrup2y1tB6JCnfoWzvO47obTO2l5G94/vZN+Rh+SGDhAxM/MIbgLkeJRwinc5dIjjjFKH/7SsDFFGscZkJOCPEECgYEA+wP1dfN/mBthn3OCQoXy0eESLPs+OnnsOb2UAwbWhtJt3C3g0sYVIX9boheLLVFKFteXxKb8sYvar/rcGYWb4NqEUrUkE7rU1Wc1c3+vefkQgQDpGw+C6WxTu2R1CN0dAKFNJhAt8Y7am/IRwQ8GtiC6ZB0wCRJdtnEOh3tY3KMCgYEArY3J/4z3+MmFJ36doNzL8VH7gxyn/t7WWg8CQPxnuXfg8dPzecIlpffOg+ogF/WnyMJndmAvrLMqm45wNXugK/7L5Q3nNi5aQxnUNBKR9OF4gVNpA+WqiURQu0dX8W7X4MEbLSlK7dKSvF6xgGl7irZcvy+J7+WjpO3kpRoCxykCgYEA3hsJpZ2Yh/mqbtFhMAWf7s/NKloBHPLm78xsXCU/yzoFFW7O4RF+fw5XAVz9vFKSURXg55OuQldyYIaDHF4ZqGMR9LDtRDvruQxCwnj3xuB4Fv28RUJ7XpmImWZxHC+ySu1u+0ATtwrlKlmOFhp6oMriwlUFPPGwVA9DgOPGpCsCgYEAl8Z0zbCZNw/9C3H2JbSSAeOBQcLww7pa2H1hqjVR/Gdy5+lEdsmrpP4Ws1QvZuSlK7OfSW2Ko8w3ybwCMzN++pqj8xd6iO52b1IKOtO2ouH/QZIRN6BEyDBfLbiRlwkpb1tRCeifp2XMf0GMl1EoAkHC0EuhpJM9zkFhAGXD4ikCgYACznL3/xEcE8QCo9fI0ZrbvxMjvILOR9TEK/SCiBYhU33dT+oWkk5ById+n8NEfFdecQTubTHWtsaafJkO0WYXvC7SQuSNBNgmqccmWTB0KnfuM7A6VqaaRzQD/FsGjHnUpzXatBiMQ0HIoJdRkZImduEdF7Dg6gNb7qxziyFEtw==
-----END RSA PRIVATE KEY-----
`

func RsaDecryptWithPrivateKey(ciphertext, privateKey []byte) (result []byte, err error) {
	if block, _ := pem.Decode(privateKey); block != nil {
		var priInterface interface{}
		if priInterface, err = x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			priv := priInterface.(*rsa.PrivateKey)
			max := (priv.N.BitLen() + 7) / 8
			for len(ciphertext) > 0 {
				var temp []byte
				if len(ciphertext) > max {
					temp = ciphertext[:max]
					ciphertext = ciphertext[max:]
				} else {
					temp = ciphertext[:]
					ciphertext = ciphertext[:0]
				}
				if out, err := priv.Decrypt(rand.Reader, temp, new(rsa.PKCS1v15DecryptOptions)); err == nil {
					result = append(result, out...)
				} else {
					return nil, err
				}
			}
		}
	} else {
		err = errors.New("private key error!")
	}
	return
}

// func TestEncrpty(t *testing.T) {
// 	data := []byte(`{"name":"%E5%BC%A0%E5%A4%A7%E6%B5%B7","password":"maimaitivip168188","idNumber":"413026199107207276","loginName":"maimaitivip"}`)
// 	fmt.Println(len(data))
// 	res, err := RsaEncryptWithPublickKey(data, []byte(PublicKey))
// 	if err != nil {
// 		log.Println(err)
// 		t.Failed()
// 	} else {
// 		for i := 0; i < len(res); i++ {
// 			fmt.Printf("%d ", int8(res[i]))
// 		}
// 	}
// 	log.Println("----------")
// 	log.Println(base64.StdEncoding.EncodeToString(res))
// 	log.Println("----------")

// 	if dec, err := RsaDecryptWithPrivateKey(res, []byte(PrivateKey)); err != nil {
// 		t.Fatal(err)
// 	} else {
// 		if string(dec) != string(data) {
// 			log.Println("decrypt false", string(dec))
// 			t.Failed()
// 		}
// 	}

// 	signdata := `WKEAUnIgzp5CEJ3ZkRXWZ7Jxvt7GGdycoMkJewai5X/JhMJlaLHkbeF2yU24ZZKIuG+SgBLLXyrEeiFT+VI7gqPR8UF7JHNj3or3FygIDoNYWgykQrGYXDva2zgFaJHW8BfikvLkKFQTOA2f2PJIJ2Hp0mO3FKnTJZC8DFKIs0Y=`
// 	sign, _ := RsaSignWithPrivateKey([]byte(signdata), []byte(PrivateKey))
// 	signres := `LFIsQodqjq4kLmy2j1+PrQ+MA7DGsApXEuiHrzq40w+xZDN1KsvL5ESe6HzwSTqYUG6XJ2WkxLKyXXBHO8VhbiU1SMWsTnQNq7eert7Tr/lk9GXVVnjfBc9EDrrya2FTr6UJO1SlPzeAFCgTtqQEL4DqHqV1qxXb4pHhrPmLmvc=`
// 	if base64.StdEncoding.EncodeToString(sign) != signres {
// 		log.Println("false")
// 		t.Failed()
// 	}
// }

func TestRegexp(t *testing.T) {
	var wg sync.WaitGroup
	executeTime := time.Now().UnixNano()
	n := runtime.NumCPU()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			r, _ := regexp.Compile("a+")
			str := []byte("aaaaaa")
			for j := 0; j < 100000/n; j++ {
				r.Match(str)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println("execute time : ", (time.Now().UnixNano()-executeTime)/1e6, " MS")
}
