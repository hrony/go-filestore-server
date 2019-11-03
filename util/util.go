package util

import (
	"bytes"
	"os/exec"
)

// "root:123456@tcp(127.0.0.1:3306)/fileServer?charset=utf8"
func GetMysqlSource(username, password, host, port, db, charset string) string {
	return username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?charset=" + charset
}

// 执行 linux shell command
func ExecLinuxShell(s string) (string, error) {
	// 函数返回一个io.Writer类型的*Cmd
	cmd := exec.Command("/bin/bash", "-c", s)

	// 通过bytes.Buffer将byte类型转化为string类型
	var result bytes.Buffer
	cmd.Stdout = &result

	// Run执行cmd包含的命令，并阻塞直至完成
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return result.String(), err
}
