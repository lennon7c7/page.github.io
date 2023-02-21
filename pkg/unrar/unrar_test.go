package unrar_test

import (
	"fmt"
	"page.github.io/pkg/unrar"
	"testing"
	"time"
)

// go test -v pkg/unrar/unrar_test.go
func TestCommand(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	input := "../../test.rar"
	output := "../../images/test/1"
	unrar.Command(input, output)

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}
