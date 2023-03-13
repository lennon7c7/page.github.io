package ffmpeg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

func Img2Video(input string, output string) {
	err := os.MkdirAll(path.Dir(output), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	command := "ffmpeg -framerate 1/2 -start_number 1 -i \"" + input + "\" -r 30 -y \"" + output + "\""
	var msg []byte
	switch runtime.GOOS {
	case "windows":
		msg, err = exec.Command("powershell", command).Output()
	case "linux":
		msg, err = exec.Command("/bin/sh", "-c", command).Output()
	default:
		log.Fatalln("I don't support other os")
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(msg))
}

func AddAudio2Video(inputVideo string, inputAudio string, output string) {
	tempOutput := time.Now().Format("20060102150405") + path.Ext(output)
	command := "ffmpeg -i \"" + inputVideo + "\" -stream_loop -1 -i \"" + inputAudio + "\" -shortest -map 0:v:0 -map 1:a:0 -c:v copy -y \"" + tempOutput + "\""

	var msg []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		msg, err = exec.Command("powershell", command).Output()
	case "linux":
		msg, err = exec.Command("/bin/sh", "-c", command).Output()
	default:
		log.Fatalln("I don't support other os")
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	if inputVideo == output {
		_ = os.Remove(inputVideo)
	}
	_ = os.Rename(tempOutput, output)

	fmt.Println(string(msg))
}
