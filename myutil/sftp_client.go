// sftp_client.go sftp 客户端的二次封装
package myutil

import (
	"fmt"
	"log"

	//"net"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

// Connect 基于ssh连接sftp服务器，成功返回sftp客户端会话
func Connect(host string, port int, user string, password string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		//需要验证服务端，不做验证返回nil就可以
		//		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		//			return nil
		//		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		log.Printf("call ssh.Dial() fail. cause=%s\n", err.Error())
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Printf("call sftp.NewClient() fail. cause=%s\n", err.Error())
		return nil, err
	}

	log.Printf("connect to sftp server succeed!\n")
	return sftpClient, nil
}

// UploadFile 上传本地文件 localFile 到sftp远程目录 remoteDir 下
func UploadFile(sftpClient *sftp.Client, localFile string, remoteDir string) error {
	srcFile, err := os.Open(localFile)
	if err != nil {
		log.Printf("打开文件[%s]失败, 原因: [%s]\n", localFile, err.Error())
		return err
	}
	defer srcFile.Close()

	//远程目录在sftp上若不存在，自动创建
	err = sftpClient.MkdirAll(remoteDir)
	if err != nil {
		log.Printf("创建sftp目录[%s]失败, 原因: [%s]\n", remoteDir, err.Error())
		return err
	}

	var remoteFileName = path.Base(localFile) //取本地文件的文件名
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Printf("在sftp服务器上创建文件[%s]失败, 原因: [%s]\n", remoteDir+remoteFileName, err.Error())
		return err
	}
	defer dstFile.Close()

	//上传文件内容
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	log.Printf("copy file to remote server finished!\n")
	return nil
}

// DownloadFile 下载sftp远程文件 remoteFile 到本地目录 localDir
func DownloadFile(sftpClient *sftp.Client, remoteFile string, localDir string) error {
	srcFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		log.Printf("打开sftp远程文件[%s]失败, 原因: [%s]\n", remoteFile, err.Error())
		return err
	}
	defer srcFile.Close()

	var localFileName = path.Base(remoteFile) //获取文件名
	dstFile, err := os.Create(path.Join(localDir, localFileName))
	if err != nil {
		log.Printf("创建本地文件[%s]失败, 原因: [%s]\n", localDir+localFileName, err.Error())
		return err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		log.Printf("数据写入本地文件[%s]失败, 原因: [%s]\n", localDir+localFileName, err.Error())
		return err
	}

	log.Printf("copy file from remote server finished!\n")
	return nil
}

//var host = flag.String("host", "127.0.0.1", "sftp server's ip")
//var port = flag.Int("port", 22, "sftp server's port")
//var user = flag.String("user", "test", "sftp login name")
//var password = flag.String("password", "", "password")

//func main() {
//	flag.Parse()

//	var sftpClient *sftp.Client

//	sftpClient, err := Connect(*host, *port, *user, *password)
//	if err != nil {
//		log.Printf("连接sftp服务器[%s@%s:%d]失败\n", *user, *host, *port)
//		os.Exit(1)
//	}
//	defer sftpClient.Close()

//	localFile := "test.txt"
//	remoteDir := "20180717/"
//	err = UploadFile(sftpClient, localFile, remoteDir)
//	if err != nil {
//		log.Printf("上传文件[%s]到sftp失败\n", localFile)
//		os.Exit(1)
//	}

//	remoteFile := "20180717/abc.txt"
//	localDir := ""
//	err = DownloadFile(sftpClient, remoteFile, localDir)
//	if err != nil {
//		log.Printf("下载文件[%s]到本地[%s]失败\n", remoteFile, localDir)
//		os.Exit(1)
//	}

//	return
//}
