// config_parse.go 用于解析 config.json 配置文件
package myutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

//对应 config.json 文件
type Config struct {
	PacketCharSet string `json:"pkt_charset"` //通讯报文的字符编码 GBK|UTF8
	//CompFileVer   string `json:"compfile_version"` //对账文件接口版本 1.0/1.1/1.2
	//WantResult    string `json:"want_result"` //期望返回的处理结果 0-表示返回成功; 1-表示返回失败
	//ErrCode       string `json:"err_code"`
	//ErrMsg        string `json:"err_msg"`
	//SFTPHost   string `json:"sftp_host"`
	//SFTPPort   int    `json:"sftp_port"`
	//SFTPUser   string `json:"sftp_user"`
	//SFTPPasswd string `json:"sftp_pwd"`
	DirectRespFuncNo string `json:"direct_resp_funcno"` //基于service.xml接口配置文件中的预设值直接返回的功能号列表
	FileRespFuncNo   string `json:"file_resp_funcno"`   //基于特定数据文件返回响应报文体的功能号列表
	CheckMD5         string `json:"check_md5"`          //是否需要进行MD5签名，Y-是;N-否
	MD5Salt          string `json:"md5_salt"`           //MD5签名使用的盐值
}

var ConfigFile = "config.json"

// stripComments 去除注释行，以#开头
func stripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // Windows
	lines := bytes.Split(data, []byte("\n"))                //split to multi lines
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^#\s*`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}

	return bytes.Join(filtered, []byte("\n")), nil
}

// GetConfig 读取 config.json 配置文件并保存到Config结构体
func GetConfig(filename string) (*Config, error) {
	config := new(Config)
	jsonStr, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("读取文件失败: %s, 错误: %s\n", filename, err.Error())
		return nil, err
	}

	jsonStr, err = stripComments(jsonStr) //去掉#注释
	if err != nil {
		log.Printf("对配置文件去注释处理失败: %s, 错误: %s\n", jsonStr, err.Error())
		return nil, err
	}

	err = json.Unmarshal(jsonStr, config)
	if err != nil {
		log.Printf("json字符串转换为go结构体失败: %s, 错误: %s\n", jsonStr, err.Error())
		return nil, err
	}

	log.Printf("配置文件加载成功，内容: %#v\n", config)

	return config, nil
}
