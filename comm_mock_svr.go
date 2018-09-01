// comm_mock_svr.go 是efund通讯机程序的模拟程序，用于直销系统单元测试、集成测试等场景。
// 通讯机请求/响应报文为定长报文，格式为: 长度(4byte)+报文内容+MD5(32byte)
// 长度 = 4 + len(报文内容+MD5)
// MD5 基于 报文内容+固定的盐值 生成
// 报文内容按HS_WEB接口格式定义，前2字节为功能号
// author: pj
// date  : 201807

package main

import (
	//"bufio"
	"comm_mock_svr/myutil"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var port = flag.Int("port", 6610, "listening port")
var output = flag.String("output", "file", "logging location, console|file")

// getFunctionNo 从 regmsg 中读取功能号，通过 funcno 返回
func getFunctionNo(reqmsg string) (funcno string, err error) {
	if len(reqmsg) < 2 {
		log.Printf("报文格式非法，长度不足")
		return "", fmt.Errorf("invalid request packet")
	}
	funcno = reqmsg[:2] //前面2个字节是功能号
	return funcno, nil
}

// packRespMsgEx 基于外部数据 data 生成功能号为 funcno 的返回报文头/报文体(单个)
// packType: 报文类型 resp_head:返回报文头; resp_body: 返回报文体
func packRespMsgEx(funcno string, packType string, data map[string]string, config *myutil.Config) []byte {

	var msgContentStr []byte
	msgContentStr = make([]byte, 0)

	respFieldArr, _ := myutil.GetInterface(funcno, packType)
	//打包返回报文头/报文体
	for i := range respFieldArr {
		//根据接口字段名取对应的返回值
		value := data[respFieldArr[i].Name]

		//中文转成GBK编码，以[]byte保存
		valueByte := []byte(value)
		//log.Printf("old_valueByte=[%s],len(valueByte)=[%d]\n", valueByte, len(valueByte))
		if config != nil && config.PacketCharSet == "GBK" {
			gbkmsg, err := myutil.UTF8ToGBK(valueByte)
			if err != nil {
				log.Printf("返回报文转成GBK编码失败[%s]", value)
				return []byte("")
			}

			valueByte = gbkmsg
		}

		//按接口字段长度要求进行补空格或补0处理
		switch respFieldArr[i].Type {
		case "string":
			valueByte = myutil.RightPad(valueByte, respFieldArr[i].Len)
		case "int":
			valueByte = myutil.LeftPadZero(valueByte, respFieldArr[i].Len)
		case "float":
			valueByte = myutil.LeftPadZeroForFloat(valueByte, respFieldArr[i].Len, respFieldArr[i].Precision)
		default:
			log.Printf("不支持的接口字段类型[%s]\n", respFieldArr[i].Type)
			return []byte("")
		}

		//log.Printf("valueByte=[%s],len(valueByte)=[%d]\n", valueByte, len(valueByte))
		//拼接成一个字符串
		msgContentStr = append(msgContentStr, valueByte...)
		//log.Printf("msgContentStr=[%s],len(msgContentStr)=[%d]\n", msgContentStr, len(msgContentStr))

	}

	return msgContentStr
}

// packRespMsg 基于接口配置文件中的预设数据生成功能号为 funcno 的返回报文头或报文体
// packType 取值: resp_head, resp_body
func packRespMsg(funcno string, packType string, config *myutil.Config) []byte {

	var msgContentStr []byte
	msgContentStr = make([]byte, 0)

	respFieldArr, _ := myutil.GetInterface(funcno, packType)
	//打包返回报文头/报文体
	for i := range respFieldArr {
		//根据接口字段名取对应的返回值
		value := respFieldArr[i].Value

		//中文转成GBK编码，以[]byte保存
		valueByte := []byte(value)
		//log.Printf("old_valueByte=[%s],len(valueByte)=[%d]\n", valueByte, len(valueByte))
		if config != nil && config.PacketCharSet == "GBK" {
			gbkmsg, err := myutil.UTF8ToGBK(valueByte)
			if err != nil {
				log.Printf("返回报文转成GBK编码失败[%s]", value)
				return []byte("")
			}

			valueByte = gbkmsg
		}

		//按接口字段长度要求进行补空格或补0处理
		switch respFieldArr[i].Type {
		case "string":
			valueByte = myutil.RightPad(valueByte, respFieldArr[i].Len)
		case "int":
			valueByte = myutil.LeftPadZero(valueByte, respFieldArr[i].Len)
		case "float":
			valueByte = myutil.LeftPadZeroForFloat(valueByte, respFieldArr[i].Len, respFieldArr[i].Precision)
		default:
			log.Printf("不支持的接口字段类型[%s]\n", respFieldArr[i].Type)
			return []byte("")
		}

		//log.Printf("valueByte=[%s],len(valueByte)=[%d]\n", valueByte, len(valueByte))
		//拼接成一个字符串
		msgContentStr = append(msgContentStr, valueByte...)
		//log.Printf("msgContentStr=[%s],len(msgContentStr)=[%d]\n", msgContentStr, len(msgContentStr))

	}

	return msgContentStr
}

// genRespMsg 生成最终返回报文，报文格式：长度(4byte)+报文内容+MD5(32byte)
func genRespMsg(msgcontent []byte, config *myutil.Config) []byte {

	resp := msgcontent

	//生成MD5
	var md5Out []byte
	var lenstr string
	if config.CheckMD5 == "Y" {
		md5Salt := []byte(config.MD5Salt)
		md5Input := append(resp, md5Salt...)
		md5Out = myutil.GetMD5Hash(md5Input) //基于报文内容拼接上MD5盐值计算MD5
		resp = append(resp, md5Out...)

		//计算报文长度
		lenstr = fmt.Sprintf("%04d", 4+len(resp)) //基于最后要发送的报文内容+MD5值计算
	}

	//将gbk报文转成utf-8报文，以便打印到日志中
	var utf8msg = []byte(msgcontent)
	if config != nil && config.PacketCharSet == "GBK" {
		var err error
		utf8msg, err = myutil.GBKToUTF8(msgcontent)
		if err != nil {
			log.Printf("返回报文转成UTF-8编码失败[%s]", msgcontent)
		}
	}
	log.Printf("响应报文: [%s]", lenstr+string(utf8msg)+string(md5Out))

	resp1 := []byte(lenstr)
	resp1 = append(resp1, resp...)

	return resp1
}

// genErrRespMsg 生成出错时的通用返回报文，此时尚未正确解析出请求报文
// 通用失败报文 报文长度(4byte)+失败标志(1byte)+错误代码(4byte)+错误原因(不定长)+MD5(32byte)
func genErrRespMsg(errcode string, errmsg string, config *myutil.Config) []byte {
	msgcontent := "1" + errcode + errmsg

	var resp []byte
	resp = []byte(msgcontent)
	//GBK编码转换
	if config != nil && config.PacketCharSet == "GBK" {
		gbkmsg, err := myutil.UTF8ToGBK(resp)
		if err != nil {
			log.Printf("返回报文转成GBK编码失败[%s]", resp)
			return resp
		}

		resp = gbkmsg
	}

	return genRespMsg(resp, config)
}

// handleConn 处理一个连接请求并返回相应的处理结果
func handleConn(c net.Conn) {
	defer c.Close()

	//加载配置文件
	config, err := myutil.GetConfig(myutil.ConfigFile)
	if err != nil {
		respmsg := genErrRespMsg("9999", "服务端内部异常", config)
		c.Write(respmsg)
		return
	}

	var serviceXmlFile = "service.xml"
	err = myutil.ReadInterfaces(serviceXmlFile)
	if err != nil {
		log.Printf("读入接口列表到内存失败\n")
		respmsg := genErrRespMsg("9999", "服务端内部异常", config)
		c.Write(respmsg)
		return
	}

	for true {

		//处理请求报文，先读长度
		var msglen [4]byte //报文前四字节表示报文长度
		length, errRead := c.Read(msglen[0:])
		if errRead != nil {
			if errRead != io.EOF {
				log.Printf("从连接读取请求报文长度失败[%s]", errRead.Error())
				//respmsg := genErrRespMsg("9999", "读取请求数据失败", config)
				//c.Write(respmsg)
			} else {
				log.Printf("客户端关闭连接，结束[%s]", errRead.Error())
			}

			return
		}

		if length != 4 { //不足4个字节，也不正确
			log.Printf("请求报文长度不足[%d]", length)
			respmsg := genErrRespMsg("9999", "读取请求数据失败", config)
			c.Write(respmsg)
			return
		}

		mlen, errLen := strconv.Atoi(string(msglen[:]))
		if errLen != nil {
			log.Printf("请求报文格式非法，前四字节必须为长度[%s]", string(msglen[:]))
			respmsg := genErrRespMsg("9999", "请求报文格式非法", config)
			c.Write(respmsg)
			return
		}

		//继续读取报文后续内容
		var data []byte
		data = make([]byte, mlen)
		length, errRead = c.Read(data)
		if errRead != nil {
			if errRead != io.EOF {
				log.Printf("从连接读取数据失败[%s]", errRead.Error())
				respmsg := genErrRespMsg("9999", "读取请求数据失败", config)
				c.Write(respmsg)

			} else {
				log.Printf("客户端关闭连接，结束", errRead.Error())
			}

			return
		}

		log.Printf("期望读取数据字节[%d], 实际读取数据字节[%d]", mlen, length)

		//处理报文编码
		//input := bufio.NewScanner(c)
		//for input.Scan() {
		var reqmsg string
		if config.PacketCharSet == "GBK" {
			//gbkmsg := input.Bytes()
			gbkmsg := data[:length]
			utf8msg, err := myutil.GBKToUTF8(gbkmsg)
			if err != nil {
				respmsg := genErrRespMsg("9999", "请求报文编码格式非法", config)
				c.Write(respmsg)
				return
			}
			log.Printf("原始报文编码为GBK，转为UTF8\n")
			reqmsg = string(utf8msg)
		} else { //UTF-8
			//reqmsg = input.Text()
			reqmsg = string(data[:length])
		}

		log.Printf("请求报文: [%s]", reqmsg)

		//解析功能号
		funcno, err1 := getFunctionNo(reqmsg)
		if err1 != nil {
			respmsg := genErrRespMsg("9999", err1.Error(), config)
			c.Write(respmsg)
			return
		}
		log.Printf("功能号＝%s\n", funcno)

		var respmsg = []byte("处理成功")
		var msgHeadStr = []byte("")
		var msgBodyStr = []byte("")

		//if config.WantResult == "0" { //要求按成功返回，则需要做后续处理，否则直接按失败报文打包即可

		//根据功能号进行具体的逻辑处理
		if strings.Contains(config.DirectRespFuncNo, funcno) {

			log.Printf("按照接口配置文件中的预设值进行返回\n")
			msgHeadStr = packRespMsg(funcno, "resp_head", config)
			msgBodyStr = packRespMsg(funcno, "resp_body", config)
		} else if strings.Contains(config.FileRespFuncNo, funcno) {

			var returnDirect = 0

			respHeadArr, _ := myutil.GetInterface(funcno, "resp_head")
			//检查接口配置文件中预设的返回结果
			for i := range respHeadArr {
				if respHeadArr[i].Name == "succ_flag" && respHeadArr[i].Value == "1" {
					//需要返回失败，无需返回报文体
					log.Printf("按照接口配置文件中的预设错误进行返回\n")
					returnDirect = 1
					break
				}
			}

			if returnDirect == 1 {
				msgHeadStr = packRespMsg(funcno, "resp_head", config) //基于接口配置文件中预设值返回
			} else {

				//从本地文件读入数据

				dataFileName := funcno + "_file.txt"

				ResData := myutil.ReadDataFile(funcno, dataFileName)
				if ResData == nil {
					log.Printf("读取对账文件[%s]失败\n", dataFileName)

					var respData map[string]string
					respData = make(map[string]string)
					respData["succ_flag"] = "1"
					respData["err_code"] = "9999"
					respData["err_msg"] = "读取对账文件失败"
					respData["total_num"] = "0"
					respData["return_num"] = "0"

					//只打包返回报文头
					msgHeadStr = packRespMsgEx(funcno, "resp_head", respData, config)
					break
				}

				log.Printf("读取数据文件[%s]成功\n", dataFileName)

				//循环打包返回报文体
				for i := range ResData {
					respBody := packRespMsgEx(funcno, "resp_body", ResData[i], config)
					msgBodyStr = append(msgBodyStr, respBody...)
				}

				//打包返回报文头
				var respData map[string]string
				respData = make(map[string]string)

				//基于接口配置文件的返回报文头打包
				for i := range respHeadArr {
					respData[respHeadArr[i].Name] = respHeadArr[i].Value
				}

				//返回记录数要以文件中的实际记录数返回
				respData["total_num"] = strconv.Itoa(len(ResData))
				respData["return_num"] = strconv.Itoa(len(ResData))

				msgHeadStr = packRespMsgEx(funcno, "resp_head", respData, config)
			}
		} else {
			//case "0M", "0P":
			//do nothing
			//			//从本地对账文件 datafile.json 读入数据
			//			//生成对账文件
			//			compFileName := "compfile_20180717.txt"
			//			compFile, err := os.Create(compFileName)
			//			if err != nil {
			//				log.Printf("创建本地文件[%s]失败, 原因: [%s]\n", compFileName, err.Error())
			//			}
			//			defer compFile.Close()

			//			//写对账数据
			//			compFile.Write([]byte("对账数据在此"))

			//			//连接到sftp
			//			sftpClient, err := myutil.Connect(config.SFTPHost, config.SFTPPort,
			//				config.SFTPUser, config.SFTPPasswd)
			//			if err != nil {
			//				log.Printf("连接sftp服务器[%s@%s:%d]失败\n", config.SFTPUser,
			//					config.SFTPHost, config.SFTPPort)

			//			}
			//			defer sftpClient.Close()

			//			//上传文件到sftp
			//			remoteDir := "ftp/test"
			//			err = myutil.UploadFile(sftpClient, compFileName, remoteDir)
			//			if err != nil {
			//				log.Printf("上传文件[%s]到sftp失败\n", compFileName)
			//				os.Exit(1)
			//			}

			log.Printf("非法的功能号[%s]", funcno)
			respmsg = genErrRespMsg("9999", "非法的功能号", config)
			c.Write(respmsg)
			return
		}
		//}

		//组装返回报文内容

		msgContentStr := append(msgHeadStr, msgBodyStr...)
		//生成最终返回报文
		respmsg = genRespMsg(msgContentStr, config)

		//发送给客户端
		c.Write(respmsg)
		log.Printf("发送返回报文成功")
		//}end of for
	} //end of for
}

func main() {
	flag.Parse()

	if *output != "console" { //输出到日志文件
		var logFileName = "server.log"
		//创建全局日志输出
		logFile, logErr := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if logErr != nil {
			fmt.Println("打开日志文件失败: ", logFileName)
			os.Exit(1)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
		log.SetFlags(log.Ldate | log.Ltime)
	} else { //输出到控制台
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	//启动侦听
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("开始监听tcp端口: ", strconv.Itoa(*port))
	//!+
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
	//!-
}
