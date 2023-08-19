package unrar

import (
	"fmt"
	"os"
	"os/exec"
	"page.github.io/pkg/log"
	"path/filepath"
	"runtime"
)

func Command(input string, output string) {
	input, _ = filepath.Abs(input)
	output, _ = filepath.Abs(output)

	err := os.MkdirAll(output, os.ModePerm)
	if err != nil {
		fmt.Println(err, input, output)
		return
	}

	var msg []byte
	switch runtime.GOOS {
	case "windows":
		// e 提取文件不带压缩路径
		// x 解压文件带完整路径
		command := `D:\"Program Files"\WinRAR\UnRAR.exe x ` + input + ` -p"mrcong.com" -inul -y ` + output
		msg, err = exec.Command("powershell", command).Output()
	case "linux":
		command := "unrar e -pmrcong.com -inul -y " + input + " " + output
		msg, err = exec.Command("/bin/sh", "-c", command).Output()
	default:
		log.Fatalln("I don't support other os")
	}
	if err != nil {
		fmt.Println(err, input, output)
		return
	}

	if string(msg) != "" {
		fmt.Println(string(msg), input, output)
	}
}
