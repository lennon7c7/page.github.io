package ffmpeg

import (
	"fmt"
	"os/exec"
)

func Img2Video(input string, output string) {
	command := "ffmpeg -framerate 1/2 -start_number 1 -i \"" + input + "\" -r 30 -y \"" + output + "\""
	msg, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(msg)
}

func AddAudio2Video(inputVideo string, inputAudio string, output string) {
	command := "ffmpeg -i \"" + inputVideo + "\" -stream_loop -1 -i \"" + inputAudio + "\" -shortest -map 0:v:0 -map 1:a:0 -c:v copy -y \"" + output + "\""
	fmt.Println(command)
	msg, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(msg)
}
