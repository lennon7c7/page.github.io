package unrar_test

import (
	"fmt"
	"page.github.io/pkg/unrar"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/unrar/unrar_test.go
func TestCommand(t *testing.T) {
	input := "../../test.rar"
	output := "../../images/test/1"
	unrar.Command(input, output)
}
