// datafile_parse.go 用于解析对账文件
package myutil

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// ReadDataFile 读取指定接口的数据文件 dataFileName 作为返回报文体
func ReadDataFile(funcno string, dataFileName string) []map[string]string {
	file, err := os.Open(dataFileName)
	if err != nil {
		log.Printf("打开文件[%s]失败，原因: [%s]\n", dataFileName, err.Error())
		return nil
	}
	defer file.Close()

	var resDataArr []map[string]string
	resDataArr = make([]map[string]string, 0)
	var lineno = 0

	var resData map[string]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() //返回一行
		lineno++
		if lineno == 1 { //第一行不是数据行，跳过
			continue
		}
		//log.Printf("line=[%s]\n", line)
		fields := strings.Split(line, "|")
		resData = make(map[string]string)

		if funcno == "0A" {
			resData["busin_type"] = fields[0]
			resData["fund_serial_no"] = fields[1]
			resData["bank_serial_no"] = fields[2]
			resData["bank_date"] = fields[3]
			resData["bank_time"] = fields[4]
			resData["req_date"] = fields[5]
			resData["order_status"] = fields[6]
			resData["occur_bala"] = fields[7]
			resData["bank_acco"] = fields[8]
			resData["explain"] = fields[9]
		} else if funcno == "0J" {
			resData["acco_id"] = fields[0]
			resData["other_serial_no"] = fields[1]
			resData["direct"] = fields[2]
			resData["name_in_bank"] = fields[3]
			resData["bank_name"] = fields[4]
			resData["bank_acco"] = fields[5]
			resData["home_name_in_bank"] = fields[6]
			resData["home_bank_acco"] = fields[7]
			resData["occur_bala"] = fields[8]
			resData["occur_date"] = fields[9]
			resData["occur_time"] = fields[10]
			resData["summary"] = fields[11]
			resData["usage"] = fields[12]
			resData["post_script"] = fields[13]
			resData["money_type"] = fields[14]
		} else if funcno == "0K" {
			resData["trade_acco"] = fields[0]
			resData["sign_no"] = fields[1]
			resData["bank_balance"] = fields[2]
			resData["fund_code"] = fields[3]
			resData["third_request_no"] = fields[4]
			resData["capital_source"] = fields[5]
			resData["request_bala"] = fields[6]
			resData["custom_name"] = fields[7]
			resData["identity_type"] = fields[8]
			resData["identity_no"] = fields[9]
			resData["name_in_bank"] = fields[10]
			resData["bank_name"] = fields[11]
			resData["bank_acco"] = fields[12]
			resData["trade_type"] = fields[13]
		}

		resDataArr = append(resDataArr, resData)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取文件出现异常，原因: [%s]", err.Error())
		return nil
	}

	return resDataArr
}
