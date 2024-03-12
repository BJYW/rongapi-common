package utils

import (
	"fmt"
	"github.com/pborman/uuid"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var IdCardAlpha = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2, 0}
var IdCardCheckSum = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

var Random = []int64{77006791947779465, 74665223082153637, 29484611666145882,
	37200794235010091, 16589616287113976, 34824724549167383, 5394647632969764,
	43635317331776162, 94385949183117224, 75422040480279476, 14588224549167383}

var RandomAlpha = []string{"a", "b", "c", "d", "e", "f"}

func Encode(idcode string) string {
	res := idcode
	if len(idcode) == 15 {
		idcode = ConvertOldIdCard(idcode).CardNum
	}
	if len(idcode) == 18 {
		num, err := strconv.ParseInt(idcode[:17], 10, 64)
		var index int64
		index, err = strconv.ParseInt(idcode[17:], 10, 8)
		if err != nil {
			index = 10
		}
		num = num<<2 + Random[index]
		cont := strconv.FormatInt(num, 16)
		offset := idcode[17:]
		if strings.ToLower(offset) == "x" {
			offset = RandomAlpha[rand.Int31n(6)]
		}
		ci := 8
		if len(cont) == 14 {
			offset = RandomAlpha[rand.Int31n(6)] + offset
		} else {
			ci = 7
			offset = fmt.Sprintf("%d%s", rand.Int31n(9), offset)
		}
		uid := uuid.NewRandom().String()
		return uid[:ci] + cont + offset + uid[28:36]
	}
	return res
}

func Decode(codex string) string {
	if len(codex) == 32 {
		offset := 8
		last := string(codex[23])
		if codex[22] > 57 {
			offset = 8
		} else {
			offset = 7
		}
		index, err := strconv.ParseInt(string(codex[23]), 10, 8)
		if err != nil {
			index = 10
			last = "X"
		}

		contx, _ := strconv.ParseInt(codex[offset:22], 16, 64)

		tmpx := (contx - Random[index]) >> 2
		return fmt.Sprintf("%d%s", tmpx, last)
	} else {
		return codex
	}
}

type IdCard struct {
	CardNum string
	idNum   [18]byte
}

func ParseIdCard(idcard string) *IdCard {
	if len(idcard) == 15 {
		return ConvertOldIdCard(idcard)
	} else if len(idcard) == 18 {
		res := new(IdCard)
		copy(res.idNum[:], idcard)
		if res.idNum[17] != checkSum(res.idNum) {
			return nil
		} else {
			res.CardNum = idcard
			return res
		}
	} else {
		return nil
	}
}

func ConvertOldIdCard(idcard string) *IdCard {
	res := new(IdCard)
	copy(res.idNum[0:6], idcard[0:6])
	copy(res.idNum[8:17], idcard[6:15])
	if res.idNum[8] >= 2 {
		copy(res.idNum[6:8], []byte{49, 57})
	}
	res.idNum[17] = checkSum(res.idNum)
	res.CardNum = string(res.idNum[:])
	return res
}

//check  request

func CheckName(name string) (string, bool) {
	if len(name) < 2 {
		return "", false
	}
	return strings.Replace(name, ".", "·", -1), true
}

//CheckIdCode 判断身份证号是否合法，返回true合法，返回false不合法
func CheckIdCode(idcode string) bool {
	fmt.Println(string(idcode))

	length := len(idcode)
	if length != 15 && length != 18 {
		return false
	}
	if length == 18 {
		var bytes [18]byte
		copy(bytes[:], idcode)
		chk := checkSum(bytes)
		fmt.Println(string(chk))
		return (bytes[17] != 'x' && bytes[17] == chk) || (bytes[17] == 'x' && chk == 'X')
	}
	//默认返回错误
	return false
}

func checkSum(idnum [18]byte) byte {
	sum := 0
	for k, v := range idnum {
		sum += (int(v) - 48) * IdCardAlpha[k]
	}
	i := sum % 11
	return IdCardCheckSum[i]
}

func (i *IdCard) GetAge() int {
	ye := time.Now().Year()
	br := 0
	for k, v := range i.idNum[6:10] {
		br += (int(v) - 48) * int(math.Pow10(3-k))
	}
	return ye - br
}

func (i *IdCard) GetAgeLevel(mode int) int {
	age := i.GetAge()
	switch mode {
	default:
		l := age/5 - 7
		return int(math.Abs(math.Abs(float64(l)) - 3))

	}
	return 0
}
