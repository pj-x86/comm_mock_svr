// myutil_test.go
package myutil

import (
	"fmt"
	"log"
	"testing"
)

func TestRightPad(t *testing.T) {
	var str []byte
	str = []byte("123456")
	fmt.Printf("str=[%s],after RightPad, str=[%s]\n", str, RightPad(str, 8))
}

func TestLeftPadZero(t *testing.T) {
	var str []byte
	str = []byte("123456")
	fmt.Printf("str=[%s],after LeftPadZero, str=[%s]\n", str, LeftPadZero(str, 10))
}

func TestGetInterface(t *testing.T) {
	serviceXmlFile := "../service.xml"

	err := ReadInterfaces(serviceXmlFile)
	if err != nil {
		log.Printf("读入接口列表到内存失败\n")
	}
	log.Printf("%#v\n", interfaceListMap)

	var respBodyArr []Field
	respBodyArr, _ = GetInterface("03", "resp_body")
	for i := range respBodyArr {
		log.Printf("funcno=%s, name=%s,len=%d,value=[%s]\n", "03", respBodyArr[i].Name, respBodyArr[i].Len, respBodyArr[i].Value)
	}

	var requestArr []Field
	requestArr, _ = GetInterface("04", "resp_head")
	for i := range requestArr {
		log.Printf("funcno=%s, name=%s,len=%d,value=[%s]\n", "04", requestArr[i].Name, requestArr[i].Len, requestArr[i].Value)
	}
}

func TestGetMD5Hash(t *testing.T) {
	var data = []byte("this is a picture.")

	md5Out := GetMD5Hash(data)
	log.Printf("MD5 signature: [%s]\n", md5Out)
}

func TestLeftPadZeroForFloat(t *testing.T) {
	s := LeftPadZeroForFloat([]byte("1000.56"), 16, 2)
	fmt.Printf("s=%s\n", s)
}
