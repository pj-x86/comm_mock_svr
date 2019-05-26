// myutil.go 实现一些常用的功能，如左补零、右补空格对齐等
package myutil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"log"
	"math"
	"strconv"
	"strings"
)

// RightPad 在字符数组 str 右边补空格直到长度 length
func RightPad(str []byte, length int) []byte {
	var padStr []byte
	padLen := length - len(str)

	if padLen > 0 {
		padStr = []byte(strings.Repeat(" ", padLen))
	} else {
		padStr = []byte("")
	}

	return append(str, padStr...)
}

// LeftPadZero 在字符数组 str 左边补数字0直到长度 length ，str 通常是数值字段的字符串表示
func LeftPadZero(str []byte, length int) []byte {
	var padStr []byte
	padLen := length - len(str)

	if padLen > 0 {
		padStr = []byte(strings.Repeat("0", padLen))
	} else {
		padStr = []byte("")
	}

	return append(padStr, str...)
}

// LeftPadZeroForFloat 在字符数组 str 左边补数字0直到长度 length ，str 通常是浮点数的字符串表示。
// str 会先转成浮点数再乘以Pow10(precision)，再进行左补零操作。
func LeftPadZeroForFloat(str []byte, length int, precision int) []byte {
	//转成浮点数再乘以Pow10(precision)
	f, _ := strconv.ParseFloat(string(str), 64)
	f = f * math.Pow10(precision)
	s := strconv.FormatFloat(f, 'f', 0, 64)

	str = []byte(s)

	var padStr []byte
	padLen := length - len(str)

	if padLen > 0 {
		padStr = []byte(strings.Repeat("0", padLen))
	} else {
		padStr = []byte("")
	}

	return append(padStr, str...)
}

// GetMD5Hash 对字符数组 text 计算其MD5值，以大写返回32字节长的字符串值
func GetMD5Hash(text []byte) []byte {
	hasher := md5.New()
	hasher.Write(text)
	cipherText := hasher.Sum(nil) //生成16byte的MD5
	hexText := make([]byte, 32)
	hex.Encode(hexText, cipherText) //以十六进制对MD5进行编码，保存到32字节数组中
	log.Printf("待签名数据[%s]\n", text)
	return bytes.ToUpper(hexText) //转成大写
}
