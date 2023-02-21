package unrar

import (
	"fmt"
	"os"
	"os/exec"
)

func Command(input string, output string) {
	err := os.MkdirAll(output, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	command := "unrar e -pmrcong.com -inul -y " + input + " " + output
	_, err = exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
}
