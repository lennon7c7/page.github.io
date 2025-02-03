package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"page.github.io/pkg/file"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Img2Video(input string, output string) {
	err := os.MkdirAll(path.Dir(output), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	if file.Exists(output) {
		fmt.Println(output + " file exists, skip")
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

func Audio2Video(inputImg string, inputAudio string, outputVideo string) {
	command := "ffmpeg -loop 1 -i '" + inputImg + "' -i '" + inputAudio + "' -c:a copy -c:v libx264 -shortest '" + outputVideo + "'"

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
		fmt.Println(err, command)
		return
	}

	fmt.Println(string(msg))
}

func ConcatVideo2Video(inputVideoDir string, outputVideoFile string) (err error) {
	if outputVideoFile == "" {
		outputVideoBase := filepath.Base(inputVideoDir)
		outputVideoFile = inputVideoDir + "/" + outputVideoBase + ".mp4"
	}

	if file.Exists(outputVideoFile) {
		err = errors.New("输出文件已存在" + inputVideoDir)
		return
	}

	// 修复路径中的空格
	outputVideoFile, _ = filepath.Abs(outputVideoFile)
	outputVideoDir := filepath.Base(filepath.Dir(outputVideoFile))
	outputVideoBase := filepath.Base(outputVideoFile)
	outputVideoFile = strings.ReplaceAll(outputVideoFile, outputVideoDir, "\""+outputVideoDir+"\"")
	outputVideoFile = strings.ReplaceAll(outputVideoFile, outputVideoBase, "\""+outputVideoBase+"\"")

	var videoFiles []string
	for _, f := range file.GetFiles(inputVideoDir) {
		fileSuffix := path.Ext(f)
		fileSuffix = strings.ToLower(fileSuffix)
		if fileSuffix != ".mp4" && fileSuffix != ".mov" && fileSuffix != ".mkv" {
			continue
		}

		// 修复路径中的空格
		outputFilename, _ := filepath.Abs(f)
		baseDir := filepath.Base(filepath.Dir(outputFilename))
		outputFilename = strings.ReplaceAll(outputFilename, baseDir, "\""+baseDir+"\"")

		videoFiles = append(videoFiles, outputFilename)
	}

	if len(videoFiles) <= 1 {
		err = errors.New("视频文件数量不足" + inputVideoDir)
		return
	}

	minWidth, minHeight := 1080, 1920
	commandArg := `cd ` + filepath.Dir(outputVideoFile) + `; ffmpeg `
	for _, videoFile := range videoFiles {
		command := "ffmpeg -i " + videoFile
		commandArg += ` -i ` + filepath.Base(videoFile)
		var msg []byte
		switch runtime.GOOS {
		case "windows":
			msg, _ = exec.Command("powershell", command).CombinedOutput()
		case "linux":
			msg, _ = exec.Command("/bin/sh", "-c", command).CombinedOutput()
		default:
			err = errors.New("I don't support other os")
			return
		}

		re := regexp.MustCompile(`, (\d+)x(\d+), `)
		match := re.FindStringSubmatch(string(msg))
		if len(match) == 3 {
			width, _ := strconv.Atoi(match[1])
			height, _ := strconv.Atoi(match[2])
			if width < minWidth {
				minWidth = width
			}
			if height < minHeight {
				minHeight = height
			}
			//fmt.Printf("宽度：%d，高度：%d\n", width, height)
		} else {
			err = errors.New("未找到宽度和高度信息 " + videoFile)
			return
		}
	}

	commandArg += " -filter_complex \""
	for i := range videoFiles {
		commandArg += "[" + strconv.Itoa(i) + ":v]scale='min(" + strconv.Itoa(minWidth) + ",iw)':min'(" + strconv.Itoa(minHeight) + ",ih)':force_original_aspect_ratio=decrease,pad=" + strconv.Itoa(minWidth) + ":" + strconv.Itoa(minHeight) + ":(ow-iw)/2:(oh-ih)/2:black[v" + strconv.Itoa(i) + "]; "
	}
	for i := range videoFiles {
		commandArg += "[v" + strconv.Itoa(i) + "][" + strconv.Itoa(i) + ":a]"
	}
	commandArg += "concat=n=" + strconv.Itoa(len(videoFiles)) + ":v=1:a=1[outv][outa]\" -map \"[outv]\" -map \"[outa]\" " + outputVideoFile
	fmt.Println(commandArg)

	var msg []byte
	switch runtime.GOOS {
	case "windows":
		msg, err = exec.Command("powershell", commandArg).CombinedOutput()
	case "linux":
		msg, err = exec.Command("/bin/sh", "-c", commandArg).CombinedOutput()
	default:
		log.Fatalln("I don't support other os")
	}

	fmt.Printf("%v\n", outputVideoFile)
	msg = bytes.ReplaceAll(msg, []byte("\n"), []byte(""))
	return
}
