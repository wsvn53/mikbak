package main

import (
	"encoding/base64"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var mikOpts struct {
	Host 	string	`short:"s" long:"server" required:"true" description:"Mikrotik RouterOS ip address, example: 127.0.0.1"`
	User	string	`short:"u" long:"user" required:"true" description:"Username to login."`
	Pass	string	`short:"p" long:"password" required:"true" description:"Password to login, support base64 encoded password with prefix 'B:'."`
	Output  string  `short:"o" long:"output" default:"." description:"Target directory to save backup file, by default will save to current directory."`
	Prefix  string  `long:"prefix" default:"ros" description:"Add prefix to backup filename."`
}

func main() {
	// options
	_, _ = flags.ParseArgs(&mikOpts, os.Args)
	if mikOpts.User == "" || mikOpts.Pass == "" || mikOpts.Host == "" {
		os.Exit(1)
	}

	// resolve output dir
	mikOpts.Output, _ = homedir.Expand(mikOpts.Output)

	// decode base64 encoded password
	if strings.HasPrefix(mikOpts.Pass, "B:") {
		encodedPass := string([]byte(mikOpts.Pass)[2:])
		decodedPass, err := base64.StdEncoding.DecodeString(encodedPass)
		checkErr(err, "Password", true)
		mikOpts.Pass = strings.TrimSuffix(string(decodedPass), "\n")
	}

	// ssh login
	sshConfig := &ssh.ClientConfig{
		Timeout: 10 * time.Second,
		User: mikOpts.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{ ssh.Password(mikOpts.Pass) },
	}
	sshClient, err := ssh.Dial("tcp", mikOpts.Host + ":22", sshConfig)
	checkErr(err, "SSH", true)
	defer sshClient.Close()
	sshSession, err := sshClient.NewSession()
	checkErr(err, "SSH", true)
	defer sshSession.Close()

	// system info
	systemOutput, err := sshSession.CombinedOutput("/system resource print")
	checkErr(err, "SysInfo", true)
	fmt.Println("==> System Information:")
	fmt.Println(strings.TrimSuffix(string(systemOutput), "\n"))

	// do backup
	fmt.Println("==> Backup:")
	backupDate := time.Now().Format("20060102")
	backupFile := mikOpts.Prefix + "-" + backupDate
	backupCommand := "/system backup save dont-encrypt=yes name=" + backupFile
	fmt.Println(" - Run:", backupCommand)
	sshSession, err = sshClient.NewSession()
	checkErr(err, "Backup", true)
	defer sshSession.Close()
	err = sshSession.Run(backupCommand)
	checkErr(err, "Backup", true)

	// sftp
	backupFile += ".backup"
	if _, err := os.Stat(mikOpts.Output); os.IsNotExist(err) {
		err := os.MkdirAll(mikOpts.Output, 0755)
		checkErr(err, "Mkdir", true)
	}
	mikOpts.Output = filepath.Join(mikOpts.Output, backupFile)
	sftpClient, err := sftp.NewClient(sshClient)
	checkErr(err, "Download", true)
	defer sftpClient.Close()

	// save file
	destFile, err := os.Create(mikOpts.Output)
	checkErr(err, "Download", true)
	defer destFile.Close()
	remoteFile, err := sftpClient.Open(backupFile)
	checkErr(err, "Remote", true)
	defer remoteFile.Close()
	n, err := io.Copy(destFile, remoteFile)
	checkErr(err, "Copy", true)
	fmt.Println(" - Saved:", mikOpts.Output, ",", n, "bytes")
	_ = destFile.Sync()
}

func checkErr(err error, tag string, abort bool) {
	if err == nil {
		return
	}

	fmt.Printf("\n[!] %s: %s\n", tag, err)
	if abort {
		os.Exit(1)
	}
}