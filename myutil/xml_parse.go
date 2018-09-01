// xml_parse.go
package myutil

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"strconv"
)

type Field struct {
	//Order     int
	Name      string
	Len       int
	Precision int
	Type      string //字段类型, string/int/float
	Value     string //字段值
}

var interfaceListMap map[string]map[string][]Field //全局变量

// GetInterface 根据功能号 funcno 和报文打包类型 packType 获取对应的打包字段数组
// packType 取值: request, resp_head, resp_body
func GetInterface(funcno string, packType string) ([]Field, error) {
	return interfaceListMap[funcno][packType], nil
}

// ReadInterfaces 从配置文件 serviceXmlFile 读入全部接口列表，保存在全局变量 interfaceListMap 中
func ReadInterfaces(serviceXmlFile string) error {
	filePtr, err := os.Open(serviceXmlFile)
	if err != nil {
		log.Printf("打开配置文件[%s]失败, 原因: [%s]\n", serviceXmlFile, err.Error())
		return err
	}
	defer filePtr.Close()

	//创建一个接口列表map，以功能号为key
	interfaceListMap = make(map[string]map[string][]Field)
	var interfaceMap map[string][]Field
	var fieldArr []Field
	var field Field

	dec := xml.NewDecoder(filePtr)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("解析xml出错, 原因: [%s]\n", err.Error())
			return err
		}

		switch tok := tok.(type) {
		case xml.StartElement: //<name>标签
			if tok.Name.Local == "interface" { //开始一个接口
				//创建具体接口map，以request/response为key
				interfaceMap = make(map[string][]Field)

				//tok.Attr[0].Value (funcno的值)作为key
				interfaceListMap[tok.Attr[0].Value] = interfaceMap

			} else if tok.Name.Local == "request" || tok.Name.Local == "resp_head" || tok.Name.Local == "resp_body" { //请求报文

				//创建一个slice，保存字段列表
				fieldArr = make([]Field, 0)

			} else if tok.Name.Local == "field" { //报文字段

				for i := range tok.Attr {

					/*if tok.Attr[i].Name.Local == "order" {
						field.Order, _ = strconv.Atoi(tok.Attr[i].Value)
					} else */if tok.Attr[i].Name.Local == "name" {
						field.Name = tok.Attr[i].Value
					} else if tok.Attr[i].Name.Local == "len" {
						field.Len, _ = strconv.Atoi(tok.Attr[i].Value)
					} else if tok.Attr[i].Name.Local == "precision" {
						field.Precision, _ = strconv.Atoi(tok.Attr[i].Value)
					} else if tok.Attr[i].Name.Local == "type" {
						field.Type = tok.Attr[i].Value
					}
				}

			}

		case xml.EndElement: //</name>标签
			if tok.Name.Local == "request" || tok.Name.Local == "resp_head" || tok.Name.Local == "resp_body" {
				interfaceMap[tok.Name.Local] = fieldArr
			} else if tok.Name.Local == "field" { //报文字段
				fieldArr = append(fieldArr, field)
				//log.Printf("%#v\n", fieldArr)
			}
		case xml.CharData:
			if string(tok) != "\n" {
				field.Value = string(tok)
			} else {
				field.Value = ""
			}

		}

	} //end of for

	return nil
}

//func main() {
//	serviceXmlFile := "../service.xml"

//	err := ReadInterfaces(serviceXmlFile)
//	if err != nil {
//		log.Printf("读入接口列表到内存失败\n")
//		os.Exit(1)
//	}
//	log.Printf("%#v\n", interfaceListMap)

//	var respBodyArr []Field
//	respBodyArr, _ = GetInterface("03", "resp_body")
//	for i := range respBodyArr {
//		log.Printf("funcno=%s, name=%s,len=%d\n", "03", respBodyArr[i].Name, respBodyArr[i].Len)
//	}

//	var requestArr []Field
//	requestArr, _ = GetInterface("04", "request")
//	for i := range requestArr {
//		log.Printf("funcno=%s, name=%s,len=%d\n", "04", requestArr[i].Name, requestArr[i].Len)
//	}

//	return
//}
