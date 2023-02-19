package unrar

import (
	"fmt"
	"os/exec"
)

func Command(input string, output string) {
	command := "unrar e -pmrcong.com -inul -y " + input + " " + output
	_, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
}
