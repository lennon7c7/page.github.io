package ffmpeg_test

import (
	"fmt"
	"page.github.io/pkg/ffmpeg"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/ffmpeg/ffmpeg_test.go -run TestImg2Video
func TestImg2Video(t *testing.T) {
	input := "../../images/nature-%d.jpg"
	output := "../../images/test.mp4"
	ffmpeg.Img2Video(input, output)
}

// go test -v pkg/ffmpeg/ffmpeg_test.go -run TestAddAudio2Video
func TestAddAudio2Video(t *testing.T) {
	inputVideo := "../../images/test.mp4"
	inputAudio := "../../../audio/01-The Beatles-yesterday.mp3"
	output := "../../images/test-music-01-The Beatles-yesterday.mp4"
	ffmpeg.AddAudio2Video(inputVideo, inputAudio, output)

	inputVideo = "../../images/test.mp4"
	inputAudio = "../../../audio/11-The Eagles-Hotel California.mp3"
	output = "../../images/test-music-11-The Eagles-Hotel California.mp4"
	ffmpeg.AddAudio2Video(inputVideo, inputAudio, output)
}

// go test -v pkg/ffmpeg/ffmpeg_test.go -run TestAudio2Video
func TestAudio2Video(t *testing.T) {
	inputImg := `D:\Download\novel\让你模拟作案，你上演完美犯罪？\banner.jpg`
	inputAudio := `D:\Download\第100章 腐烂的脸.mp3`
	outputVideo := `D:\Download\第100章 腐烂的脸.mp4`
	ffmpeg.Audio2Video(inputImg, inputAudio, outputVideo)
}

func TestConcatVideo2Video(t *testing.T) {
	inputVideo := "../../images/test/coser@白银81 纯白圣女"
	outputVideo := ""
	t.Log(ffmpeg.ConcatVideo2Video(inputVideo, outputVideo))
}
